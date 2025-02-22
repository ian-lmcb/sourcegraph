// Code generated by go-mockgen 1.2.0; DO NOT EDIT.

package resolvers

import (
	"context"
	"sync"

	lsifstore "github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/stores/lsifstore"
)

// MockPositionAdjuster is a mock implementation of the PositionAdjuster
// interface (from the package
// github.com/sourcegraph/sourcegraph/enterprise/cmd/frontend/internal/codeintel/resolvers)
// used for unit testing.
type MockPositionAdjuster struct {
	// AdjustPathFunc is an instance of a mock function object controlling
	// the behavior of the method AdjustPath.
	AdjustPathFunc *PositionAdjusterAdjustPathFunc
	// AdjustPositionFunc is an instance of a mock function object
	// controlling the behavior of the method AdjustPosition.
	AdjustPositionFunc *PositionAdjusterAdjustPositionFunc
	// AdjustRangeFunc is an instance of a mock function object controlling
	// the behavior of the method AdjustRange.
	AdjustRangeFunc *PositionAdjusterAdjustRangeFunc
}

// NewMockPositionAdjuster creates a new mock of the PositionAdjuster
// interface. All methods return zero values for all results, unless
// overwritten.
func NewMockPositionAdjuster() *MockPositionAdjuster {
	return &MockPositionAdjuster{
		AdjustPathFunc: &PositionAdjusterAdjustPathFunc{
			defaultHook: func(context.Context, string, string, bool) (r0 string, r1 bool, r2 error) {
				return
			},
		},
		AdjustPositionFunc: &PositionAdjusterAdjustPositionFunc{
			defaultHook: func(context.Context, string, string, lsifstore.Position, bool) (r0 string, r1 lsifstore.Position, r2 bool, r3 error) {
				return
			},
		},
		AdjustRangeFunc: &PositionAdjusterAdjustRangeFunc{
			defaultHook: func(context.Context, string, string, lsifstore.Range, bool) (r0 string, r1 lsifstore.Range, r2 bool, r3 error) {
				return
			},
		},
	}
}

// NewStrictMockPositionAdjuster creates a new mock of the PositionAdjuster
// interface. All methods panic on invocation, unless overwritten.
func NewStrictMockPositionAdjuster() *MockPositionAdjuster {
	return &MockPositionAdjuster{
		AdjustPathFunc: &PositionAdjusterAdjustPathFunc{
			defaultHook: func(context.Context, string, string, bool) (string, bool, error) {
				panic("unexpected invocation of MockPositionAdjuster.AdjustPath")
			},
		},
		AdjustPositionFunc: &PositionAdjusterAdjustPositionFunc{
			defaultHook: func(context.Context, string, string, lsifstore.Position, bool) (string, lsifstore.Position, bool, error) {
				panic("unexpected invocation of MockPositionAdjuster.AdjustPosition")
			},
		},
		AdjustRangeFunc: &PositionAdjusterAdjustRangeFunc{
			defaultHook: func(context.Context, string, string, lsifstore.Range, bool) (string, lsifstore.Range, bool, error) {
				panic("unexpected invocation of MockPositionAdjuster.AdjustRange")
			},
		},
	}
}

// NewMockPositionAdjusterFrom creates a new mock of the
// MockPositionAdjuster interface. All methods delegate to the given
// implementation, unless overwritten.
func NewMockPositionAdjusterFrom(i PositionAdjuster) *MockPositionAdjuster {
	return &MockPositionAdjuster{
		AdjustPathFunc: &PositionAdjusterAdjustPathFunc{
			defaultHook: i.AdjustPath,
		},
		AdjustPositionFunc: &PositionAdjusterAdjustPositionFunc{
			defaultHook: i.AdjustPosition,
		},
		AdjustRangeFunc: &PositionAdjusterAdjustRangeFunc{
			defaultHook: i.AdjustRange,
		},
	}
}

// PositionAdjusterAdjustPathFunc describes the behavior when the AdjustPath
// method of the parent MockPositionAdjuster instance is invoked.
type PositionAdjusterAdjustPathFunc struct {
	defaultHook func(context.Context, string, string, bool) (string, bool, error)
	hooks       []func(context.Context, string, string, bool) (string, bool, error)
	history     []PositionAdjusterAdjustPathFuncCall
	mutex       sync.Mutex
}

// AdjustPath delegates to the next hook function in the queue and stores
// the parameter and result values of this invocation.
func (m *MockPositionAdjuster) AdjustPath(v0 context.Context, v1 string, v2 string, v3 bool) (string, bool, error) {
	r0, r1, r2 := m.AdjustPathFunc.nextHook()(v0, v1, v2, v3)
	m.AdjustPathFunc.appendCall(PositionAdjusterAdjustPathFuncCall{v0, v1, v2, v3, r0, r1, r2})
	return r0, r1, r2
}

// SetDefaultHook sets function that is called when the AdjustPath method of
// the parent MockPositionAdjuster instance is invoked and the hook queue is
// empty.
func (f *PositionAdjusterAdjustPathFunc) SetDefaultHook(hook func(context.Context, string, string, bool) (string, bool, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// AdjustPath method of the parent MockPositionAdjuster instance invokes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *PositionAdjusterAdjustPathFunc) PushHook(hook func(context.Context, string, string, bool) (string, bool, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *PositionAdjusterAdjustPathFunc) SetDefaultReturn(r0 string, r1 bool, r2 error) {
	f.SetDefaultHook(func(context.Context, string, string, bool) (string, bool, error) {
		return r0, r1, r2
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *PositionAdjusterAdjustPathFunc) PushReturn(r0 string, r1 bool, r2 error) {
	f.PushHook(func(context.Context, string, string, bool) (string, bool, error) {
		return r0, r1, r2
	})
}

func (f *PositionAdjusterAdjustPathFunc) nextHook() func(context.Context, string, string, bool) (string, bool, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *PositionAdjusterAdjustPathFunc) appendCall(r0 PositionAdjusterAdjustPathFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of PositionAdjusterAdjustPathFuncCall objects
// describing the invocations of this function.
func (f *PositionAdjusterAdjustPathFunc) History() []PositionAdjusterAdjustPathFuncCall {
	f.mutex.Lock()
	history := make([]PositionAdjusterAdjustPathFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// PositionAdjusterAdjustPathFuncCall is an object that describes an
// invocation of method AdjustPath on an instance of MockPositionAdjuster.
type PositionAdjusterAdjustPathFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 string
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 string
	// Arg3 is the value of the 4th argument passed to this method
	// invocation.
	Arg3 bool
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 string
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 bool
	// Result2 is the value of the 3rd result returned from this method
	// invocation.
	Result2 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c PositionAdjusterAdjustPathFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2, c.Arg3}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c PositionAdjusterAdjustPathFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1, c.Result2}
}

// PositionAdjusterAdjustPositionFunc describes the behavior when the
// AdjustPosition method of the parent MockPositionAdjuster instance is
// invoked.
type PositionAdjusterAdjustPositionFunc struct {
	defaultHook func(context.Context, string, string, lsifstore.Position, bool) (string, lsifstore.Position, bool, error)
	hooks       []func(context.Context, string, string, lsifstore.Position, bool) (string, lsifstore.Position, bool, error)
	history     []PositionAdjusterAdjustPositionFuncCall
	mutex       sync.Mutex
}

// AdjustPosition delegates to the next hook function in the queue and
// stores the parameter and result values of this invocation.
func (m *MockPositionAdjuster) AdjustPosition(v0 context.Context, v1 string, v2 string, v3 lsifstore.Position, v4 bool) (string, lsifstore.Position, bool, error) {
	r0, r1, r2, r3 := m.AdjustPositionFunc.nextHook()(v0, v1, v2, v3, v4)
	m.AdjustPositionFunc.appendCall(PositionAdjusterAdjustPositionFuncCall{v0, v1, v2, v3, v4, r0, r1, r2, r3})
	return r0, r1, r2, r3
}

// SetDefaultHook sets function that is called when the AdjustPosition
// method of the parent MockPositionAdjuster instance is invoked and the
// hook queue is empty.
func (f *PositionAdjusterAdjustPositionFunc) SetDefaultHook(hook func(context.Context, string, string, lsifstore.Position, bool) (string, lsifstore.Position, bool, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// AdjustPosition method of the parent MockPositionAdjuster instance invokes
// the hook at the front of the queue and discards it. After the queue is
// empty, the default hook function is invoked for any future action.
func (f *PositionAdjusterAdjustPositionFunc) PushHook(hook func(context.Context, string, string, lsifstore.Position, bool) (string, lsifstore.Position, bool, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *PositionAdjusterAdjustPositionFunc) SetDefaultReturn(r0 string, r1 lsifstore.Position, r2 bool, r3 error) {
	f.SetDefaultHook(func(context.Context, string, string, lsifstore.Position, bool) (string, lsifstore.Position, bool, error) {
		return r0, r1, r2, r3
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *PositionAdjusterAdjustPositionFunc) PushReturn(r0 string, r1 lsifstore.Position, r2 bool, r3 error) {
	f.PushHook(func(context.Context, string, string, lsifstore.Position, bool) (string, lsifstore.Position, bool, error) {
		return r0, r1, r2, r3
	})
}

func (f *PositionAdjusterAdjustPositionFunc) nextHook() func(context.Context, string, string, lsifstore.Position, bool) (string, lsifstore.Position, bool, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *PositionAdjusterAdjustPositionFunc) appendCall(r0 PositionAdjusterAdjustPositionFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of PositionAdjusterAdjustPositionFuncCall
// objects describing the invocations of this function.
func (f *PositionAdjusterAdjustPositionFunc) History() []PositionAdjusterAdjustPositionFuncCall {
	f.mutex.Lock()
	history := make([]PositionAdjusterAdjustPositionFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// PositionAdjusterAdjustPositionFuncCall is an object that describes an
// invocation of method AdjustPosition on an instance of
// MockPositionAdjuster.
type PositionAdjusterAdjustPositionFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 string
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 string
	// Arg3 is the value of the 4th argument passed to this method
	// invocation.
	Arg3 lsifstore.Position
	// Arg4 is the value of the 5th argument passed to this method
	// invocation.
	Arg4 bool
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 string
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 lsifstore.Position
	// Result2 is the value of the 3rd result returned from this method
	// invocation.
	Result2 bool
	// Result3 is the value of the 4th result returned from this method
	// invocation.
	Result3 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c PositionAdjusterAdjustPositionFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2, c.Arg3, c.Arg4}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c PositionAdjusterAdjustPositionFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1, c.Result2, c.Result3}
}

// PositionAdjusterAdjustRangeFunc describes the behavior when the
// AdjustRange method of the parent MockPositionAdjuster instance is
// invoked.
type PositionAdjusterAdjustRangeFunc struct {
	defaultHook func(context.Context, string, string, lsifstore.Range, bool) (string, lsifstore.Range, bool, error)
	hooks       []func(context.Context, string, string, lsifstore.Range, bool) (string, lsifstore.Range, bool, error)
	history     []PositionAdjusterAdjustRangeFuncCall
	mutex       sync.Mutex
}

// AdjustRange delegates to the next hook function in the queue and stores
// the parameter and result values of this invocation.
func (m *MockPositionAdjuster) AdjustRange(v0 context.Context, v1 string, v2 string, v3 lsifstore.Range, v4 bool) (string, lsifstore.Range, bool, error) {
	r0, r1, r2, r3 := m.AdjustRangeFunc.nextHook()(v0, v1, v2, v3, v4)
	m.AdjustRangeFunc.appendCall(PositionAdjusterAdjustRangeFuncCall{v0, v1, v2, v3, v4, r0, r1, r2, r3})
	return r0, r1, r2, r3
}

// SetDefaultHook sets function that is called when the AdjustRange method
// of the parent MockPositionAdjuster instance is invoked and the hook queue
// is empty.
func (f *PositionAdjusterAdjustRangeFunc) SetDefaultHook(hook func(context.Context, string, string, lsifstore.Range, bool) (string, lsifstore.Range, bool, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// AdjustRange method of the parent MockPositionAdjuster instance invokes
// the hook at the front of the queue and discards it. After the queue is
// empty, the default hook function is invoked for any future action.
func (f *PositionAdjusterAdjustRangeFunc) PushHook(hook func(context.Context, string, string, lsifstore.Range, bool) (string, lsifstore.Range, bool, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *PositionAdjusterAdjustRangeFunc) SetDefaultReturn(r0 string, r1 lsifstore.Range, r2 bool, r3 error) {
	f.SetDefaultHook(func(context.Context, string, string, lsifstore.Range, bool) (string, lsifstore.Range, bool, error) {
		return r0, r1, r2, r3
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *PositionAdjusterAdjustRangeFunc) PushReturn(r0 string, r1 lsifstore.Range, r2 bool, r3 error) {
	f.PushHook(func(context.Context, string, string, lsifstore.Range, bool) (string, lsifstore.Range, bool, error) {
		return r0, r1, r2, r3
	})
}

func (f *PositionAdjusterAdjustRangeFunc) nextHook() func(context.Context, string, string, lsifstore.Range, bool) (string, lsifstore.Range, bool, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *PositionAdjusterAdjustRangeFunc) appendCall(r0 PositionAdjusterAdjustRangeFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of PositionAdjusterAdjustRangeFuncCall objects
// describing the invocations of this function.
func (f *PositionAdjusterAdjustRangeFunc) History() []PositionAdjusterAdjustRangeFuncCall {
	f.mutex.Lock()
	history := make([]PositionAdjusterAdjustRangeFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// PositionAdjusterAdjustRangeFuncCall is an object that describes an
// invocation of method AdjustRange on an instance of MockPositionAdjuster.
type PositionAdjusterAdjustRangeFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 string
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 string
	// Arg3 is the value of the 4th argument passed to this method
	// invocation.
	Arg3 lsifstore.Range
	// Arg4 is the value of the 5th argument passed to this method
	// invocation.
	Arg4 bool
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 string
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 lsifstore.Range
	// Result2 is the value of the 3rd result returned from this method
	// invocation.
	Result2 bool
	// Result3 is the value of the 4th result returned from this method
	// invocation.
	Result3 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c PositionAdjusterAdjustRangeFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2, c.Arg3, c.Arg4}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c PositionAdjusterAdjustRangeFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1, c.Result2, c.Result3}
}
