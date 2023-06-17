// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	icsfilemanager "github.com/nikita5637/quiz-ics-manager-api/pkg/pb/ics_file_manager"

	mock "github.com/stretchr/testify/mock"
)

// ICSFileManagerAPIServiceClient is an autogenerated mock type for the ICSFileManagerAPIServiceClient type
type ICSFileManagerAPIServiceClient struct {
	mock.Mock
}

type ICSFileManagerAPIServiceClient_Expecter struct {
	mock *mock.Mock
}

func (_m *ICSFileManagerAPIServiceClient) EXPECT() *ICSFileManagerAPIServiceClient_Expecter {
	return &ICSFileManagerAPIServiceClient_Expecter{mock: &_m.Mock}
}

// GetICSFileByGameID provides a mock function with given fields: ctx, in, opts
func (_m *ICSFileManagerAPIServiceClient) GetICSFileByGameID(ctx context.Context, in *icsfilemanager.GetICSFileByGameIDRequest, opts ...grpc.CallOption) (*icsfilemanager.ICSFile, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *icsfilemanager.ICSFile
	if rf, ok := ret.Get(0).(func(context.Context, *icsfilemanager.GetICSFileByGameIDRequest, ...grpc.CallOption) *icsfilemanager.ICSFile); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*icsfilemanager.ICSFile)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *icsfilemanager.GetICSFileByGameIDRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ICSFileManagerAPIServiceClient_GetICSFileByGameID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetICSFileByGameID'
type ICSFileManagerAPIServiceClient_GetICSFileByGameID_Call struct {
	*mock.Call
}

// GetICSFileByGameID is a helper method to define mock.On call
//  - ctx context.Context
//  - in *icsfilemanager.GetICSFileByGameIDRequest
//  - opts ...grpc.CallOption
func (_e *ICSFileManagerAPIServiceClient_Expecter) GetICSFileByGameID(ctx interface{}, in interface{}, opts ...interface{}) *ICSFileManagerAPIServiceClient_GetICSFileByGameID_Call {
	return &ICSFileManagerAPIServiceClient_GetICSFileByGameID_Call{Call: _e.mock.On("GetICSFileByGameID",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *ICSFileManagerAPIServiceClient_GetICSFileByGameID_Call) Run(run func(ctx context.Context, in *icsfilemanager.GetICSFileByGameIDRequest, opts ...grpc.CallOption)) *ICSFileManagerAPIServiceClient_GetICSFileByGameID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*icsfilemanager.GetICSFileByGameIDRequest), variadicArgs...)
	})
	return _c
}

func (_c *ICSFileManagerAPIServiceClient_GetICSFileByGameID_Call) Return(_a0 *icsfilemanager.ICSFile, _a1 error) *ICSFileManagerAPIServiceClient_GetICSFileByGameID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewICSFileManagerAPIServiceClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewICSFileManagerAPIServiceClient creates a new instance of ICSFileManagerAPIServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewICSFileManagerAPIServiceClient(t mockConstructorTestingTNewICSFileManagerAPIServiceClient) *ICSFileManagerAPIServiceClient {
	mock := &ICSFileManagerAPIServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
