package users

import (
	"errors"
	"testing"

	usermanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/user_manager"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestFacade_GetUserByID(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		fx := tearUp(t)

		fx.userManagerServiceClient.EXPECT().GetUser(fx.ctx, &usermanagerpb.GetUserRequest{
			Id: 1,
		}).Return(nil, errors.New("some error"))

		got, err := fx.facade.GetUserByID(fx.ctx, 1)
		assert.Equal(t, model.User{}, got)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.userManagerServiceClient.EXPECT().GetUser(fx.ctx, &usermanagerpb.GetUserRequest{
			Id: 1,
		}).Return(&usermanagerpb.User{
			Id: 1,
			Email: &wrapperspb.StringValue{
				Value: "email",
			},
			Name: "name",
			Phone: &wrapperspb.StringValue{
				Value: "phone",
			},
			State: usermanagerpb.UserState_USER_STATE_CHANGING_NAME,
		}, nil)

		got, err := fx.facade.GetUserByID(fx.ctx, 1)
		assert.Equal(t, model.User{
			ID:    1,
			Email: "email",
			Name:  "name",
			Phone: "phone",
			State: int32(usermanagerpb.UserState_USER_STATE_CHANGING_NAME),
		}, got)
		assert.NoError(t, err)
	})
}

func TestFacade_GetUserByTelegramID(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		fx := tearUp(t)

		fx.userManagerServiceClient.EXPECT().GetUserByTelegramID(fx.ctx, &usermanagerpb.GetUserByTelegramIDRequest{
			TelegramId: -100,
		}).Return(nil, errors.New("some error"))

		got, err := fx.facade.GetUserByTelegramID(fx.ctx, -100)
		assert.Equal(t, model.User{}, got)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.userManagerServiceClient.EXPECT().GetUserByTelegramID(fx.ctx, &usermanagerpb.GetUserByTelegramIDRequest{
			TelegramId: -100,
		}).Return(&usermanagerpb.User{
			Id: 1,
			Email: &wrapperspb.StringValue{
				Value: "email",
			},
			Name: "name",
			Phone: &wrapperspb.StringValue{
				Value: "phone",
			},
			State: usermanagerpb.UserState_USER_STATE_CHANGING_NAME,
		}, nil)

		got, err := fx.facade.GetUserByTelegramID(fx.ctx, -100)
		assert.Equal(t, model.User{
			ID:    1,
			Email: "email",
			Name:  "name",
			Phone: "phone",
			State: int32(usermanagerpb.UserState_USER_STATE_CHANGING_NAME),
		}, got)
		assert.NoError(t, err)
	})
}
