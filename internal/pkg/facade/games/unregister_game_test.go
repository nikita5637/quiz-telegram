package games

import (
	"errors"
	"testing"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestFacade_UnregisterGame(t *testing.T) {
	t.Run("error while unregister game", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().UnregisterGame(fx.ctx, &registrator.UnregisterGameRequest{
			GameId: 1,
		}).Return(nil, errors.New("some error"))

		got, err := fx.facade.UnregisterGame(fx.ctx, 1)
		assert.Equal(t, int32(registrator.UnregisterGameStatus_UNREGISTER_GAME_STATUS_INVALID), got)
		assert.Error(t, err)
	})

	t.Run("error game not found while unregister game", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().UnregisterGame(fx.ctx, &registrator.UnregisterGameRequest{
			GameId: 1,
		}).Return(nil, status.New(codes.NotFound, "").Err())

		got, err := fx.facade.UnregisterGame(fx.ctx, 1)
		assert.Equal(t, int32(registrator.UnregisterGameStatus_UNREGISTER_GAME_STATUS_INVALID), got)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrGameNotFound)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().UnregisterGame(fx.ctx, &registrator.UnregisterGameRequest{
			GameId: 1,
		}).Return(&registrator.UnregisterGameResponse{
			Status: registrator.UnregisterGameStatus_UNREGISTER_GAME_STATUS_OK,
		}, nil)

		got, err := fx.facade.UnregisterGame(fx.ctx, 1)
		assert.Equal(t, int32(registrator.UnregisterGameStatus_UNREGISTER_GAME_STATUS_OK), got)
		assert.NoError(t, err)
	})
}
