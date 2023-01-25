package gamephotos

import (
	"errors"
	"testing"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFacade_GetGamesWithPhotos(t *testing.T) {
	zeroDateTime := timestamppb.Timestamp{}
	t.Run("error while get games with photos", func(t *testing.T) {
		fx := tearUp(t)

		fx.photographerServiceClient.EXPECT().GetGamesWithPhotos(fx.ctx, &registrator.GetGamesWithPhotosRequest{
			Limit:  4,
			Offset: 0,
		}).Return(nil, errors.New("some error"))

		got1, got2, err := fx.facade.GetGamesWithPhotos(fx.ctx, 4, 0)
		assert.Nil(t, got1)
		assert.Equal(t, uint32(0), got2)
		assert.Error(t, err)
	})

	t.Run("error while get league", func(t *testing.T) {
		fx := tearUp(t)

		fx.photographerServiceClient.EXPECT().GetGamesWithPhotos(fx.ctx, &registrator.GetGamesWithPhotosRequest{
			Limit:  4,
			Offset: 0,
		}).Return(&registrator.GetGamesWithPhotosResponse{
			Games: []*registrator.Game{
				{
					Id:       1,
					LeagueId: 1,
					PlaceId:  1,
				},
				{
					Id:       2,
					LeagueId: 1,
					PlaceId:  2,
				},
				{
					Id:       3,
					LeagueId: 2,
					PlaceId:  3,
				},
			},
			Total: 3,
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(1)).Return(model.League{
			ID: 1,
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(2)).Return(model.League{}, errors.New("some error"))

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(1)).Return(model.Place{
			ID: 1,
		}, nil)

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(2)).Return(model.Place{
			ID: 2,
		}, nil)

		got1, got2, err := fx.facade.GetGamesWithPhotos(fx.ctx, 4, 0)
		assert.Nil(t, got1)
		assert.Equal(t, uint32(0), got2)
		assert.Error(t, err)
	})

	t.Run("error while get place", func(t *testing.T) {
		fx := tearUp(t)

		fx.photographerServiceClient.EXPECT().GetGamesWithPhotos(fx.ctx, &registrator.GetGamesWithPhotosRequest{
			Limit:  4,
			Offset: 0,
		}).Return(&registrator.GetGamesWithPhotosResponse{
			Games: []*registrator.Game{
				{
					Id:       1,
					LeagueId: 1,
					PlaceId:  1,
				},
				{
					Id:       2,
					LeagueId: 1,
					PlaceId:  2,
				},
				{
					Id:       3,
					LeagueId: 2,
					PlaceId:  3,
				},
			},
			Total: 3,
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(1)).Return(model.League{
			ID: 1,
		}, nil)

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(1)).Return(model.Place{
			ID: 1,
		}, nil)

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(2)).Return(model.Place{}, errors.New("some error"))

		got1, got2, err := fx.facade.GetGamesWithPhotos(fx.ctx, 4, 0)
		assert.Nil(t, got1)
		assert.Equal(t, uint32(0), got2)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.photographerServiceClient.EXPECT().GetGamesWithPhotos(fx.ctx, &registrator.GetGamesWithPhotosRequest{
			Limit:  4,
			Offset: 0,
		}).Return(&registrator.GetGamesWithPhotosResponse{
			Games: []*registrator.Game{
				{
					Id:       1,
					LeagueId: 1,
					PlaceId:  1,
				},
				{
					Id:       2,
					LeagueId: 1,
					PlaceId:  2,
				},
				{
					Id:       3,
					LeagueId: 2,
					PlaceId:  3,
				},
			},
			Total: 3,
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(1)).Return(model.League{
			ID: 1,
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(2)).Return(model.League{
			ID: 2,
		}, nil)

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(1)).Return(model.Place{
			ID: 1,
		}, nil)

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(2)).Return(model.Place{
			ID: 2,
		}, nil)

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(3)).Return(model.Place{
			ID: 3,
		}, nil)

		got1, got2, err := fx.facade.GetGamesWithPhotos(fx.ctx, 4, 0)
		assert.Equal(t, []model.Game{
			{
				ID: 1,
				League: model.League{
					ID: 1,
				},
				Place: model.Place{
					ID: 1,
				},
				Date: model.DateTime(zeroDateTime.AsTime()),
			},
			{
				ID: 2,
				League: model.League{
					ID: 1,
				},
				Place: model.Place{
					ID: 2,
				},
				Date: model.DateTime(zeroDateTime.AsTime()),
			},
			{
				ID: 3,
				League: model.League{
					ID: 2,
				},
				Place: model.Place{
					ID: 3,
				},
				Date: model.DateTime(zeroDateTime.AsTime()),
			},
		}, got1)
		assert.Equal(t, uint32(3), got2)
		assert.NoError(t, err)
	})
}
