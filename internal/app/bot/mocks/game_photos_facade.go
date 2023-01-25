// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/nikita5637/quiz-telegram/internal/pkg/model"
	mock "github.com/stretchr/testify/mock"
)

// GamePhotosFacade is an autogenerated mock type for the GamePhotosFacade type
type GamePhotosFacade struct {
	mock.Mock
}

type GamePhotosFacade_Expecter struct {
	mock *mock.Mock
}

func (_m *GamePhotosFacade) EXPECT() *GamePhotosFacade_Expecter {
	return &GamePhotosFacade_Expecter{mock: &_m.Mock}
}

// GetGamesWithPhotos provides a mock function with given fields: ctx, limit, offset
func (_m *GamePhotosFacade) GetGamesWithPhotos(ctx context.Context, limit uint32, offset uint32) ([]model.Game, uint32, error) {
	ret := _m.Called(ctx, limit, offset)

	var r0 []model.Game
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) []model.Game); ok {
		r0 = rf(ctx, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Game)
		}
	}

	var r1 uint32
	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32) uint32); ok {
		r1 = rf(ctx, limit, offset)
	} else {
		r1 = ret.Get(1).(uint32)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, uint32, uint32) error); ok {
		r2 = rf(ctx, limit, offset)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GamePhotosFacade_GetGamesWithPhotos_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetGamesWithPhotos'
type GamePhotosFacade_GetGamesWithPhotos_Call struct {
	*mock.Call
}

// GetGamesWithPhotos is a helper method to define mock.On call
//  - ctx context.Context
//  - limit uint32
//  - offset uint32
func (_e *GamePhotosFacade_Expecter) GetGamesWithPhotos(ctx interface{}, limit interface{}, offset interface{}) *GamePhotosFacade_GetGamesWithPhotos_Call {
	return &GamePhotosFacade_GetGamesWithPhotos_Call{Call: _e.mock.On("GetGamesWithPhotos", ctx, limit, offset)}
}

func (_c *GamePhotosFacade_GetGamesWithPhotos_Call) Run(run func(ctx context.Context, limit uint32, offset uint32)) *GamePhotosFacade_GetGamesWithPhotos_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint32), args[2].(uint32))
	})
	return _c
}

func (_c *GamePhotosFacade_GetGamesWithPhotos_Call) Return(_a0 []model.Game, _a1 uint32, _a2 error) *GamePhotosFacade_GetGamesWithPhotos_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

// GetPhotosByGameID provides a mock function with given fields: ctx, gameID
func (_m *GamePhotosFacade) GetPhotosByGameID(ctx context.Context, gameID int32) ([]string, error) {
	ret := _m.Called(ctx, gameID)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, int32) []string); ok {
		r0 = rf(ctx, gameID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int32) error); ok {
		r1 = rf(ctx, gameID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GamePhotosFacade_GetPhotosByGameID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPhotosByGameID'
type GamePhotosFacade_GetPhotosByGameID_Call struct {
	*mock.Call
}

// GetPhotosByGameID is a helper method to define mock.On call
//  - ctx context.Context
//  - gameID int32
func (_e *GamePhotosFacade_Expecter) GetPhotosByGameID(ctx interface{}, gameID interface{}) *GamePhotosFacade_GetPhotosByGameID_Call {
	return &GamePhotosFacade_GetPhotosByGameID_Call{Call: _e.mock.On("GetPhotosByGameID", ctx, gameID)}
}

func (_c *GamePhotosFacade_GetPhotosByGameID_Call) Run(run func(ctx context.Context, gameID int32)) *GamePhotosFacade_GetPhotosByGameID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int32))
	})
	return _c
}

func (_c *GamePhotosFacade_GetPhotosByGameID_Call) Return(_a0 []string, _a1 error) *GamePhotosFacade_GetPhotosByGameID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewGamePhotosFacade interface {
	mock.TestingT
	Cleanup(func())
}

// NewGamePhotosFacade creates a new instance of GamePhotosFacade. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGamePhotosFacade(t mockConstructorTestingTNewGamePhotosFacade) *GamePhotosFacade {
	mock := &GamePhotosFacade{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
