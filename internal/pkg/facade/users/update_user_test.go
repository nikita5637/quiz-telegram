package users

import (
	"errors"
	"testing"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/stretchr/testify/assert"
)

func TestFacade_UpdateUserEmail(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().UpdateUserEmail(fx.ctx, &registrator.UpdateUserEmailRequest{
			UserId: 1,
			Email:  "email",
		}).Return(nil, errors.New("some error"))

		err := fx.facade.UpdateUserEmail(fx.ctx, 1, "email")
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().UpdateUserEmail(fx.ctx, &registrator.UpdateUserEmailRequest{
			UserId: 1,
			Email:  "email",
		}).Return(nil, nil)

		err := fx.facade.UpdateUserEmail(fx.ctx, 1, "email")
		assert.NoError(t, err)
	})
}

func TestFacade_UpdateUserName(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().UpdateUserName(fx.ctx, &registrator.UpdateUserNameRequest{
			UserId: 1,
			Name:   "name",
		}).Return(nil, errors.New("some error"))

		err := fx.facade.UpdateUserName(fx.ctx, 1, "name")
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().UpdateUserName(fx.ctx, &registrator.UpdateUserNameRequest{
			UserId: 1,
			Name:   "name",
		}).Return(nil, nil)

		err := fx.facade.UpdateUserName(fx.ctx, 1, "name")
		assert.NoError(t, err)
	})
}

func TestFacade_UpdateUserPhone(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().UpdateUserPhone(fx.ctx, &registrator.UpdateUserPhoneRequest{
			UserId: 1,
			Phone:  "phone",
		}).Return(nil, errors.New("some error"))

		err := fx.facade.UpdateUserPhone(fx.ctx, 1, "phone")
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().UpdateUserPhone(fx.ctx, &registrator.UpdateUserPhoneRequest{
			UserId: 1,
			Phone:  "phone",
		}).Return(nil, nil)

		err := fx.facade.UpdateUserPhone(fx.ctx, 1, "phone")
		assert.NoError(t, err)
	})
}

func TestFacade_UpdateUserState(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().UpdateUserState(fx.ctx, &registrator.UpdateUserStateRequest{
			UserId: 1,
			State:  registrator.UserState_USER_STATE_CHANGING_NAME,
		}).Return(nil, errors.New("some error"))

		err := fx.facade.UpdateUserState(fx.ctx, 1, int32(registrator.UserState_USER_STATE_CHANGING_NAME))
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().UpdateUserState(fx.ctx, &registrator.UpdateUserStateRequest{
			UserId: 1,
			State:  registrator.UserState_USER_STATE_CHANGING_NAME,
		}).Return(nil, nil)

		err := fx.facade.UpdateUserState(fx.ctx, 1, int32(registrator.UserState_USER_STATE_CHANGING_NAME))
		assert.NoError(t, err)
	})
}
