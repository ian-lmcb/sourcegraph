package store

import (
	"fmt"

	"github.com/sourcegraph/sourcegraph/internal/metrics"
	"github.com/sourcegraph/sourcegraph/internal/observation"
)

type operations struct {
	deleteDependencyReposByID  *observation.Operation
	listDependencyRepos        *observation.Operation
	lockfileDependencies       *observation.Operation
	upsertDependencyRepos      *observation.Operation
	upsertLockfileDependencies *observation.Operation
}

func newOperations(observationContext *observation.Context) *operations {
	metrics := metrics.NewREDMetrics(
		observationContext.Registerer,
		"codeintel_dependencies_store",
		metrics.WithLabels("op"),
		metrics.WithCountHelp("Total number of method invocations."),
	)

	op := func(name string) *observation.Operation {
		return observationContext.Operation(observation.Op{
			Name:              fmt.Sprintf("codeintel.dependencies.store.%s", name),
			MetricLabelValues: []string{name},
			Metrics:           metrics,
		})
	}

	return &operations{
		deleteDependencyReposByID:  op("DeleteDependencyReposByID"),
		listDependencyRepos:        op("ListDependencyRepos"),
		lockfileDependencies:       op("LockfileDependencies"),
		upsertDependencyRepos:      op("UpsertDependencyRepos"),
		upsertLockfileDependencies: op("UpsertLockfileDependencies"),
	}
}
