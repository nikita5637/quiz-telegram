// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	croupier "github.com/nikita5637/quiz-registrator-api/pkg/pb/croupier"
	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// CroupierServiceClient is an autogenerated mock type for the CroupierServiceClient type
type CroupierServiceClient struct {
	mock.Mock
}

type CroupierServiceClient_Expecter struct {
	mock *mock.Mock
}

func (_m *CroupierServiceClient) EXPECT() *CroupierServiceClient_Expecter {
	return &CroupierServiceClient_Expecter{mock: &_m.Mock}
}

// GetLotteryStatus provides a mock function with given fields: ctx, in, opts
func (_m *CroupierServiceClient) GetLotteryStatus(ctx context.Context, in *croupier.GetLotteryStatusRequest, opts ...grpc.CallOption) (*croupier.GetLotteryStatusResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *croupier.GetLotteryStatusResponse
	if rf, ok := ret.Get(0).(func(context.Context, *croupier.GetLotteryStatusRequest, ...grpc.CallOption) *croupier.GetLotteryStatusResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*croupier.GetLotteryStatusResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *croupier.GetLotteryStatusRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CroupierServiceClient_GetLotteryStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLotteryStatus'
type CroupierServiceClient_GetLotteryStatus_Call struct {
	*mock.Call
}

// GetLotteryStatus is a helper method to define mock.On call
//  - ctx context.Context
//  - in *croupier.GetLotteryStatusRequest
//  - opts ...grpc.CallOption
func (_e *CroupierServiceClient_Expecter) GetLotteryStatus(ctx interface{}, in interface{}, opts ...interface{}) *CroupierServiceClient_GetLotteryStatus_Call {
	return &CroupierServiceClient_GetLotteryStatus_Call{Call: _e.mock.On("GetLotteryStatus",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *CroupierServiceClient_GetLotteryStatus_Call) Run(run func(ctx context.Context, in *croupier.GetLotteryStatusRequest, opts ...grpc.CallOption)) *CroupierServiceClient_GetLotteryStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*croupier.GetLotteryStatusRequest), variadicArgs...)
	})
	return _c
}

func (_c *CroupierServiceClient_GetLotteryStatus_Call) Return(_a0 *croupier.GetLotteryStatusResponse, _a1 error) *CroupierServiceClient_GetLotteryStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// RegisterForLottery provides a mock function with given fields: ctx, in, opts
func (_m *CroupierServiceClient) RegisterForLottery(ctx context.Context, in *croupier.RegisterForLotteryRequest, opts ...grpc.CallOption) (*croupier.RegisterForLotteryResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *croupier.RegisterForLotteryResponse
	if rf, ok := ret.Get(0).(func(context.Context, *croupier.RegisterForLotteryRequest, ...grpc.CallOption) *croupier.RegisterForLotteryResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*croupier.RegisterForLotteryResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *croupier.RegisterForLotteryRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CroupierServiceClient_RegisterForLottery_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterForLottery'
type CroupierServiceClient_RegisterForLottery_Call struct {
	*mock.Call
}

// RegisterForLottery is a helper method to define mock.On call
//  - ctx context.Context
//  - in *croupier.RegisterForLotteryRequest
//  - opts ...grpc.CallOption
func (_e *CroupierServiceClient_Expecter) RegisterForLottery(ctx interface{}, in interface{}, opts ...interface{}) *CroupierServiceClient_RegisterForLottery_Call {
	return &CroupierServiceClient_RegisterForLottery_Call{Call: _e.mock.On("RegisterForLottery",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *CroupierServiceClient_RegisterForLottery_Call) Run(run func(ctx context.Context, in *croupier.RegisterForLotteryRequest, opts ...grpc.CallOption)) *CroupierServiceClient_RegisterForLottery_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*croupier.RegisterForLotteryRequest), variadicArgs...)
	})
	return _c
}

func (_c *CroupierServiceClient_RegisterForLottery_Call) Return(_a0 *croupier.RegisterForLotteryResponse, _a1 error) *CroupierServiceClient_RegisterForLottery_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewCroupierServiceClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewCroupierServiceClient creates a new instance of CroupierServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCroupierServiceClient(t mockConstructorTestingTNewCroupierServiceClient) *CroupierServiceClient {
	mock := &CroupierServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
