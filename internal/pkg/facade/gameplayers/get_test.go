package gameplayers

import (
	"errors"
	"testing"

	"github.com/mono83/maybe"
	gameplayerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/game_player"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestFacade_GetPlayersByGameID(t *testing.T) {
	t.Run("error while get players by game ID", func(t *testing.T) {
		fx := tearUp(t)

		fx.gamePlayerServiceClient.EXPECT().GetGamePlayersByGameID(fx.ctx, &gameplayerpb.GetGamePlayersByGameIDRequest{
			GameId: 1,
		}).Return(nil, errors.New("some error"))

		got, err := fx.facade.GetGamePlayersByGameID(fx.ctx, 1)
		assert.Nil(t, got)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.gamePlayerServiceClient.EXPECT().GetGamePlayersByGameID(fx.ctx, &gameplayerpb.GetGamePlayersByGameIDRequest{
			GameId: 1,
		}).Return(&gameplayerpb.GetGamePlayersByGameIDResponse{
			GamePlayers: []*gameplayerpb.GamePlayer{
				{
					Id:           1,
					GameId:       1,
					UserId:       wrapperspb.Int32(1),
					RegisteredBy: 1,
					Degree:       gameplayerpb.Degree_DEGREE_LIKELY,
				},
				{
					Id:           2,
					GameId:       2,
					UserId:       wrapperspb.Int32(2),
					RegisteredBy: 2,
					Degree:       gameplayerpb.Degree_DEGREE_UNLIKELY,
				},
				{
					Id:           3,
					GameId:       2,
					RegisteredBy: 1,
					Degree:       gameplayerpb.Degree_DEGREE_UNLIKELY,
				},
			},
		}, nil)

		got, err := fx.facade.GetGamePlayersByGameID(fx.ctx, 1)
		assert.Equal(t, []model.GamePlayer{
			{
				ID:           1,
				GameID:       1,
				UserID:       maybe.Just(int32(1)),
				RegisteredBy: 1,
				Degree:       model.DegreeLikely,
			},
			{
				ID:           2,
				GameID:       2,
				UserID:       maybe.Just(int32(2)),
				RegisteredBy: 2,
				Degree:       model.DegreeUnlikely,
			},
			{
				ID:           3,
				GameID:       2,
				UserID:       maybe.Nothing[int32](),
				RegisteredBy: 1,
				Degree:       model.DegreeUnlikely,
			},
		}, got)
		assert.NoError(t, err)
	})
}
