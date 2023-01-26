package users

import (
	"errors"
	"testing"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/stretchr/testify/assert"
)

func TestFacade_CreateUser(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().CreateUser(fx.ctx, &registrator.CreateUserRequest{
			Name:       "name",
			TelegramId: -100,
			State:      registrator.UserState_USER_STATE_WELCOME,
		}).Return(nil, errors.New("some error"))

		got, err := fx.facade.CreateUser(fx.ctx, "name", -100, int32(registrator.UserState_USER_STATE_WELCOME))
		assert.Equal(t, int32(0), got)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().CreateUser(fx.ctx, &registrator.CreateUserRequest{
			Name:       "name",
			TelegramId: -100,
			State:      registrator.UserState_USER_STATE_WELCOME,
		}).Return(&registrator.CreateUserResponse{
			Id: 1,
		}, nil)

		got, err := fx.facade.CreateUser(fx.ctx, "name", -100, int32(registrator.UserState_USER_STATE_WELCOME))
		assert.Equal(t, int32(1), got)
		assert.NoError(t, err)
	})
}
