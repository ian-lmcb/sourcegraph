package resolvers

import (
	"context"
	"strconv"
	"strings"

	"github.com/sourcegraph/sourcegraph/cmd/frontend/graphqlbackend"
	"github.com/sourcegraph/sourcegraph/enterprise/internal/batches/search"
	"github.com/sourcegraph/sourcegraph/enterprise/internal/batches/store"
	btypes "github.com/sourcegraph/sourcegraph/enterprise/internal/batches/types"
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

type batchSpecWorkspaceResolutionResolver struct {
	store      *store.Store
	resolution *btypes.BatchSpecResolutionJob
}

var _ graphqlbackend.BatchSpecWorkspaceResolutionResolver = &batchSpecWorkspaceResolutionResolver{}

func (r *batchSpecWorkspaceResolutionResolver) State() string {
	return r.resolution.State.ToGraphQL()
}

func (r *batchSpecWorkspaceResolutionResolver) StartedAt() *graphqlbackend.DateTime {
	if r.resolution.StartedAt.IsZero() {
		return nil
	}
	return &graphqlbackend.DateTime{Time: r.resolution.StartedAt}
}

func (r *batchSpecWorkspaceResolutionResolver) FinishedAt() *graphqlbackend.DateTime {
	if r.resolution.FinishedAt.IsZero() {
		return nil
	}
	return &graphqlbackend.DateTime{Time: r.resolution.FinishedAt}
}

func (r *batchSpecWorkspaceResolutionResolver) FailureMessage() *string {
	return r.resolution.FailureMessage
}

func (r *batchSpecWorkspaceResolutionResolver) Workspaces(ctx context.Context, args *graphqlbackend.ListWorkspacesArgs) (graphqlbackend.BatchSpecWorkspaceConnectionResolver, error) {
	opts, err := workspacesListArgsToDBOpts(args)
	if err != nil {
		return nil, err
	}
	opts.BatchSpecID = r.resolution.BatchSpecID

	return &batchSpecWorkspaceConnectionResolver{store: r.store, opts: opts}, nil
}

func (r *batchSpecWorkspaceResolutionResolver) RecentlyCompleted(ctx context.Context, args *graphqlbackend.ListRecentlyCompletedWorkspacesArgs) graphqlbackend.BatchSpecWorkspaceConnectionResolver {
	// TODO(ssbc): not implemented
	return nil
}

func (r *batchSpecWorkspaceResolutionResolver) RecentlyErrored(ctx context.Context, args *graphqlbackend.ListRecentlyErroredWorkspacesArgs) graphqlbackend.BatchSpecWorkspaceConnectionResolver {
	// TODO(ssbc): not implemented
	return nil
}

func workspacesListArgsToDBOpts(args *graphqlbackend.ListWorkspacesArgs) (opts store.ListBatchSpecWorkspacesOpts, err error) {
	if err := validateFirstParamDefaults(args.First); err != nil {
		return opts, err
	}
	opts.Limit = int(args.First)
	if args.After != nil {
		id, err := strconv.Atoi(*args.After)
		if err != nil {
			return opts, err
		}
		opts.Cursor = int64(id)
	}

	if args.Search != nil {
		var err error
		opts.TextSearch, err = search.ParseTextSearch(*args.Search)
		if err != nil {
			return opts, errors.Wrap(err, "parsing search")
		}
	}

	if args.State != nil {
		if *args.State == "COMPLETED" {
			opts.OnlyCachedOrCompleted = true
		} else if *args.State == "PENDING" {
			opts.OnlyWithoutExecution = true
		} else {
			// Convert the GQL type into the DB type: we just need to lowercase it. Magic 🪄.
			opts.State = btypes.BatchSpecWorkspaceExecutionJobState(strings.ToLower(*args.State))
		}
	}

	return opts, nil
}
