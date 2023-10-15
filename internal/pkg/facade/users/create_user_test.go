package users

import (
	"errors"
	"testing"

	usermanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/user_manager"
	telegramutils "github.com/nikita5637/quiz-telegram/utils/telegram"
	"github.com/stretchr/testify/assert"
)

func TestFacade_CreateUser(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		fx := tearUp(t)

		ctx := telegramutils.NewContextWithClientID(fx.ctx, 0)

		fx.userManagerServiceClient.EXPECT().CreateUser(ctx, &usermanagerpb.CreateUserRequest{
			User: &usermanagerpb.User{
				Name:       "name",
				TelegramId: -100,
				State:      usermanagerpb.UserState_USER_STATE_WELCOME,
			},
		}).Return(nil, errors.New("some error"))

		got, err := fx.facade.CreateUser(fx.ctx, "name", -100, int32(usermanagerpb.UserState_USER_STATE_WELCOME))
		assert.Equal(t, int32(0), got)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		ctx := telegramutils.NewContextWithClientID(fx.ctx, 0)

		fx.userManagerServiceClient.EXPECT().CreateUser(ctx, &usermanagerpb.CreateUserRequest{
			User: &usermanagerpb.User{
				Name:       "name",
				TelegramId: -100,
				State:      usermanagerpb.UserState_USER_STATE_WELCOME,
			},
		}).Return(&usermanagerpb.User{
			Id:         1,
			TelegramId: -100,
			State:      usermanagerpb.UserState_USER_STATE_WELCOME,
		}, nil)

		got, err := fx.facade.CreateUser(fx.ctx, "name", -100, int32(usermanagerpb.UserState_USER_STATE_WELCOME))
		assert.Equal(t, int32(1), got)
		assert.NoError(t, err)
	})
}
