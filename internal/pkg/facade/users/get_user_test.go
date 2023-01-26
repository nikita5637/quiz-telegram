package users

import (
	"errors"
	"testing"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestFacade_GetUserByID(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetUserByID(fx.ctx, &registrator.GetUserByIDRequest{
			Id: 1,
		}).Return(nil, errors.New("some error"))

		got, err := fx.facade.GetUserByID(fx.ctx, 1)
		assert.Equal(t, model.User{}, got)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetUserByID(fx.ctx, &registrator.GetUserByIDRequest{
			Id: 1,
		}).Return(&registrator.GetUserByIDResponse{
			User: &registrator.User{
				Id:    1,
				Email: "email",
				Name:  "name",
				Phone: "phone",
				State: registrator.UserState_USER_STATE_CHANGINE_NAME,
			},
		}, nil)

		got, err := fx.facade.GetUserByID(fx.ctx, 1)
		assert.Equal(t, model.User{
			ID:    1,
			Email: "email",
			Name:  "name",
			Phone: "phone",
			State: int32(registrator.UserState_USER_STATE_CHANGINE_NAME),
		}, got)
		assert.NoError(t, err)
	})
}

func TestFacade_GetUserByTelegramID(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetUserByTelegramID(fx.ctx, &registrator.GetUserByTelegramIDRequest{
			TelegramId: -100,
		}).Return(nil, errors.New("some error"))

		got, err := fx.facade.GetUserByTelegramID(fx.ctx, -100)
		assert.Equal(t, model.User{}, got)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetUserByTelegramID(fx.ctx, &registrator.GetUserByTelegramIDRequest{
			TelegramId: -100,
		}).Return(&registrator.GetUserByTelegramIDResponse{
			User: &registrator.User{
				Id:    1,
				Email: "email",
				Name:  "name",
				Phone: "phone",
				State: registrator.UserState_USER_STATE_CHANGINE_NAME,
			},
		}, nil)

		got, err := fx.facade.GetUserByTelegramID(fx.ctx, -100)
		assert.Equal(t, model.User{
			ID:    1,
			Email: "email",
			Name:  "name",
			Phone: "phone",
			State: int32(registrator.UserState_USER_STATE_CHANGINE_NAME),
		}, got)
		assert.NoError(t, err)
	})
}
