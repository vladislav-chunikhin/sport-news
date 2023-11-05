// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// FeedFetcher is an autogenerated mock type for the FeedFetcher type
type FeedFetcher struct {
	mock.Mock
}

type FeedFetcher_Expecter struct {
	mock *mock.Mock
}

func (_m *FeedFetcher) EXPECT() *FeedFetcher_Expecter {
	return &FeedFetcher_Expecter{mock: &_m.Mock}
}

// Fetch provides a mock function with given fields: ctx
func (_m *FeedFetcher) Fetch(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FeedFetcher_Fetch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Fetch'
type FeedFetcher_Fetch_Call struct {
	*mock.Call
}

// Fetch is a helper method to define mock.On call
//   - ctx context.Context
func (_e *FeedFetcher_Expecter) Fetch(ctx interface{}) *FeedFetcher_Fetch_Call {
	return &FeedFetcher_Fetch_Call{Call: _e.mock.On("Fetch", ctx)}
}

func (_c *FeedFetcher_Fetch_Call) Run(run func(ctx context.Context)) *FeedFetcher_Fetch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *FeedFetcher_Fetch_Call) Return(_a0 error) *FeedFetcher_Fetch_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *FeedFetcher_Fetch_Call) RunAndReturn(run func(context.Context) error) *FeedFetcher_Fetch_Call {
	_c.Call.Return(run)
	return _c
}

// NewFeedFetcher creates a new instance of FeedFetcher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFeedFetcher(t interface {
	mock.TestingT
	Cleanup(func())
}) *FeedFetcher {
	mock := &FeedFetcher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
