package lsifstore

import (
	"context"
	"fmt"
	"strings"

	"github.com/keegancsmith/sqlf"
	"github.com/opentracing/opentracing-go/log"

	"github.com/sourcegraph/sourcegraph/internal/observation"
	"github.com/sourcegraph/sourcegraph/lib/codeintel/precise"
)

// MonikersByPosition returns all monikers attached ranges containing the given position. If multiple
// ranges contain the position, then this method will return multiple sets of monikers. Each slice
// of monikers are attached to a single range. The order of the output slice is "outside-in", so that
// the range attached to earlier monikers enclose the range attached to later monikers.
func (s *Store) MonikersByPosition(ctx context.Context, bundleID int, path string, line, character int) (_ [][]precise.MonikerData, err error) {
	ctx, trace, endObservation := s.operations.monikersByPosition.With(ctx, &err, observation.Args{LogFields: []log.Field{
		log.Int("bundleID", bundleID),
		log.String("path", path),
		log.Int("line", line),
		log.Int("character", character),
	}})
	defer endObservation(1, observation.Args{})

	documentData, exists, err := s.scanFirstDocumentData(s.Store.Query(ctx, sqlf.Sprintf(monikersDocumentQuery, bundleID, path)))
	if err != nil || !exists {
		return nil, err
	}

	trace.Log(log.Int("numRanges", len(documentData.Document.Ranges)))
	ranges := precise.FindRanges(documentData.Document.Ranges, line, character)
	trace.Log(log.Int("numIntersectingRanges", len(ranges)))

	monikerData := make([][]precise.MonikerData, 0, len(ranges))
	for _, r := range ranges {
		batch := make([]precise.MonikerData, 0, len(r.MonikerIDs))
		for _, monikerID := range r.MonikerIDs {
			if moniker, exists := documentData.Document.Monikers[monikerID]; exists {
				batch = append(batch, moniker)
			}
		}
		trace.Log(log.Int("numMonikersForRange", len(batch)))

		monikerData = append(monikerData, batch)
	}
	trace.Log(log.Int("numMonikers", len(monikerData)))

	return monikerData, nil
}

const monikersDocumentQuery = `
-- source: enterprise/internal/codeintel/stores/lsifstore/monikers.go:MonikersByPosition
SELECT
	dump_id,
	path,
	data,
	ranges,
	NULL AS hovers,
	monikers,
	NULL AS packages,
	NULL AS diagnostics
FROM
	lsif_data_documents
WHERE
	dump_id = %s AND
	path = %s
LIMIT 1
`

// BulkMonikerResults returns the locations (within one of the given uploads) with an attached moniker
// whose scheme+identifier matches one of the given monikers. This method also returns the size of the
// complete result set to aid in pagination.
func (s *Store) BulkMonikerResults(ctx context.Context, tableName string, uploadIDs []int, monikers []precise.MonikerData, limit, offset int) (_ []Location, _ int, err error) {
	ctx, trace, endObservation := s.operations.bulkMonikerResults.With(ctx, &err, observation.Args{LogFields: []log.Field{
		log.String("tableName", tableName),
		log.Int("numUploadIDs", len(uploadIDs)),
		log.String("uploadIDs", intsToString(uploadIDs)),
		log.Int("numMonikers", len(monikers)),
		log.String("monikers", monikersToString(monikers)),
		log.Int("limit", limit),
		log.Int("offset", offset),
	}})
	defer endObservation(1, observation.Args{})

	if len(uploadIDs) == 0 || len(monikers) == 0 {
		return nil, 0, nil
	}

	idQueries := make([]*sqlf.Query, 0, len(uploadIDs))
	for _, id := range uploadIDs {
		idQueries = append(idQueries, sqlf.Sprintf("%s", id))
	}

	monikerQueries := make([]*sqlf.Query, 0, len(monikers))
	for _, arg := range monikers {
		monikerQueries = append(monikerQueries, sqlf.Sprintf("(%s, %s)", arg.Scheme, arg.Identifier))
	}

	locationData, err := s.scanQualifiedMonikerLocations(s.Store.Query(ctx, sqlf.Sprintf(
		bulkMonikerResultsQuery,
		sqlf.Sprintf(fmt.Sprintf("lsif_data_%s", tableName)),
		sqlf.Join(idQueries, ", "),
		sqlf.Join(monikerQueries, ", "),
	)))
	if err != nil {
		return nil, 0, err
	}

	totalCount := 0
	for _, monikerLocations := range locationData {
		totalCount += len(monikerLocations.Locations)
	}
	trace.Log(
		log.Int("numDumps", len(locationData)),
		log.Int("totalCount", totalCount),
	)

	max := totalCount
	if totalCount > limit {
		max = limit
	}

	locations := make([]Location, 0, max)
outer:
	for _, monikerLocations := range locationData {
		for _, row := range monikerLocations.Locations {
			offset--
			if offset >= 0 {
				continue
			}

			locations = append(locations, Location{
				DumpID: monikerLocations.DumpID,
				Path:   row.URI,
				Range:  newRange(row.StartLine, row.StartCharacter, row.EndLine, row.EndCharacter),
			})

			if len(locations) >= limit {
				break outer
			}
		}
	}
	trace.Log(log.Int("numLocations", len(locations)))

	return locations, totalCount, nil
}

const bulkMonikerResultsQuery = `
-- source: enterprise/internal/codeintel/stores/lsifstore/monikers.go:BulkMonikerResults
SELECT dump_id, scheme, identifier, data FROM %s WHERE dump_id IN (%s) AND (scheme, identifier) IN (%s) ORDER BY (dump_id, scheme, identifier)
`

func monikersToString(vs []precise.MonikerData) string {
	strs := make([]string, 0, len(vs))
	for _, v := range vs {
		strs = append(strs, fmt.Sprintf("%s:%s:%s", v.Kind, v.Scheme, v.Identifier))
	}

	return strings.Join(strs, ", ")
}
