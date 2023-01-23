package games

import (
	"errors"
	"reflect"
	"testing"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	time_utils "github.com/nikita5637/quiz-telegram/utils/time"
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

		fx.registratorServiceClient.EXPECT().GetLeagueByID(fx.ctx, &registrator.GetLeagueByIDRequest{
			Id: 1,
		}).Return(nil, errors.New("some error"))

		got, err := fx.facade.GetGameByID(fx.ctx, 1)
		assert.Equal(t, model.Game{}, got)
		assert.Error(t, err)
	})

	t.Run("error while get place", func(t *testing.T) {
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

		fx.registratorServiceClient.EXPECT().GetLeagueByID(fx.ctx, &registrator.GetLeagueByIDRequest{
			Id: 1,
		}).Return(&registrator.GetLeagueByIDResponse{
			League: &registrator.League{
				Id: 1,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetPlaceByID(fx.ctx, &registrator.GetPlaceByIDRequest{
			Id: 1,
		}).Return(nil, errors.New("some error"))

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

		fx.registratorServiceClient.EXPECT().GetLeagueByID(fx.ctx, &registrator.GetLeagueByIDRequest{
			Id: 1,
		}).Return(&registrator.GetLeagueByIDResponse{
			League: &registrator.League{
				Id: 1,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetPlaceByID(fx.ctx, &registrator.GetPlaceByIDRequest{
			Id: 1,
		}).Return(&registrator.GetPlaceByIDResponse{
			Place: &registrator.Place{
				Id: 1,
			},
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

		fx.registratorServiceClient.EXPECT().GetLeagueByID(fx.ctx, &registrator.GetLeagueByIDRequest{
			Id: 1,
		}).Once().Return(&registrator.GetLeagueByIDResponse{
			League: &registrator.League{
				Id: 1,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetLeagueByID(fx.ctx, &registrator.GetLeagueByIDRequest{
			Id: 2,
		}).Once().Return(&registrator.GetLeagueByIDResponse{
			League: &registrator.League{
				Id: 2,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetLeagueByID(fx.ctx, &registrator.GetLeagueByIDRequest{
			Id: 3,
		}).Once().Return(nil, errors.New("some error"))

		fx.registratorServiceClient.EXPECT().GetPlaceByID(fx.ctx, &registrator.GetPlaceByIDRequest{
			Id: 1,
		}).Once().Return(&registrator.GetPlaceByIDResponse{
			Place: &registrator.Place{
				Id: 1,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetPlaceByID(fx.ctx, &registrator.GetPlaceByIDRequest{
			Id: 2,
		}).Once().Return(&registrator.GetPlaceByIDResponse{
			Place: &registrator.Place{
				Id: 2,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetPlaceByID(fx.ctx, &registrator.GetPlaceByIDRequest{
			Id: 3,
		}).Once().Return(&registrator.GetPlaceByIDResponse{
			Place: &registrator.Place{
				Id: 3,
			},
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

		fx.registratorServiceClient.EXPECT().GetLeagueByID(fx.ctx, &registrator.GetLeagueByIDRequest{
			Id: 1,
		}).Once().Return(&registrator.GetLeagueByIDResponse{
			League: &registrator.League{
				Id: 1,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetPlaceByID(fx.ctx, &registrator.GetPlaceByIDRequest{
			Id: 1,
		}).Once().Return(&registrator.GetPlaceByIDResponse{
			Place: &registrator.Place{
				Id: 1,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetPlaceByID(fx.ctx, &registrator.GetPlaceByIDRequest{
			Id: 2,
		}).Once().Return(nil, errors.New("some error"))

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

		fx.registratorServiceClient.EXPECT().GetLeagueByID(fx.ctx, &registrator.GetLeagueByIDRequest{
			Id: 1,
		}).Once().Return(&registrator.GetLeagueByIDResponse{
			League: &registrator.League{
				Id: 1,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetLeagueByID(fx.ctx, &registrator.GetLeagueByIDRequest{
			Id: 2,
		}).Once().Return(&registrator.GetLeagueByIDResponse{
			League: &registrator.League{
				Id: 2,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetLeagueByID(fx.ctx, &registrator.GetLeagueByIDRequest{
			Id: 3,
		}).Once().Return(&registrator.GetLeagueByIDResponse{
			League: &registrator.League{
				Id: 3,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetPlaceByID(fx.ctx, &registrator.GetPlaceByIDRequest{
			Id: 1,
		}).Once().Return(&registrator.GetPlaceByIDResponse{
			Place: &registrator.Place{
				Id: 1,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetPlaceByID(fx.ctx, &registrator.GetPlaceByIDRequest{
			Id: 2,
		}).Once().Return(&registrator.GetPlaceByIDResponse{
			Place: &registrator.Place{
				Id: 2,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetPlaceByID(fx.ctx, &registrator.GetPlaceByIDRequest{
			Id: 3,
		}).Once().Return(&registrator.GetPlaceByIDResponse{
			Place: &registrator.Place{
				Id: 3,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetPlaceByID(fx.ctx, &registrator.GetPlaceByIDRequest{
			Id: 4,
		}).Once().Return(&registrator.GetPlaceByIDResponse{
			Place: &registrator.Place{
				Id: 4,
			},
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

		fx.registratorServiceClient.EXPECT().GetLeagueByID(fx.ctx, &registrator.GetLeagueByIDRequest{
			Id: 1,
		}).Once().Return(&registrator.GetLeagueByIDResponse{
			League: &registrator.League{
				Id: 1,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetLeagueByID(fx.ctx, &registrator.GetLeagueByIDRequest{
			Id: 2,
		}).Once().Return(&registrator.GetLeagueByIDResponse{
			League: &registrator.League{
				Id: 2,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetLeagueByID(fx.ctx, &registrator.GetLeagueByIDRequest{
			Id: 3,
		}).Once().Return(nil, errors.New("some error"))

		fx.registratorServiceClient.EXPECT().GetPlaceByID(fx.ctx, &registrator.GetPlaceByIDRequest{
			Id: 1,
		}).Once().Return(&registrator.GetPlaceByIDResponse{
			Place: &registrator.Place{
				Id: 1,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetPlaceByID(fx.ctx, &registrator.GetPlaceByIDRequest{
			Id: 2,
		}).Once().Return(&registrator.GetPlaceByIDResponse{
			Place: &registrator.Place{
				Id: 2,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetPlaceByID(fx.ctx, &registrator.GetPlaceByIDRequest{
			Id: 3,
		}).Once().Return(&registrator.GetPlaceByIDResponse{
			Place: &registrator.Place{
				Id: 3,
			},
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

		fx.registratorServiceClient.EXPECT().GetLeagueByID(fx.ctx, &registrator.GetLeagueByIDRequest{
			Id: 1,
		}).Once().Return(&registrator.GetLeagueByIDResponse{
			League: &registrator.League{
				Id: 1,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetPlaceByID(fx.ctx, &registrator.GetPlaceByIDRequest{
			Id: 1,
		}).Once().Return(&registrator.GetPlaceByIDResponse{
			Place: &registrator.Place{
				Id: 1,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetPlaceByID(fx.ctx, &registrator.GetPlaceByIDRequest{
			Id: 2,
		}).Once().Return(nil, errors.New("some error"))

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

		fx.registratorServiceClient.EXPECT().GetLeagueByID(fx.ctx, &registrator.GetLeagueByIDRequest{
			Id: 1,
		}).Once().Return(&registrator.GetLeagueByIDResponse{
			League: &registrator.League{
				Id: 1,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetLeagueByID(fx.ctx, &registrator.GetLeagueByIDRequest{
			Id: 2,
		}).Once().Return(&registrator.GetLeagueByIDResponse{
			League: &registrator.League{
				Id: 2,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetLeagueByID(fx.ctx, &registrator.GetLeagueByIDRequest{
			Id: 3,
		}).Once().Return(&registrator.GetLeagueByIDResponse{
			League: &registrator.League{
				Id: 3,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetPlaceByID(fx.ctx, &registrator.GetPlaceByIDRequest{
			Id: 1,
		}).Once().Return(&registrator.GetPlaceByIDResponse{
			Place: &registrator.Place{
				Id: 1,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetPlaceByID(fx.ctx, &registrator.GetPlaceByIDRequest{
			Id: 2,
		}).Once().Return(&registrator.GetPlaceByIDResponse{
			Place: &registrator.Place{
				Id: 2,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetPlaceByID(fx.ctx, &registrator.GetPlaceByIDRequest{
			Id: 3,
		}).Once().Return(&registrator.GetPlaceByIDResponse{
			Place: &registrator.Place{
				Id: 3,
			},
		}, nil)

		fx.registratorServiceClient.EXPECT().GetPlaceByID(fx.ctx, &registrator.GetPlaceByIDRequest{
			Id: 4,
		}).Once().Return(&registrator.GetPlaceByIDResponse{
			Place: &registrator.Place{
				Id: 4,
			},
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

func TestFacade_getModelLeague(t *testing.T) {
	t.Run("error while get league", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetLeagueByID(fx.ctx, &registrator.GetLeagueByIDRequest{
			Id: 1,
		}).Return(nil, errors.New("some error"))

		got, err := fx.facade.getModelLeague(fx.ctx, 1)
		assert.Equal(t, model.League{}, got)
		assert.Error(t, err)
	})

	t.Run("ok without cache", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetLeagueByID(fx.ctx, &registrator.GetLeagueByIDRequest{
			Id: 1,
		}).Return(&registrator.GetLeagueByIDResponse{
			League: &registrator.League{
				Id:   1,
				Name: "name",
			},
		}, nil)

		got, err := fx.facade.getModelLeague(fx.ctx, 1)
		assert.Equal(t, model.League{
			ID:   1,
			Name: "name",
		}, got)
		assert.NoError(t, err)
	})

	t.Run("ok with cache", func(t *testing.T) {
		fx := tearUp(t)
		fx.facade.leagueCache[1] = model.League{
			ID:   1,
			Name: "name",
		}

		got, err := fx.facade.getModelLeague(fx.ctx, 1)
		assert.Equal(t, model.League{
			ID:   1,
			Name: "name",
		}, got)
		assert.NoError(t, err)
	})
}

func TestFacade_getModelPlace(t *testing.T) {
	t.Run("error while get place", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetPlaceByID(fx.ctx, &registrator.GetPlaceByIDRequest{
			Id: 1,
		}).Return(nil, errors.New("some error"))

		got, err := fx.facade.getModelPlace(fx.ctx, 1)
		assert.Equal(t, model.Place{}, got)
		assert.Error(t, err)
	})

	t.Run("ok without cache", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().GetPlaceByID(fx.ctx, &registrator.GetPlaceByIDRequest{
			Id: 1,
		}).Return(&registrator.GetPlaceByIDResponse{
			Place: &registrator.Place{
				Id:   1,
				Name: "name",
			},
		}, nil)

		got, err := fx.facade.getModelPlace(fx.ctx, 1)
		assert.Equal(t, model.Place{
			ID:   1,
			Name: "name",
		}, got)
		assert.NoError(t, err)
	})

	t.Run("ok with cache", func(t *testing.T) {
		fx := tearUp(t)
		fx.facade.placeCache[1] = model.Place{
			ID:   1,
			Name: "name",
		}

		got, err := fx.facade.getModelPlace(fx.ctx, 1)
		assert.Equal(t, model.Place{
			ID:   1,
			Name: "name",
		}, got)
		assert.NoError(t, err)
	})
}

func Test_convertPBLeagueToModelLeague(t *testing.T) {
	type args struct {
		pbLeague *registrator.League
	}
	tests := []struct {
		name string
		args args
		want model.League
	}{
		{
			name: "test case 1",
			args: args{
				pbLeague: &registrator.League{
					Id:        1,
					Name:      "name",
					ShortName: "short_name",
					LogoLink:  "link",
					WebSite:   "site",
				},
			},
			want: model.League{
				ID:        1,
				Name:      "name",
				ShortName: "short_name",
				LogoLink:  "link",
				WebSite:   "site",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertPBLeagueToModelLeague(tt.args.pbLeague); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertPBLeagueToModelLeague() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertPBPlaceToModelPlace(t *testing.T) {
	type args struct {
		pbPlace *registrator.Place
	}
	tests := []struct {
		name string
		args args
		want model.Place
	}{
		{
			name: "test case 1",
			args: args{
				pbPlace: &registrator.Place{
					Id:        1,
					Address:   "address",
					Name:      "name",
					ShortName: "short_name",
					Latitude:  1.1,
					Longitude: 2.2,
					MenuLink:  "menu",
				},
			},
			want: model.Place{
				ID:        1,
				Address:   "address",
				Name:      "name",
				ShortName: "short_name",
				Latitude:  1.1,
				Longitude: 2.2,
				MenuLink:  "menu",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertPBPlaceToModelPlace(tt.args.pbPlace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertPBPlaceToModelPlace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertPBGameToModelGame(t *testing.T) {
	timeNow := time_utils.TimeNow()
	type args struct {
		pbGame *registrator.Game
	}
	tests := []struct {
		name string
		args args
		want model.Game
	}{
		{
			name: "test case 1",
			args: args{
				pbGame: &registrator.Game{
					Id:                  1,
					ExternalId:          2,
					LeagueId:            3,
					Type:                registrator.GameType_GAME_TYPE_CLASSIC,
					Number:              "1",
					Name:                "name",
					PlaceId:             4,
					Date:                timestamppb.New(timeNow),
					Price:               400,
					PaymentType:         "cash,card",
					MaxPlayers:          9,
					Payment:             registrator.Payment_PAYMENT_CASH,
					Registered:          true,
					My:                  true,
					NumberOfMyLegioners: 3,
					NumberOfLegioners:   4,
					NumberOfPlayers:     5,
					ResultPlace:         1,
				},
			},
			want: model.Game{
				ID:                  1,
				ExternalID:          2,
				Type:                1,
				Number:              "1",
				Name:                "name",
				Date:                model.DateTime(timestamppb.New(timeNow).AsTime()),
				Price:               400,
				PaymentType:         "cash,card",
				MaxPlayers:          9,
				Payment:             model.PaymentTypeCash,
				Registered:          true,
				My:                  true,
				NumberOfMyLegioners: 3,
				NumberOfLegioners:   4,
				NumberOfPlayers:     5,
				ResultPlace:         1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertPBGameToModelGame(tt.args.pbGame)
			assert.Equal(t, tt.want, got)
		})
	}
}
