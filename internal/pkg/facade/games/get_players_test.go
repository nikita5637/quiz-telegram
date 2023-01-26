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

func TestFacade_GetPlayersByGameID(t *testing.T) {
	t.Run("error while get players by game ID", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetPlayersByGameID(fx.ctx, &registrator.GetPlayersByGameIDRequest{
			GameId: 1,
		}).Return(nil, errors.New("some error"))

		got, err := fx.facade.GetPlayersByGameID(fx.ctx, 1)
		assert.Nil(t, got)
		assert.Error(t, err)
	})

	t.Run("error game not found while get players by game ID", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetPlayersByGameID(fx.ctx, &registrator.GetPlayersByGameIDRequest{
			GameId: 1,
		}).Return(nil, status.New(codes.NotFound, "").Err())

		got, err := fx.facade.GetPlayersByGameID(fx.ctx, 1)
		assert.Nil(t, got)
		assert.Error(t, err)
		assert.ErrorIs(t, err, model.ErrGameNotFound)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetPlayersByGameID(fx.ctx, &registrator.GetPlayersByGameIDRequest{
			GameId: 1,
		}).Return(&registrator.GetPlayersByGameIDResponse{
			Players: []*registrator.Player{
				{
					UserId:       1,
					RegisteredBy: 1,
					Degree:       registrator.Degree_DEGREE_LIKELY,
				},
				{
					UserId:       2,
					RegisteredBy: 2,
					Degree:       registrator.Degree_DEGREE_UNLIKELY,
				},
				{
					UserId:       0,
					RegisteredBy: 1,
					Degree:       registrator.Degree_DEGREE_UNLIKELY,
				},
			},
		}, nil)

		got, err := fx.facade.GetPlayersByGameID(fx.ctx, 1)
		assert.Equal(t, []model.Player{
			{
				UserID:       1,
				RegisteredBy: 1,
				Degree:       int32(registrator.Degree_DEGREE_LIKELY),
			},
			{
				UserID:       2,
				RegisteredBy: 2,
				Degree:       int32(registrator.Degree_DEGREE_UNLIKELY),
			},
			{
				UserID:       0,
				RegisteredBy: 1,
				Degree:       int32(registrator.Degree_DEGREE_UNLIKELY),
			},
		}, got)
		assert.NoError(t, err)
	})
}
