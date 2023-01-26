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

func TestFacade_UnregisterPlayer(t *testing.T) {
	t.Run("error while unregister player", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().UnregisterPlayer(fx.ctx, &registrator.UnregisterPlayerRequest{
			GameId:     1,
			PlayerType: registrator.PlayerType_PLAYER_TYPE_LEGIONER,
		}).Return(nil, errors.New("some error"))

		got, err := fx.facade.UnregisterPlayer(fx.ctx, 1, int32(registrator.PlayerType_PLAYER_TYPE_LEGIONER))
		assert.Equal(t, int32(registrator.UnregisterPlayerStatus_UNREGISTER_PLAYER_STATUS_INVALID), got)
		assert.Error(t, err)
	})

	t.Run("error game not found while unregister player", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().UnregisterPlayer(fx.ctx, &registrator.UnregisterPlayerRequest{
			GameId:     1,
			PlayerType: registrator.PlayerType_PLAYER_TYPE_LEGIONER,
		}).Return(nil, status.New(codes.NotFound, "").Err())

		got, err := fx.facade.UnregisterPlayer(fx.ctx, 1, int32(registrator.PlayerType_PLAYER_TYPE_LEGIONER))
		assert.Equal(t, int32(registrator.UnregisterPlayerStatus_UNREGISTER_PLAYER_STATUS_INVALID), got)
		assert.Error(t, err)
		assert.Error(t, err, model.ErrGameNotFound)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().UnregisterPlayer(fx.ctx, &registrator.UnregisterPlayerRequest{
			GameId:     1,
			PlayerType: registrator.PlayerType_PLAYER_TYPE_LEGIONER,
		}).Return(&registrator.UnregisterPlayerResponse{
			Status: registrator.UnregisterPlayerStatus_UNREGISTER_PLAYER_STATUS_OK,
		}, nil)

		got, err := fx.facade.UnregisterPlayer(fx.ctx, 1, int32(registrator.PlayerType_PLAYER_TYPE_LEGIONER))
		assert.Equal(t, int32(registrator.UnregisterPlayerStatus_UNREGISTER_PLAYER_STATUS_OK), got)
		assert.NoError(t, err)
	})
}
