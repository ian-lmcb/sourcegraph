package lsifstore

import (
	"context"

	"github.com/keegancsmith/sqlf"
	"github.com/opentracing/opentracing-go/log"

	"github.com/sourcegraph/sourcegraph/internal/observation"
	"github.com/sourcegraph/sourcegraph/lib/codeintel/precise"
)

// Hover returns the hover text of the symbol at the given position.
func (s *Store) Hover(ctx context.Context, bundleID int, path string, line, character int) (_ string, _ Range, _ bool, err error) {
	ctx, trace, endObservation := s.operations.hover.With(ctx, &err, observation.Args{LogFields: []log.Field{
		log.Int("bundleID", bundleID),
		log.String("path", path),
		log.Int("line", line),
		log.Int("character", character),
	}})
	defer endObservation(1, observation.Args{})

	documentData, exists, err := s.scanFirstDocumentData(s.Store.Query(ctx, sqlf.Sprintf(hoverDocumentQuery, bundleID, path)))
	if err != nil || !exists {
		return "", Range{}, false, err
	}

	trace.Log(log.Int("numRanges", len(documentData.Document.Ranges)))
	ranges := precise.FindRanges(documentData.Document.Ranges, line, character)
	trace.Log(log.Int("numIntersectingRanges", len(ranges)))

	for _, r := range ranges {
		if text, ok := documentData.Document.HoverResults[r.HoverResultID]; ok {
			return text, newRange(r.StartLine, r.StartCharacter, r.EndLine, r.EndCharacter), true, nil
		}
	}

	return "", Range{}, false, nil
}

const hoverDocumentQuery = `
-- source: enterprise/internal/codeintel/stores/lsifstore/hover.go:Hover
SELECT
	dump_id,
	path,
	data,
	ranges,
	hovers,
	NULL AS monikers,
	NULL AS packages,
	NULL AS diagnostics
FROM
	lsif_data_documents
WHERE
	dump_id = %s AND
	path = %s
LIMIT 1
`
