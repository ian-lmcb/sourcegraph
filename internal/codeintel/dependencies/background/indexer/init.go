package indexer

import (
	"context"

	"github.com/sourcegraph/sourcegraph/internal/codeintel/dependencies"
	livedependencies "github.com/sourcegraph/sourcegraph/internal/codeintel/dependencies/live"
	"github.com/sourcegraph/sourcegraph/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/goroutine"
)

func NewIndexer(db database.DB, syncer dependencies.Syncer) goroutine.BackgroundRoutine {
	return goroutine.NewPeriodicGoroutine(context.Background(), ConfigInst.Interval, &indexer{
		dependenciesSvc: livedependencies.GetService(db, syncer),
	})
}
