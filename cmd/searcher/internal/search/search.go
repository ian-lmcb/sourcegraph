// Package search is a service which exposes an API to text search a repo at
// a specific commit.
//
// Architecture Notes:
//  * Archive is fetched from gitserver
//  * Simple HTTP API exposed
//  * Currently no concept of authorization
//  * On disk cache of fetched archives to reduce load on gitserver
//  * Run search on archive. Rely on OS file buffers
//  * Simple to scale up since stateless
//  * Use ingress with affinity to increase local cache hit ratio
package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	zoektquery "github.com/google/zoekt/query"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	nettrace "golang.org/x/net/trace"

	"github.com/sourcegraph/sourcegraph/cmd/searcher/protocol"
	"github.com/sourcegraph/sourcegraph/internal/api"
	"github.com/sourcegraph/sourcegraph/internal/errcode"
	"github.com/sourcegraph/sourcegraph/internal/search/searcher"
	streamhttp "github.com/sourcegraph/sourcegraph/internal/search/streaming/http"
	zoektutil "github.com/sourcegraph/sourcegraph/internal/search/zoekt"
	"github.com/sourcegraph/sourcegraph/internal/trace/ot"
	"github.com/sourcegraph/sourcegraph/lib/errors"
	"github.com/sourcegraph/sourcegraph/lib/log"
	sglog "github.com/sourcegraph/sourcegraph/lib/log"
)

const (
	// numWorkers is how many concurrent readerGreps run in the case of
	// regexSearch, and the number of parallel workers in the case of
	// structuralSearch.
	numWorkers = 8
)

// Service is the search service. It is an http.Handler.
type Service struct {
	Store *Store
	Log   log.Logger

	// GitOutput returns the stdout of running git with args against the repo.
	//
	// TODO pick a design which doesn't directly depend on Command. Probably
	// adding a relevant function to the gitserver client. This is only used
	// by FeatHybrid.
	GitOutput func(ctx context.Context, repo api.RepoName, args ...string) ([]byte, error)
}

// ServeHTTP handles HTTP based search requests
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	running.Inc()
	defer running.Dec()

	var p protocol.Request
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&p); err != nil {
		http.Error(w, "failed to decode form: "+err.Error(), http.StatusBadRequest)
		return
	}

	if !p.PatternMatchesContent && !p.PatternMatchesPath {
		// BACKCOMPAT: Old frontends send neither of these fields, but we still want to
		// search file content in that case.
		p.PatternMatchesContent = true
	}
	if err := validateParams(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.streamSearch(ctx, w, p)
}

func (s *Service) streamSearch(ctx context.Context, w http.ResponseWriter, p protocol.Request) {
	if p.Limit == 0 {
		// No limit for streaming search since upstream limits
		// will either be sent in the request, or propagated by
		// a cancelled context.
		p.Limit = math.MaxInt32
	}
	eventWriter, err := streamhttp.NewWriter(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	matchesBuf := streamhttp.NewJSONArrayBuf(32*1024, func(data []byte) error {
		return eventWriter.EventBytes("matches", data)
	})
	onMatches := func(match protocol.FileMatch) {
		if err := matchesBuf.Append(match); err != nil {
			s.Log.Warn("failed appending match to buffer", log.Error(err))
		}
	}

	ctx, cancel, stream := newLimitedStream(ctx, p.Limit, onMatches)
	defer cancel()

	err = s.search(ctx, &p, stream)
	doneEvent := searcher.EventDone{
		LimitHit: stream.LimitHit(),
	}
	if err != nil {
		doneEvent.Error = err.Error()
	}

	// Flush remaining matches before sending a different event
	if err := matchesBuf.Flush(); err != nil {
		s.Log.Warn("failed to flush matches", log.Error(err))
	}
	if err := eventWriter.Event("done", doneEvent); err != nil {
		s.Log.Warn("failed to send done event", log.Error(err))
	}
}

func (s *Service) search(ctx context.Context, p *protocol.Request, sender matchSender) (err error) {
	tr := nettrace.New("search", fmt.Sprintf("%s@%s", p.Repo, p.Commit))
	tr.LazyPrintf("%s", p.Pattern)

	// TODO observability around FeatHybrid

	span, ctx := ot.StartSpanFromContext(ctx, "Search")
	ext.Component.Set(span, "service")
	span.SetTag("repo", p.Repo)
	span.SetTag("url", p.URL)
	span.SetTag("commit", p.Commit)
	span.SetTag("pattern", p.Pattern)
	span.SetTag("isRegExp", strconv.FormatBool(p.IsRegExp))
	span.SetTag("isStructuralPat", strconv.FormatBool(p.IsStructuralPat))
	span.SetTag("languages", p.Languages)
	span.SetTag("isWordMatch", strconv.FormatBool(p.IsWordMatch))
	span.SetTag("isCaseSensitive", strconv.FormatBool(p.IsCaseSensitive))
	span.SetTag("pathPatternsAreRegExps", strconv.FormatBool(p.PathPatternsAreRegExps))
	span.SetTag("pathPatternsAreCaseSensitive", strconv.FormatBool(p.PathPatternsAreCaseSensitive))
	span.SetTag("limit", p.Limit)
	span.SetTag("patternMatchesContent", p.PatternMatchesContent)
	span.SetTag("patternMatchesPath", p.PatternMatchesPath)
	span.SetTag("indexerEndpoints", p.IndexerEndpoints)
	span.SetTag("select", p.Select)
	defer func(start time.Time) {
		code := "200"
		// We often have canceled and timed out requests. We do not want to
		// record them as errors to avoid noise
		if ctx.Err() == context.Canceled {
			code = "canceled"
			span.SetTag("err", err)
		} else if err != nil {
			tr.LazyPrintf("error: %v", err)
			tr.SetError()
			ext.Error.Set(span, true)
			span.SetTag("err", err.Error())
			if errcode.IsBadRequest(err) {
				code = "400"
			} else if errcode.IsTemporary(err) {
				code = "503"
			} else {
				code = "500"
			}
		}
		tr.LazyPrintf("code=%s matches=%d limitHit=%v", code, sender.SentCount(), sender.LimitHit())
		tr.Finish()
		requestTotal.WithLabelValues(code).Inc()
		span.LogFields(otlog.Int("matches.len", sender.SentCount()))
		span.SetTag("limitHit", sender.LimitHit())
		span.Finish()
		s.Log.Debug("search request",
			sglog.String("repo", string(p.Repo)),
			sglog.String("commit", string(p.Commit)),
			sglog.String("pattern", p.Pattern),
			sglog.Bool("isRegExp", p.IsRegExp),
			sglog.Bool("isStructuralPat", p.IsStructuralPat),
			sglog.Strings("languages", p.Languages),
			sglog.Bool("isWordMatch", p.IsWordMatch),
			sglog.Bool("isCaseSensitive", p.IsCaseSensitive),
			sglog.Bool("patternMatchesContent", p.PatternMatchesContent),
			sglog.Bool("patternMatchesPath", p.PatternMatchesPath),
			sglog.Int("matches", sender.SentCount()),
			sglog.String("code", code),
			sglog.Duration("duration", time.Since(start)),
			sglog.Strings("indexerEndpoints", p.IndexerEndpoints),
			sglog.Error(err))
	}(time.Now())

	if p.IsStructuralPat && p.Indexed {
		// Execute the new structural search path that directly calls Zoekt.
		// TODO use limit in indexed structural search
		return structuralSearchWithZoekt(ctx, p, sender)
	}

	// Compile pattern before fetching from store incase it is bad.
	var rg *readerGrep
	if !p.IsStructuralPat {
		rg, err = compile(&p.PatternInfo)
		if err != nil {
			return badRequestError{err.Error()}
		}
	}

	if p.FetchTimeout == "" {
		p.FetchTimeout = "500ms"
	}
	fetchTimeout, err := time.ParseDuration(p.FetchTimeout)
	if err != nil {
		return err
	}
	prepareCtx, cancel := context.WithTimeout(ctx, fetchTimeout)
	defer cancel()

	getZf := func() (string, *zipFile, error) {
		path, err := s.Store.PrepareZip(prepareCtx, p.Repo, p.Commit)
		if err != nil {
			return "", nil, err
		}
		zf, err := s.Store.zipCache.Get(path)
		return path, zf, err
	}

	hybrid := !p.IsStructuralPat && p.FeatHybrid
	if hybrid {
		// TODO which ctx?
		indexed, ok, err := zoektIndexedCommit(ctx, p.IndexerEndpoints, p.Repo)
		if err != nil {
			return err
		}
		// TODO from this point onwards we should cache keyed by indexed + p.Commit
		if ok {
			out, err := s.GitOutput(ctx, p.Repo, "diff", "-z", "--name-status", "--no-renames", string(p.Commit), string(indexed))
			if err != nil {
				return err
			}

			slices := bytes.Split(bytes.TrimRight(out, "\x00"), []byte{0})
			if len(slices)%2 != 0 {
				return errors.New("uneven pairs")
			}

			var unindexedSearch []string
			var indexedIgnore []string
			for i := 0; i < len(slices); i += 2 {
				path := string(slices[i+1])
				switch slices[i][0] {
				case 'A':
					unindexedSearch = append(unindexedSearch, path)
				case 'M':
					unindexedSearch = append(unindexedSearch, path)
					indexedIgnore = append(indexedIgnore, path)
				case 'D':
					indexedIgnore = append(indexedIgnore, path)
				}
			}

			{
				// TODO race between calling List and searching now. Index may have changed.
				qText, err := zoektCompile(&p.PatternInfo)
				if err != nil {
					return errors.Wrap(err, "failed to compile query for zoekt")
				}
				q := zoektquery.Simplify(zoektquery.NewAnd(
					zoektquery.NewSingleBranchesRepos("HEAD", uint32(p.RepoID)),
					qText,
					zoektIgnorePaths(indexedIgnore),
				))

				k := zoektutil.ResultCountFactor(1, int32(p.Limit), false)
				opts := zoektutil.SearchOpts(ctx, k, int32(p.Limit), nil)
				if deadline, ok := ctx.Deadline(); ok {
					opts.MaxWallTime = time.Until(deadline) - 100*time.Millisecond
				}

				client := getZoektClient(p.IndexerEndpoints)
				res, err := client.Search(ctx, q, &opts)
				if err != nil {
					return err
				}
				for _, fm := range res.Files {
					var lineMatches []protocol.LineMatch
					var matchCount int
					for _, lm := range fm.LineMatches {
						var offs [][2]int
						for _, lf := range lm.LineFragments {
							offs = append(offs, [2]int{lf.LineOffset, lf.MatchLength})
						}
						lineMatches = append(lineMatches, protocol.LineMatch{
							Preview:          string(lm.Line),
							LineNumber:       lm.LineNumber - 1,
							OffsetAndLengths: offs,
						})
						matchCount += len(offs)
						if len(offs) == 0 {
							matchCount++
						}
					}
					sender.Send(protocol.FileMatch{
						Path:        fm.FileName,
						LineMatches: lineMatches,
						MatchCount:  matchCount,
					})
				}

				return nil
			}

			// TODO fetch paths unindexedSearch and search them
		}
	}

	zipPath, zf, err := getZipFileWithRetry(getZf)
	if err != nil {
		return errors.Wrap(err, "failed to get archive")
	}
	defer zf.Close()

	nFiles := uint64(len(zf.Files))
	bytes := int64(len(zf.Data))
	tr.LazyPrintf("files=%d bytes=%d", nFiles, bytes)
	span.LogFields(
		otlog.Uint64("archive.files", nFiles),
		otlog.Int64("archive.size", bytes))
	archiveFiles.Observe(float64(nFiles))
	archiveSize.Observe(float64(bytes))

	if p.IsStructuralPat {
		return filteredStructuralSearch(ctx, zipPath, zf, &p.PatternInfo, p.Repo, sender)
	} else {
		return regexSearch(ctx, rg, zf, p.PatternMatchesContent, p.PatternMatchesPath, p.IsNegated, sender)
	}
}

func validateParams(p *protocol.Request) error {
	if p.Repo == "" {
		return errors.New("Repo must be non-empty")
	}
	// Surprisingly this is the same sanity check used in the git source.
	if len(p.Commit) != 40 {
		return errors.Errorf("Commit must be resolved (Commit=%q)", p.Commit)
	}
	if p.Pattern == "" && p.ExcludePattern == "" && len(p.IncludePatterns) == 0 {
		return errors.New("At least one of pattern and include/exclude pattners must be non-empty")
	}
	if p.IsNegated && p.IsStructuralPat {
		return errors.New("Negated patterns are not supported for structural searches")
	}
	return nil
}

const megabyte = float64(1000 * 1000)

var (
	running = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "searcher_service_running",
		Help: "Number of running search requests.",
	})
	archiveSize = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "searcher_service_archive_size_bytes",
		Help:    "Observes the size when an archive is searched.",
		Buckets: []float64{1 * megabyte, 10 * megabyte, 100 * megabyte, 500 * megabyte, 1000 * megabyte, 5000 * megabyte},
	})
	archiveFiles = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "searcher_service_archive_files",
		Help:    "Observes the number of files when an archive is searched.",
		Buckets: []float64{100, 1000, 10000, 50000, 100000},
	})
	requestTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "searcher_service_request_total",
		Help: "Number of returned search requests.",
	}, []string{"code"})
)

type badRequestError struct{ msg string }

func (e badRequestError) Error() string    { return e.msg }
func (e badRequestError) BadRequest() bool { return true }
