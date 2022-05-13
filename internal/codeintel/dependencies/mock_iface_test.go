// Code generated by go-mockgen 1.2.0; DO NOT EDIT.

package dependencies

import (
	"context"
	"sync"

	api "github.com/sourcegraph/sourcegraph/internal/api"
	store "github.com/sourcegraph/sourcegraph/internal/codeintel/dependencies/internal/store"
	shared "github.com/sourcegraph/sourcegraph/internal/codeintel/dependencies/shared"
	reposource "github.com/sourcegraph/sourcegraph/internal/conf/reposource"
	gitdomain "github.com/sourcegraph/sourcegraph/internal/gitserver/gitdomain"
)

// MockLockfilesService is a mock implementation of the LockfilesService
// interface (from the package
// github.com/sourcegraph/sourcegraph/internal/codeintel/dependencies) used
// for unit testing.
type MockLockfilesService struct {
	// ListDependenciesFunc is an instance of a mock function object
	// controlling the behavior of the method ListDependencies.
	ListDependenciesFunc *LockfilesServiceListDependenciesFunc
}

// NewMockLockfilesService creates a new mock of the LockfilesService
// interface. All methods return zero values for all results, unless
// overwritten.
func NewMockLockfilesService() *MockLockfilesService {
	return &MockLockfilesService{
		ListDependenciesFunc: &LockfilesServiceListDependenciesFunc{
			defaultHook: func(context.Context, api.RepoName, string) (r0 []reposource.PackageDependency, r1 error) {
				return
			},
		},
	}
}

// NewStrictMockLockfilesService creates a new mock of the LockfilesService
// interface. All methods panic on invocation, unless overwritten.
func NewStrictMockLockfilesService() *MockLockfilesService {
	return &MockLockfilesService{
		ListDependenciesFunc: &LockfilesServiceListDependenciesFunc{
			defaultHook: func(context.Context, api.RepoName, string) ([]reposource.PackageDependency, error) {
				panic("unexpected invocation of MockLockfilesService.ListDependencies")
			},
		},
	}
}

// NewMockLockfilesServiceFrom creates a new mock of the
// MockLockfilesService interface. All methods delegate to the given
// implementation, unless overwritten.
func NewMockLockfilesServiceFrom(i LockfilesService) *MockLockfilesService {
	return &MockLockfilesService{
		ListDependenciesFunc: &LockfilesServiceListDependenciesFunc{
			defaultHook: i.ListDependencies,
		},
	}
}

// LockfilesServiceListDependenciesFunc describes the behavior when the
// ListDependencies method of the parent MockLockfilesService instance is
// invoked.
type LockfilesServiceListDependenciesFunc struct {
	defaultHook func(context.Context, api.RepoName, string) ([]reposource.PackageDependency, error)
	hooks       []func(context.Context, api.RepoName, string) ([]reposource.PackageDependency, error)
	history     []LockfilesServiceListDependenciesFuncCall
	mutex       sync.Mutex
}

// ListDependencies delegates to the next hook function in the queue and
// stores the parameter and result values of this invocation.
func (m *MockLockfilesService) ListDependencies(v0 context.Context, v1 api.RepoName, v2 string) ([]reposource.PackageDependency, error) {
	r0, r1 := m.ListDependenciesFunc.nextHook()(v0, v1, v2)
	m.ListDependenciesFunc.appendCall(LockfilesServiceListDependenciesFuncCall{v0, v1, v2, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the ListDependencies
// method of the parent MockLockfilesService instance is invoked and the
// hook queue is empty.
func (f *LockfilesServiceListDependenciesFunc) SetDefaultHook(hook func(context.Context, api.RepoName, string) ([]reposource.PackageDependency, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// ListDependencies method of the parent MockLockfilesService instance
// invokes the hook at the front of the queue and discards it. After the
// queue is empty, the default hook function is invoked for any future
// action.
func (f *LockfilesServiceListDependenciesFunc) PushHook(hook func(context.Context, api.RepoName, string) ([]reposource.PackageDependency, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *LockfilesServiceListDependenciesFunc) SetDefaultReturn(r0 []reposource.PackageDependency, r1 error) {
	f.SetDefaultHook(func(context.Context, api.RepoName, string) ([]reposource.PackageDependency, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *LockfilesServiceListDependenciesFunc) PushReturn(r0 []reposource.PackageDependency, r1 error) {
	f.PushHook(func(context.Context, api.RepoName, string) ([]reposource.PackageDependency, error) {
		return r0, r1
	})
}

func (f *LockfilesServiceListDependenciesFunc) nextHook() func(context.Context, api.RepoName, string) ([]reposource.PackageDependency, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *LockfilesServiceListDependenciesFunc) appendCall(r0 LockfilesServiceListDependenciesFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of LockfilesServiceListDependenciesFuncCall
// objects describing the invocations of this function.
func (f *LockfilesServiceListDependenciesFunc) History() []LockfilesServiceListDependenciesFuncCall {
	f.mutex.Lock()
	history := make([]LockfilesServiceListDependenciesFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// LockfilesServiceListDependenciesFuncCall is an object that describes an
// invocation of method ListDependencies on an instance of
// MockLockfilesService.
type LockfilesServiceListDependenciesFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 api.RepoName
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 string
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []reposource.PackageDependency
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c LockfilesServiceListDependenciesFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c LockfilesServiceListDependenciesFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// MockStore is a mock implementation of the Store interface (from the
// package
// github.com/sourcegraph/sourcegraph/internal/codeintel/dependencies) used
// for unit testing.
type MockStore struct {
	// DeleteDependencyReposByIDFunc is an instance of a mock function
	// object controlling the behavior of the method
	// DeleteDependencyReposByID.
	DeleteDependencyReposByIDFunc *StoreDeleteDependencyReposByIDFunc
	// ListDependencyReposFunc is an instance of a mock function object
	// controlling the behavior of the method ListDependencyRepos.
	ListDependencyReposFunc *StoreListDependencyReposFunc
	// UpsertDependencyReposFunc is an instance of a mock function object
	// controlling the behavior of the method UpsertDependencyRepos.
	UpsertDependencyReposFunc *StoreUpsertDependencyReposFunc
}

// NewMockStore creates a new mock of the Store interface. All methods
// return zero values for all results, unless overwritten.
func NewMockStore() *MockStore {
	return &MockStore{
		DeleteDependencyReposByIDFunc: &StoreDeleteDependencyReposByIDFunc{
			defaultHook: func(context.Context, ...int) (r0 error) {
				return
			},
		},
		ListDependencyReposFunc: &StoreListDependencyReposFunc{
			defaultHook: func(context.Context, store.ListDependencyReposOpts) (r0 []shared.Repo, r1 error) {
				return
			},
		},
		UpsertDependencyReposFunc: &StoreUpsertDependencyReposFunc{
			defaultHook: func(context.Context, []shared.Repo) (r0 []shared.Repo, r1 error) {
				return
			},
		},
	}
}

// NewStrictMockStore creates a new mock of the Store interface. All methods
// panic on invocation, unless overwritten.
func NewStrictMockStore() *MockStore {
	return &MockStore{
		DeleteDependencyReposByIDFunc: &StoreDeleteDependencyReposByIDFunc{
			defaultHook: func(context.Context, ...int) error {
				panic("unexpected invocation of MockStore.DeleteDependencyReposByID")
			},
		},
		ListDependencyReposFunc: &StoreListDependencyReposFunc{
			defaultHook: func(context.Context, store.ListDependencyReposOpts) ([]shared.Repo, error) {
				panic("unexpected invocation of MockStore.ListDependencyRepos")
			},
		},
		UpsertDependencyReposFunc: &StoreUpsertDependencyReposFunc{
			defaultHook: func(context.Context, []shared.Repo) ([]shared.Repo, error) {
				panic("unexpected invocation of MockStore.UpsertDependencyRepos")
			},
		},
	}
}

// NewMockStoreFrom creates a new mock of the MockStore interface. All
// methods delegate to the given implementation, unless overwritten.
func NewMockStoreFrom(i Store) *MockStore {
	return &MockStore{
		DeleteDependencyReposByIDFunc: &StoreDeleteDependencyReposByIDFunc{
			defaultHook: i.DeleteDependencyReposByID,
		},
		ListDependencyReposFunc: &StoreListDependencyReposFunc{
			defaultHook: i.ListDependencyRepos,
		},
		UpsertDependencyReposFunc: &StoreUpsertDependencyReposFunc{
			defaultHook: i.UpsertDependencyRepos,
		},
	}
}

// StoreDeleteDependencyReposByIDFunc describes the behavior when the
// DeleteDependencyReposByID method of the parent MockStore instance is
// invoked.
type StoreDeleteDependencyReposByIDFunc struct {
	defaultHook func(context.Context, ...int) error
	hooks       []func(context.Context, ...int) error
	history     []StoreDeleteDependencyReposByIDFuncCall
	mutex       sync.Mutex
}

// DeleteDependencyReposByID delegates to the next hook function in the
// queue and stores the parameter and result values of this invocation.
func (m *MockStore) DeleteDependencyReposByID(v0 context.Context, v1 ...int) error {
	r0 := m.DeleteDependencyReposByIDFunc.nextHook()(v0, v1...)
	m.DeleteDependencyReposByIDFunc.appendCall(StoreDeleteDependencyReposByIDFuncCall{v0, v1, r0})
	return r0
}

// SetDefaultHook sets function that is called when the
// DeleteDependencyReposByID method of the parent MockStore instance is
// invoked and the hook queue is empty.
func (f *StoreDeleteDependencyReposByIDFunc) SetDefaultHook(hook func(context.Context, ...int) error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// DeleteDependencyReposByID method of the parent MockStore instance invokes
// the hook at the front of the queue and discards it. After the queue is
// empty, the default hook function is invoked for any future action.
func (f *StoreDeleteDependencyReposByIDFunc) PushHook(hook func(context.Context, ...int) error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *StoreDeleteDependencyReposByIDFunc) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func(context.Context, ...int) error {
		return r0
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *StoreDeleteDependencyReposByIDFunc) PushReturn(r0 error) {
	f.PushHook(func(context.Context, ...int) error {
		return r0
	})
}

func (f *StoreDeleteDependencyReposByIDFunc) nextHook() func(context.Context, ...int) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *StoreDeleteDependencyReposByIDFunc) appendCall(r0 StoreDeleteDependencyReposByIDFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of StoreDeleteDependencyReposByIDFuncCall
// objects describing the invocations of this function.
func (f *StoreDeleteDependencyReposByIDFunc) History() []StoreDeleteDependencyReposByIDFuncCall {
	f.mutex.Lock()
	history := make([]StoreDeleteDependencyReposByIDFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// StoreDeleteDependencyReposByIDFuncCall is an object that describes an
// invocation of method DeleteDependencyReposByID on an instance of
// MockStore.
type StoreDeleteDependencyReposByIDFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is a slice containing the values of the variadic arguments
	// passed to this method invocation.
	Arg1 []int
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation. The variadic slice argument is flattened in this array such
// that one positional argument and three variadic arguments would result in
// a slice of four, not two.
func (c StoreDeleteDependencyReposByIDFuncCall) Args() []interface{} {
	trailing := []interface{}{}
	for _, val := range c.Arg1 {
		trailing = append(trailing, val)
	}

	return append([]interface{}{c.Arg0}, trailing...)
}

// Results returns an interface slice containing the results of this
// invocation.
func (c StoreDeleteDependencyReposByIDFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// StoreListDependencyReposFunc describes the behavior when the
// ListDependencyRepos method of the parent MockStore instance is invoked.
type StoreListDependencyReposFunc struct {
	defaultHook func(context.Context, store.ListDependencyReposOpts) ([]shared.Repo, error)
	hooks       []func(context.Context, store.ListDependencyReposOpts) ([]shared.Repo, error)
	history     []StoreListDependencyReposFuncCall
	mutex       sync.Mutex
}

// ListDependencyRepos delegates to the next hook function in the queue and
// stores the parameter and result values of this invocation.
func (m *MockStore) ListDependencyRepos(v0 context.Context, v1 store.ListDependencyReposOpts) ([]shared.Repo, error) {
	r0, r1 := m.ListDependencyReposFunc.nextHook()(v0, v1)
	m.ListDependencyReposFunc.appendCall(StoreListDependencyReposFuncCall{v0, v1, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the ListDependencyRepos
// method of the parent MockStore instance is invoked and the hook queue is
// empty.
func (f *StoreListDependencyReposFunc) SetDefaultHook(hook func(context.Context, store.ListDependencyReposOpts) ([]shared.Repo, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// ListDependencyRepos method of the parent MockStore instance invokes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *StoreListDependencyReposFunc) PushHook(hook func(context.Context, store.ListDependencyReposOpts) ([]shared.Repo, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *StoreListDependencyReposFunc) SetDefaultReturn(r0 []shared.Repo, r1 error) {
	f.SetDefaultHook(func(context.Context, store.ListDependencyReposOpts) ([]shared.Repo, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *StoreListDependencyReposFunc) PushReturn(r0 []shared.Repo, r1 error) {
	f.PushHook(func(context.Context, store.ListDependencyReposOpts) ([]shared.Repo, error) {
		return r0, r1
	})
}

func (f *StoreListDependencyReposFunc) nextHook() func(context.Context, store.ListDependencyReposOpts) ([]shared.Repo, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *StoreListDependencyReposFunc) appendCall(r0 StoreListDependencyReposFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of StoreListDependencyReposFuncCall objects
// describing the invocations of this function.
func (f *StoreListDependencyReposFunc) History() []StoreListDependencyReposFuncCall {
	f.mutex.Lock()
	history := make([]StoreListDependencyReposFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// StoreListDependencyReposFuncCall is an object that describes an
// invocation of method ListDependencyRepos on an instance of MockStore.
type StoreListDependencyReposFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 store.ListDependencyReposOpts
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []shared.Repo
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c StoreListDependencyReposFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c StoreListDependencyReposFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// StoreUpsertDependencyReposFunc describes the behavior when the
// UpsertDependencyRepos method of the parent MockStore instance is invoked.
type StoreUpsertDependencyReposFunc struct {
	defaultHook func(context.Context, []shared.Repo) ([]shared.Repo, error)
	hooks       []func(context.Context, []shared.Repo) ([]shared.Repo, error)
	history     []StoreUpsertDependencyReposFuncCall
	mutex       sync.Mutex
}

// UpsertDependencyRepos delegates to the next hook function in the queue
// and stores the parameter and result values of this invocation.
func (m *MockStore) UpsertDependencyRepos(v0 context.Context, v1 []shared.Repo) ([]shared.Repo, error) {
	r0, r1 := m.UpsertDependencyReposFunc.nextHook()(v0, v1)
	m.UpsertDependencyReposFunc.appendCall(StoreUpsertDependencyReposFuncCall{v0, v1, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the
// UpsertDependencyRepos method of the parent MockStore instance is invoked
// and the hook queue is empty.
func (f *StoreUpsertDependencyReposFunc) SetDefaultHook(hook func(context.Context, []shared.Repo) ([]shared.Repo, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// UpsertDependencyRepos method of the parent MockStore instance invokes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *StoreUpsertDependencyReposFunc) PushHook(hook func(context.Context, []shared.Repo) ([]shared.Repo, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *StoreUpsertDependencyReposFunc) SetDefaultReturn(r0 []shared.Repo, r1 error) {
	f.SetDefaultHook(func(context.Context, []shared.Repo) ([]shared.Repo, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *StoreUpsertDependencyReposFunc) PushReturn(r0 []shared.Repo, r1 error) {
	f.PushHook(func(context.Context, []shared.Repo) ([]shared.Repo, error) {
		return r0, r1
	})
}

func (f *StoreUpsertDependencyReposFunc) nextHook() func(context.Context, []shared.Repo) ([]shared.Repo, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *StoreUpsertDependencyReposFunc) appendCall(r0 StoreUpsertDependencyReposFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of StoreUpsertDependencyReposFuncCall objects
// describing the invocations of this function.
func (f *StoreUpsertDependencyReposFunc) History() []StoreUpsertDependencyReposFuncCall {
	f.mutex.Lock()
	history := make([]StoreUpsertDependencyReposFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// StoreUpsertDependencyReposFuncCall is an object that describes an
// invocation of method UpsertDependencyRepos on an instance of MockStore.
type StoreUpsertDependencyReposFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 []shared.Repo
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []shared.Repo
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c StoreUpsertDependencyReposFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c StoreUpsertDependencyReposFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// MockSyncer is a mock implementation of the Syncer interface (from the
// package
// github.com/sourcegraph/sourcegraph/internal/codeintel/dependencies) used
// for unit testing.
type MockSyncer struct {
	// SyncFunc is an instance of a mock function object controlling the
	// behavior of the method Sync.
	SyncFunc *SyncerSyncFunc
}

// NewMockSyncer creates a new mock of the Syncer interface. All methods
// return zero values for all results, unless overwritten.
func NewMockSyncer() *MockSyncer {
	return &MockSyncer{
		SyncFunc: &SyncerSyncFunc{
			defaultHook: func(context.Context, api.RepoName) (r0 error) {
				return
			},
		},
	}
}

// NewStrictMockSyncer creates a new mock of the Syncer interface. All
// methods panic on invocation, unless overwritten.
func NewStrictMockSyncer() *MockSyncer {
	return &MockSyncer{
		SyncFunc: &SyncerSyncFunc{
			defaultHook: func(context.Context, api.RepoName) error {
				panic("unexpected invocation of MockSyncer.Sync")
			},
		},
	}
}

// NewMockSyncerFrom creates a new mock of the MockSyncer interface. All
// methods delegate to the given implementation, unless overwritten.
func NewMockSyncerFrom(i Syncer) *MockSyncer {
	return &MockSyncer{
		SyncFunc: &SyncerSyncFunc{
			defaultHook: i.Sync,
		},
	}
}

// SyncerSyncFunc describes the behavior when the Sync method of the parent
// MockSyncer instance is invoked.
type SyncerSyncFunc struct {
	defaultHook func(context.Context, api.RepoName) error
	hooks       []func(context.Context, api.RepoName) error
	history     []SyncerSyncFuncCall
	mutex       sync.Mutex
}

// Sync delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockSyncer) Sync(v0 context.Context, v1 api.RepoName) error {
	r0 := m.SyncFunc.nextHook()(v0, v1)
	m.SyncFunc.appendCall(SyncerSyncFuncCall{v0, v1, r0})
	return r0
}

// SetDefaultHook sets function that is called when the Sync method of the
// parent MockSyncer instance is invoked and the hook queue is empty.
func (f *SyncerSyncFunc) SetDefaultHook(hook func(context.Context, api.RepoName) error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Sync method of the parent MockSyncer instance invokes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *SyncerSyncFunc) PushHook(hook func(context.Context, api.RepoName) error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *SyncerSyncFunc) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func(context.Context, api.RepoName) error {
		return r0
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *SyncerSyncFunc) PushReturn(r0 error) {
	f.PushHook(func(context.Context, api.RepoName) error {
		return r0
	})
}

func (f *SyncerSyncFunc) nextHook() func(context.Context, api.RepoName) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *SyncerSyncFunc) appendCall(r0 SyncerSyncFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of SyncerSyncFuncCall objects describing the
// invocations of this function.
func (f *SyncerSyncFunc) History() []SyncerSyncFuncCall {
	f.mutex.Lock()
	history := make([]SyncerSyncFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// SyncerSyncFuncCall is an object that describes an invocation of method
// Sync on an instance of MockSyncer.
type SyncerSyncFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 api.RepoName
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c SyncerSyncFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c SyncerSyncFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// MockLocalGitService is a mock implementation of the localGitService
// interface (from the package
// github.com/sourcegraph/sourcegraph/internal/codeintel/dependencies) used
// for unit testing.
type MockLocalGitService struct {
	// GetCommitsFunc is an instance of a mock function object controlling
	// the behavior of the method GetCommits.
	GetCommitsFunc *LocalGitServiceGetCommitsFunc
}

// NewMockLocalGitService creates a new mock of the localGitService
// interface. All methods return zero values for all results, unless
// overwritten.
func NewMockLocalGitService() *MockLocalGitService {
	return &MockLocalGitService{
		GetCommitsFunc: &LocalGitServiceGetCommitsFunc{
			defaultHook: func(context.Context, []api.RepoCommit, bool) (r0 []*gitdomain.Commit, r1 error) {
				return
			},
		},
	}
}

// NewStrictMockLocalGitService creates a new mock of the localGitService
// interface. All methods panic on invocation, unless overwritten.
func NewStrictMockLocalGitService() *MockLocalGitService {
	return &MockLocalGitService{
		GetCommitsFunc: &LocalGitServiceGetCommitsFunc{
			defaultHook: func(context.Context, []api.RepoCommit, bool) ([]*gitdomain.Commit, error) {
				panic("unexpected invocation of MockLocalGitService.GetCommits")
			},
		},
	}
}

// surrogateMockLocalGitService is a copy of the localGitService interface
// (from the package
// github.com/sourcegraph/sourcegraph/internal/codeintel/dependencies). It
// is redefined here as it is unexported in the source package.
type surrogateMockLocalGitService interface {
	GetCommits(context.Context, []api.RepoCommit, bool) ([]*gitdomain.Commit, error)
}

// NewMockLocalGitServiceFrom creates a new mock of the MockLocalGitService
// interface. All methods delegate to the given implementation, unless
// overwritten.
func NewMockLocalGitServiceFrom(i surrogateMockLocalGitService) *MockLocalGitService {
	return &MockLocalGitService{
		GetCommitsFunc: &LocalGitServiceGetCommitsFunc{
			defaultHook: i.GetCommits,
		},
	}
}

// LocalGitServiceGetCommitsFunc describes the behavior when the GetCommits
// method of the parent MockLocalGitService instance is invoked.
type LocalGitServiceGetCommitsFunc struct {
	defaultHook func(context.Context, []api.RepoCommit, bool) ([]*gitdomain.Commit, error)
	hooks       []func(context.Context, []api.RepoCommit, bool) ([]*gitdomain.Commit, error)
	history     []LocalGitServiceGetCommitsFuncCall
	mutex       sync.Mutex
}

// GetCommits delegates to the next hook function in the queue and stores
// the parameter and result values of this invocation.
func (m *MockLocalGitService) GetCommits(v0 context.Context, v1 []api.RepoCommit, v2 bool) ([]*gitdomain.Commit, error) {
	r0, r1 := m.GetCommitsFunc.nextHook()(v0, v1, v2)
	m.GetCommitsFunc.appendCall(LocalGitServiceGetCommitsFuncCall{v0, v1, v2, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the GetCommits method of
// the parent MockLocalGitService instance is invoked and the hook queue is
// empty.
func (f *LocalGitServiceGetCommitsFunc) SetDefaultHook(hook func(context.Context, []api.RepoCommit, bool) ([]*gitdomain.Commit, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// GetCommits method of the parent MockLocalGitService instance invokes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *LocalGitServiceGetCommitsFunc) PushHook(hook func(context.Context, []api.RepoCommit, bool) ([]*gitdomain.Commit, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *LocalGitServiceGetCommitsFunc) SetDefaultReturn(r0 []*gitdomain.Commit, r1 error) {
	f.SetDefaultHook(func(context.Context, []api.RepoCommit, bool) ([]*gitdomain.Commit, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *LocalGitServiceGetCommitsFunc) PushReturn(r0 []*gitdomain.Commit, r1 error) {
	f.PushHook(func(context.Context, []api.RepoCommit, bool) ([]*gitdomain.Commit, error) {
		return r0, r1
	})
}

func (f *LocalGitServiceGetCommitsFunc) nextHook() func(context.Context, []api.RepoCommit, bool) ([]*gitdomain.Commit, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *LocalGitServiceGetCommitsFunc) appendCall(r0 LocalGitServiceGetCommitsFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of LocalGitServiceGetCommitsFuncCall objects
// describing the invocations of this function.
func (f *LocalGitServiceGetCommitsFunc) History() []LocalGitServiceGetCommitsFuncCall {
	f.mutex.Lock()
	history := make([]LocalGitServiceGetCommitsFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// LocalGitServiceGetCommitsFuncCall is an object that describes an
// invocation of method GetCommits on an instance of MockLocalGitService.
type LocalGitServiceGetCommitsFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 []api.RepoCommit
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 bool
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []*gitdomain.Commit
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c LocalGitServiceGetCommitsFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c LocalGitServiceGetCommitsFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}
