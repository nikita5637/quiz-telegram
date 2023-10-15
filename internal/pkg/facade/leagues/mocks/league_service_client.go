// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	league "github.com/nikita5637/quiz-registrator-api/pkg/pb/league"

	mock "github.com/stretchr/testify/mock"
)

// LeagueServiceClient is an autogenerated mock type for the LeagueServiceClient type
type LeagueServiceClient struct {
	mock.Mock
}

type LeagueServiceClient_Expecter struct {
	mock *mock.Mock
}

func (_m *LeagueServiceClient) EXPECT() *LeagueServiceClient_Expecter {
	return &LeagueServiceClient_Expecter{mock: &_m.Mock}
}

// GetLeague provides a mock function with given fields: ctx, in, opts
func (_m *LeagueServiceClient) GetLeague(ctx context.Context, in *league.GetLeagueRequest, opts ...grpc.CallOption) (*league.League, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *league.League
	if rf, ok := ret.Get(0).(func(context.Context, *league.GetLeagueRequest, ...grpc.CallOption) *league.League); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*league.League)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *league.GetLeagueRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LeagueServiceClient_GetLeague_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLeague'
type LeagueServiceClient_GetLeague_Call struct {
	*mock.Call
}

// GetLeague is a helper method to define mock.On call
//  - ctx context.Context
//  - in *league.GetLeagueRequest
//  - opts ...grpc.CallOption
func (_e *LeagueServiceClient_Expecter) GetLeague(ctx interface{}, in interface{}, opts ...interface{}) *LeagueServiceClient_GetLeague_Call {
	return &LeagueServiceClient_GetLeague_Call{Call: _e.mock.On("GetLeague",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *LeagueServiceClient_GetLeague_Call) Run(run func(ctx context.Context, in *league.GetLeagueRequest, opts ...grpc.CallOption)) *LeagueServiceClient_GetLeague_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*league.GetLeagueRequest), variadicArgs...)
	})
	return _c
}

func (_c *LeagueServiceClient_GetLeague_Call) Return(_a0 *league.League, _a1 error) *LeagueServiceClient_GetLeague_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewLeagueServiceClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewLeagueServiceClient creates a new instance of LeagueServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewLeagueServiceClient(t mockConstructorTestingTNewLeagueServiceClient) *LeagueServiceClient {
	mock := &LeagueServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
