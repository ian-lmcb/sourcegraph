package dependencies

import (
	"context"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/semaphore"

	"github.com/sourcegraph/sourcegraph/internal/codeintel/dependencies/internal/lockfiles"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/dependencies/internal/store"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/dependencies/shared"
	"github.com/sourcegraph/sourcegraph/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/env"
	"github.com/sourcegraph/sourcegraph/internal/observation"
	"github.com/sourcegraph/sourcegraph/internal/trace"
	"github.com/sourcegraph/sourcegraph/lib/log"
)

var (
	svc     *Service
	svcOnce sync.Once
)

var (
	lockfilesSemaphoreWeight = env.MustGetInt("CODEINTEL_DEPENDENCIES_LOCKFILES_SEMAPHORE_WEIGHT", 64, "The maximum number of concurrent routines parsing lockfile contents.")
	syncerSemaphoreWeight    = env.MustGetInt("CODEINTEL_DEPENDENCIES_LOCKFILES_SYNCER_WEIGHT", 64, "The maximum number of concurrent routines actively syncing repositories.")
	enableUpserts            = env.Get("CODEINTEL_DEPENDENCIES_ENABLE_UPSERTS", "", "Disable writes of lockfile results to the database from the dependencies service.") != ""
)

// GetService creates or returns an already-initialized dependencies service. If the service is
// new, it will use the given database handle and git/syncer instances.
func GetService(db database.DB, gitService GitService, syncer Syncer) *Service {
	logger := log.Scoped("dependencies.service", "codeintel dependencies service")

	svcOnce.Do(func() {
		observationContext := &observation.Context{
			Logger:     logger,
			Tracer:     &trace.Tracer{Tracer: opentracing.GlobalTracer()},
			Registerer: prometheus.DefaultRegisterer,
		}

		var store Store = store.GetStore(db)
		if !enableUpserts {
			logger.Warn("Disabling dependencies.store.UpsertLockfileDependencies")
			store = &shim{store}
		}

		lockfilesService := lockfiles.GetService(gitService)
		lockfilesSemaphore := semaphore.NewWeighted(int64(lockfilesSemaphoreWeight))
		syncerSemaphore := semaphore.NewWeighted(int64(syncerSemaphoreWeight))

		svc = newService(
			store,
			gitService,
			lockfilesService,
			lockfilesSemaphore,
			syncer,
			syncerSemaphore,
			observationContext,
		)
	})

	return svc
}

// TestService creates a fresh dependencies service with the given database handle and git/syncer
// instances.
func TestService(db database.DB, gitService GitService, syncer Syncer) *Service {
	lockfilesService := lockfiles.GetService(gitService)
	lockfilesSemaphore := semaphore.NewWeighted(64)
	syncerSemaphore := semaphore.NewWeighted(64)

	return newService(
		store.GetStore(db),
		gitService,
		lockfilesService,
		lockfilesSemaphore,
		syncer,
		syncerSemaphore,
		&observation.TestContext,
	)
}

type shim struct{ Store }

func (s *shim) UpsertLockfileDependencies(ctx context.Context, repoName, commit string, deps []shared.PackageDependency) error {
	return nil
}
