package indexer

import (
	"context"
	"fmt"

	"github.com/sourcegraph/sourcegraph/internal/api"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/dependencies"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/types"
	"github.com/sourcegraph/sourcegraph/internal/goroutine"
)

type indexer struct {
	dependenciesSvc *dependencies.Service
}

var _ goroutine.Handler = &indexer{}
var _ goroutine.ErrorHandler = &indexer{}

func (i *indexer) Handle(ctx context.Context) error {
	// TODO - choose actual repos that don't have data here

	_, err := i.dependenciesSvc.Dependencies(ctx, map[api.RepoName]types.RevSpecSet{
		api.RepoName("github.com/sourcegraph/lsif-go"):     {api.RevSpec("HEAD"): struct{}{}},
		api.RepoName("github.com/sourcegraph/sourcegraph"): {api.RevSpec("HEAD"): struct{}{}},
	})

	return err
}

func (r *indexer) HandleError(err error) {
	fmt.Printf("Failed to index lockfiles: %v\n", err)
}
