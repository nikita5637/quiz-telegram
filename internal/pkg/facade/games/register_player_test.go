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

func TestFacade_RegisterPlayer(t *testing.T) {
	t.Run("error while register player", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().RegisterPlayer(fx.ctx, &registrator.RegisterPlayerRequest{
			GameId:     1,
			PlayerType: registrator.PlayerType_PLAYER_TYPE_MAIN,
			Degree:     registrator.Degree_DEGREE_LIKELY,
		}).Return(nil, errors.New("some error"))

		got, err := fx.facade.RegisterPlayer(fx.ctx,
			1,
			int32(registrator.PlayerType_PLAYER_TYPE_MAIN),
			int32(registrator.Degree_DEGREE_LIKELY),
		)
		assert.Equal(t, int32(registrator.RegisterPlayerStatus_REGISTER_PLAYER_STATUS_INVALID), got)
		assert.Error(t, err)
	})

	t.Run("error game not found while register player", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().RegisterPlayer(fx.ctx, &registrator.RegisterPlayerRequest{
			GameId:     1,
			PlayerType: registrator.PlayerType_PLAYER_TYPE_MAIN,
			Degree:     registrator.Degree_DEGREE_LIKELY,
		}).Return(nil, status.New(codes.NotFound, "").Err())

		got, err := fx.facade.RegisterPlayer(fx.ctx,
			1,
			int32(registrator.PlayerType_PLAYER_TYPE_MAIN),
			int32(registrator.Degree_DEGREE_LIKELY),
		)
		assert.Equal(t, int32(registrator.RegisterPlayerStatus_REGISTER_PLAYER_STATUS_INVALID), got)
		assert.Error(t, err)
		assert.ErrorIs(t, err, model.ErrGameNotFound)
	})

	t.Run("error no free slot while register player", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().RegisterPlayer(fx.ctx, &registrator.RegisterPlayerRequest{
			GameId:     1,
			PlayerType: registrator.PlayerType_PLAYER_TYPE_MAIN,
			Degree:     registrator.Degree_DEGREE_LIKELY,
		}).Return(nil, status.New(codes.AlreadyExists, "").Err())

		got, err := fx.facade.RegisterPlayer(fx.ctx,
			1,
			int32(registrator.PlayerType_PLAYER_TYPE_MAIN),
			int32(registrator.Degree_DEGREE_LIKELY),
		)
		assert.Equal(t, int32(registrator.RegisterPlayerStatus_REGISTER_PLAYER_STATUS_INVALID), got)
		assert.Error(t, err)
		assert.ErrorIs(t, err, model.ErrNoFreeSlot)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().RegisterPlayer(fx.ctx, &registrator.RegisterPlayerRequest{
			GameId:     1,
			PlayerType: registrator.PlayerType_PLAYER_TYPE_MAIN,
			Degree:     registrator.Degree_DEGREE_LIKELY,
		}).Return(&registrator.RegisterPlayerResponse{
			Status: registrator.RegisterPlayerStatus_REGISTER_PLAYER_STATUS_OK,
		}, nil)

		got, err := fx.facade.RegisterPlayer(fx.ctx,
			1,
			int32(registrator.PlayerType_PLAYER_TYPE_MAIN),
			int32(registrator.Degree_DEGREE_LIKELY),
		)
		assert.Equal(t, int32(registrator.RegisterPlayerStatus_REGISTER_PLAYER_STATUS_OK), got)
		assert.NoError(t, err)
	})
}
