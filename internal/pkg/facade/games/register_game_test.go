package games

import (
	"errors"
	"testing"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestFacade_RegisterGame(t *testing.T) {
	t.Run("error while register game", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().RegisterGame(fx.ctx, &registrator.RegisterGameRequest{
			GameId: 1,
		}).Return(nil, errors.New("some error"))

		got, err := fx.facade.RegisterGame(fx.ctx, 1)
		assert.Equal(t, int32(registrator.RegisterGameStatus_REGISTER_GAME_STATUS_INVALID), got)
		assert.Error(t, err)
	})

	t.Run("error game not found while register game", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().RegisterGame(fx.ctx, &registrator.RegisterGameRequest{
			GameId: 1,
		}).Return(nil, status.New(codes.NotFound, "").Err())

		got, err := fx.facade.RegisterGame(fx.ctx, 1)
		assert.Equal(t, int32(registrator.RegisterGameStatus_REGISTER_GAME_STATUS_INVALID), got)
		assert.Error(t, err)
		assert.ErrorIs(t, err, model.ErrGameNotFound)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().RegisterGame(fx.ctx, &registrator.RegisterGameRequest{
			GameId: 1,
		}).Return(&registrator.RegisterGameResponse{
			Status: registrator.RegisterGameStatus_REGISTER_GAME_STATUS_OK,
		}, nil)

		got, err := fx.facade.RegisterGame(fx.ctx, 1)
		assert.Equal(t, int32(registrator.RegisterGameStatus_REGISTER_GAME_STATUS_OK), got)
		assert.NoError(t, err)
	})
}
