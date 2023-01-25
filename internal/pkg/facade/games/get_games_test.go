package games

import (
	"errors"
	"testing"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFacade_GetGameByID(t *testing.T) {
	zeroDateTime := timestamppb.Timestamp{}
	t.Run("error while get game by ID", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetGameByID(fx.ctx, &registrator.GetGameByIDRequest{
			GameId: 1,
		}).Return(nil, errors.New("some error"))

		got, err := fx.facade.GetGameByID(fx.ctx, 1)
		assert.Equal(t, model.Game{}, got)
		assert.Error(t, err)
	})

	t.Run("error whil get league by ID", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetGameByID(fx.ctx, &registrator.GetGameByIDRequest{
			GameId: 1,
		}).Return(&registrator.GetGameByIDResponse{
			Game: &registrator.Game{
				Id:       1,
				LeagueId: 1,
				PlaceId:  1,
			},
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(1)).Return(model.League{}, errors.New("some error"))

		got, err := fx.facade.GetGameByID(fx.ctx, 1)
		assert.Equal(t, model.Game{}, got)
		assert.Error(t, err)
	})

	t.Run("error while get place by ID", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetGameByID(fx.ctx, &registrator.GetGameByIDRequest{
			GameId: 1,
		}).Return(&registrator.GetGameByIDResponse{
			Game: &registrator.Game{
				Id:       1,
				LeagueId: 1,
				PlaceId:  1,
			},
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(1)).Return(model.League{
			ID: 1,
		}, nil)

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(1)).Return(model.Place{}, errors.New("some error"))

		got, err := fx.facade.GetGameByID(fx.ctx, 1)
		assert.Equal(t, model.Game{}, got)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetGameByID(fx.ctx, &registrator.GetGameByIDRequest{
			GameId: 1,
		}).Return(&registrator.GetGameByIDResponse{
			Game: &registrator.Game{
				Id:       1,
				LeagueId: 1,
				PlaceId:  1,
			},
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(1)).Return(model.League{
			ID: 1,
		}, nil)

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(1)).Return(model.Place{
			ID: 1,
		}, nil)

		got, err := fx.facade.GetGameByID(fx.ctx, 1)
		assert.Equal(t, model.Game{
			ID: 1,
			League: model.League{
				ID: 1,
			},
			Place: model.Place{
				ID: 1,
			},
			Date: model.DateTime(zeroDateTime.AsTime()),
		}, got)
		assert.NoError(t, err)
	})
}

func TestFacade_GetGames(t *testing.T) {
	zeroDateTime := timestamppb.Timestamp{}
	t.Run("error while get games", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetGames(fx.ctx, &registrator.GetGamesRequest{
			Active: true,
		}).Return(nil, errors.New("some error"))

		got, err := fx.facade.GetGames(fx.ctx, true)
		assert.Nil(t, got)
		assert.Error(t, err)
	})

	t.Run("error while get league", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetGames(fx.ctx, &registrator.GetGamesRequest{
			Active: true,
		}).Return(&registrator.GetGamesResponse{
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
				{
					Id:       4,
					LeagueId: 2,
					PlaceId:  1,
				},
				{
					Id:       5,
					LeagueId: 3,
					PlaceId:  3,
				},
				{
					Id:       6,
					LeagueId: 3,
					PlaceId:  4,
				},
			},
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(1)).Return(model.League{
			ID: 1,
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(2)).Return(model.League{
			ID: 2,
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(3)).Return(model.League{}, errors.New("some error"))

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(1)).Return(model.Place{
			ID: 1,
		}, nil)

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(2)).Return(model.Place{
			ID: 2,
		}, nil)

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(3)).Return(model.Place{
			ID: 3,
		}, nil)

		got, err := fx.facade.GetGames(fx.ctx, true)
		assert.Nil(t, got)
		assert.Error(t, err)
	})

	t.Run("error while get place", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetGames(fx.ctx, &registrator.GetGamesRequest{
			Active: true,
		}).Return(&registrator.GetGamesResponse{
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
				{
					Id:       4,
					LeagueId: 2,
					PlaceId:  1,
				},
				{
					Id:       5,
					LeagueId: 3,
					PlaceId:  3,
				},
				{
					Id:       6,
					LeagueId: 3,
					PlaceId:  4,
				},
			},
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(1)).Return(model.League{
			ID: 1,
		}, nil)

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(1)).Return(model.Place{
			ID: 1,
		}, nil)

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(2)).Return(model.Place{}, errors.New("some error"))

		got, err := fx.facade.GetGames(fx.ctx, true)
		assert.Nil(t, got)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetGames(fx.ctx, &registrator.GetGamesRequest{
			Active: true,
		}).Return(&registrator.GetGamesResponse{
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
				{
					Id:       4,
					LeagueId: 2,
					PlaceId:  1,
				},
				{
					Id:       5,
					LeagueId: 3,
					PlaceId:  3,
				},
				{
					Id:       6,
					LeagueId: 3,
					PlaceId:  4,
				},
			},
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(1)).Return(model.League{
			ID: 1,
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(2)).Return(model.League{
			ID: 2,
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(3)).Return(model.League{
			ID: 3,
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

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(4)).Return(model.Place{
			ID: 4,
		}, nil)

		got, err := fx.facade.GetGames(fx.ctx, true)
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
			{
				ID: 4,
				League: model.League{
					ID: 2,
				},
				Place: model.Place{
					ID: 1,
				},
				Date: model.DateTime(zeroDateTime.AsTime()),
			},
			{
				ID: 5,
				League: model.League{
					ID: 3,
				},
				Place: model.Place{
					ID: 3,
				},
				Date: model.DateTime(zeroDateTime.AsTime()),
			},
			{
				ID: 6,
				League: model.League{
					ID: 3,
				},
				Place: model.Place{
					ID: 4,
				},
				Date: model.DateTime(zeroDateTime.AsTime()),
			},
		}, got)
		assert.NoError(t, err)
	})
}

func TestFacade_GetRegisteredGames(t *testing.T) {
	zeroDateTime := timestamppb.Timestamp{}
	t.Run("error while get registered games", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetRegisteredGames(fx.ctx, &registrator.GetRegisteredGamesRequest{
			Active: true,
		}).Return(nil, errors.New("some error"))

		got, err := fx.facade.GetRegisteredGames(fx.ctx, true)
		assert.Nil(t, got)
		assert.Error(t, err)
	})

	t.Run("error while get league", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetRegisteredGames(fx.ctx, &registrator.GetRegisteredGamesRequest{
			Active: true,
		}).Return(&registrator.GetRegisteredGamesResponse{
			Games: []*registrator.Game{
				{
					Id:         1,
					LeagueId:   1,
					PlaceId:    1,
					Registered: true,
				},
				{
					Id:         2,
					LeagueId:   1,
					PlaceId:    2,
					Registered: true,
				},
				{
					Id:         3,
					LeagueId:   2,
					PlaceId:    3,
					Registered: true,
				},
				{
					Id:         4,
					LeagueId:   2,
					PlaceId:    1,
					Registered: true,
				},
				{
					Id:         5,
					LeagueId:   3,
					PlaceId:    3,
					Registered: true,
				},
				{
					Id:         6,
					LeagueId:   3,
					PlaceId:    4,
					Registered: true,
				},
			},
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(1)).Return(model.League{
			ID: 1,
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(2)).Return(model.League{
			ID: 2,
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(3)).Return(model.League{}, errors.New("some error"))

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(1)).Return(model.Place{
			ID: 1,
		}, nil)

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(2)).Return(model.Place{
			ID: 2,
		}, nil)

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(3)).Return(model.Place{
			ID: 3,
		}, nil)

		got, err := fx.facade.GetRegisteredGames(fx.ctx, true)
		assert.Nil(t, got)
		assert.Error(t, err)
	})

	t.Run("error while get place", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetRegisteredGames(fx.ctx, &registrator.GetRegisteredGamesRequest{
			Active: true,
		}).Return(&registrator.GetRegisteredGamesResponse{
			Games: []*registrator.Game{
				{
					Id:         1,
					LeagueId:   1,
					PlaceId:    1,
					Registered: true,
				},
				{
					Id:         2,
					LeagueId:   1,
					PlaceId:    2,
					Registered: true,
				},
				{
					Id:         3,
					LeagueId:   2,
					PlaceId:    3,
					Registered: true,
				},
				{
					Id:         4,
					LeagueId:   2,
					PlaceId:    1,
					Registered: true,
				},
				{
					Id:         5,
					LeagueId:   3,
					PlaceId:    3,
					Registered: true,
				},
				{
					Id:         6,
					LeagueId:   3,
					PlaceId:    4,
					Registered: true,
				},
			},
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(1)).Return(model.League{
			ID: 1,
		}, nil)

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(1)).Return(model.Place{
			ID: 1,
		}, nil)

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(2)).Return(model.Place{}, errors.New("some error"))

		got, err := fx.facade.GetRegisteredGames(fx.ctx, true)
		assert.Nil(t, got)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetRegisteredGames(fx.ctx, &registrator.GetRegisteredGamesRequest{
			Active: true,
		}).Return(&registrator.GetRegisteredGamesResponse{
			Games: []*registrator.Game{
				{
					Id:         1,
					LeagueId:   1,
					PlaceId:    1,
					Registered: true,
				},
				{
					Id:         2,
					LeagueId:   1,
					PlaceId:    2,
					Registered: true,
				},
				{
					Id:         3,
					LeagueId:   2,
					PlaceId:    3,
					Registered: true,
				},
				{
					Id:         4,
					LeagueId:   2,
					PlaceId:    1,
					Registered: true,
				},
				{
					Id:         5,
					LeagueId:   3,
					PlaceId:    3,
					Registered: true,
				},
				{
					Id:         6,
					LeagueId:   3,
					PlaceId:    4,
					Registered: true,
				},
			},
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(1)).Return(model.League{
			ID: 1,
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(2)).Return(model.League{
			ID: 2,
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(3)).Return(model.League{
			ID: 3,
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

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(4)).Return(model.Place{
			ID: 4,
		}, nil)

		got, err := fx.facade.GetRegisteredGames(fx.ctx, true)
		assert.Equal(t, []model.Game{
			{
				ID: 1,
				League: model.League{
					ID: 1,
				},
				Place: model.Place{
					ID: 1,
				},
				Date:       model.DateTime(zeroDateTime.AsTime()),
				Registered: true,
			},
			{
				ID: 2,
				League: model.League{
					ID: 1,
				},
				Place: model.Place{
					ID: 2,
				},
				Date:       model.DateTime(zeroDateTime.AsTime()),
				Registered: true,
			},
			{
				ID: 3,
				League: model.League{
					ID: 2,
				},
				Place: model.Place{
					ID: 3,
				},
				Date:       model.DateTime(zeroDateTime.AsTime()),
				Registered: true,
			},
			{
				ID: 4,
				League: model.League{
					ID: 2,
				},
				Place: model.Place{
					ID: 1,
				},
				Date:       model.DateTime(zeroDateTime.AsTime()),
				Registered: true,
			},
			{
				ID: 5,
				League: model.League{
					ID: 3,
				},
				Place: model.Place{
					ID: 3,
				},
				Date:       model.DateTime(zeroDateTime.AsTime()),
				Registered: true,
			},
			{
				ID: 6,
				League: model.League{
					ID: 3,
				},
				Place: model.Place{
					ID: 4,
				},
				Date:       model.DateTime(zeroDateTime.AsTime()),
				Registered: true,
			},
		}, got)
		assert.NoError(t, err)
	})
}

func TestFacade_GetUserGames(t *testing.T) {
	zeroDateTime := timestamppb.Timestamp{}
	t.Run("error while get user games", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetUserGames(fx.ctx, &registrator.GetUserGamesRequest{
			Active: true,
			UserId: 1,
		}).Return(nil, errors.New("some error"))

		got, err := fx.facade.GetUserGames(fx.ctx, true, 1)
		assert.Nil(t, got)
		assert.Error(t, err)
	})

	t.Run("error while get league", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetUserGames(fx.ctx, &registrator.GetUserGamesRequest{
			Active: true,
			UserId: 1,
		}).Return(&registrator.GetUserGamesResponse{
			Games: []*registrator.Game{
				{
					Id:         1,
					LeagueId:   1,
					PlaceId:    1,
					Registered: true,
				},
				{
					Id:         2,
					LeagueId:   1,
					PlaceId:    2,
					Registered: true,
				},
				{
					Id:         3,
					LeagueId:   2,
					PlaceId:    3,
					Registered: true,
				},
				{
					Id:         4,
					LeagueId:   2,
					PlaceId:    1,
					Registered: true,
				},
				{
					Id:         5,
					LeagueId:   3,
					PlaceId:    3,
					Registered: true,
				},
				{
					Id:         6,
					LeagueId:   3,
					PlaceId:    4,
					Registered: true,
				},
			},
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(1)).Return(model.League{
			ID: 1,
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(2)).Return(model.League{
			ID: 2,
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(3)).Return(model.League{}, errors.New("some error"))

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(1)).Return(model.Place{
			ID: 1,
		}, nil)

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(2)).Return(model.Place{
			ID: 2,
		}, nil)

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(3)).Return(model.Place{
			ID: 3,
		}, nil)

		got, err := fx.facade.GetUserGames(fx.ctx, true, 1)
		assert.Nil(t, got)
		assert.Error(t, err)
	})

	t.Run("error while get place", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetUserGames(fx.ctx, &registrator.GetUserGamesRequest{
			Active: true,
			UserId: 1,
		}).Return(&registrator.GetUserGamesResponse{
			Games: []*registrator.Game{
				{
					Id:         1,
					LeagueId:   1,
					PlaceId:    1,
					Registered: true,
				},
				{
					Id:         2,
					LeagueId:   1,
					PlaceId:    2,
					Registered: true,
				},
				{
					Id:         3,
					LeagueId:   2,
					PlaceId:    3,
					Registered: true,
				},
				{
					Id:         4,
					LeagueId:   2,
					PlaceId:    1,
					Registered: true,
				},
				{
					Id:         5,
					LeagueId:   3,
					PlaceId:    3,
					Registered: true,
				},
				{
					Id:         6,
					LeagueId:   3,
					PlaceId:    4,
					Registered: true,
				},
			},
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(1)).Return(model.League{
			ID: 1,
		}, nil)

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(1)).Return(model.Place{
			ID: 1,
		}, nil)

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(2)).Return(model.Place{}, errors.New("some error"))

		got, err := fx.facade.GetUserGames(fx.ctx, true, 1)
		assert.Nil(t, got)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetUserGames(fx.ctx, &registrator.GetUserGamesRequest{
			Active: true,
			UserId: 1,
		}).Return(&registrator.GetUserGamesResponse{
			Games: []*registrator.Game{
				{
					Id:         1,
					LeagueId:   1,
					PlaceId:    1,
					Registered: true,
				},
				{
					Id:         2,
					LeagueId:   1,
					PlaceId:    2,
					Registered: true,
				},
				{
					Id:         3,
					LeagueId:   2,
					PlaceId:    3,
					Registered: true,
				},
				{
					Id:         4,
					LeagueId:   2,
					PlaceId:    1,
					Registered: true,
				},
				{
					Id:         5,
					LeagueId:   3,
					PlaceId:    3,
					Registered: true,
				},
				{
					Id:         6,
					LeagueId:   3,
					PlaceId:    4,
					Registered: true,
				},
			},
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(1)).Return(model.League{
			ID: 1,
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(2)).Return(model.League{
			ID: 2,
		}, nil)

		fx.leaguesFacade.EXPECT().GetLeagueByID(fx.ctx, int32(3)).Return(model.League{
			ID: 3,
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

		fx.placesFacade.EXPECT().GetPlaceByID(fx.ctx, int32(4)).Return(model.Place{
			ID: 4,
		}, nil)

		got, err := fx.facade.GetUserGames(fx.ctx, true, 1)
		assert.Equal(t, []model.Game{
			{
				ID: 1,
				League: model.League{
					ID: 1,
				},
				Place: model.Place{
					ID: 1,
				},
				Date:       model.DateTime(zeroDateTime.AsTime()),
				Registered: true,
			},
			{
				ID: 2,
				League: model.League{
					ID: 1,
				},
				Place: model.Place{
					ID: 2,
				},
				Date:       model.DateTime(zeroDateTime.AsTime()),
				Registered: true,
			},
			{
				ID: 3,
				League: model.League{
					ID: 2,
				},
				Place: model.Place{
					ID: 3,
				},
				Date:       model.DateTime(zeroDateTime.AsTime()),
				Registered: true,
			},
			{
				ID: 4,
				League: model.League{
					ID: 2,
				},
				Place: model.Place{
					ID: 1,
				},
				Date:       model.DateTime(zeroDateTime.AsTime()),
				Registered: true,
			},
			{
				ID: 5,
				League: model.League{
					ID: 3,
				},
				Place: model.Place{
					ID: 3,
				},
				Date:       model.DateTime(zeroDateTime.AsTime()),
				Registered: true,
			},
			{
				ID: 6,
				League: model.League{
					ID: 3,
				},
				Place: model.Place{
					ID: 4,
				},
				Date:       model.DateTime(zeroDateTime.AsTime()),
				Registered: true,
			},
		}, got)
		assert.NoError(t, err)
	})
}
