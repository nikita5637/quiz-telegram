package games

import (
	"context"
	"fmt"

	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
)

// GetGameByID ...
func (f *Facade) GetGameByID(ctx context.Context, gameID int32) (model.Game, error) {
	gameResp, err := f.registratorServiceClient.GetGameByID(ctx, &registrator.GetGameByIDRequest{
		GameId: gameID,
	})
	if err != nil {
		return model.Game{}, fmt.Errorf("get game by ID error: %w", err)
	}

	league, err := f.getModelLeague(ctx, gameResp.GetGame().GetLeagueId())
	if err != nil {
		return model.Game{}, fmt.Errorf("get league by ID error: %w", err)
	}

	place, err := f.getModelPlace(ctx, gameResp.GetGame().GetPlaceId())
	if err != nil {
		return model.Game{}, fmt.Errorf("get place by ID error: %w", err)
	}

	game := convertPBGameToModelGame(gameResp.GetGame())
	game.League = league
	game.Place = place

	return game, nil
}

// GetGames ...
func (f *Facade) GetGames(ctx context.Context, active bool) ([]model.Game, error) {
	gamesResp, err := f.registratorServiceClient.GetGames(ctx, &registrator.GetGamesRequest{
		Active: active,
	})
	if err != nil {
		return nil, fmt.Errorf("get games error: %w", err)
	}

	games := make([]model.Game, 0, len(gamesResp.GetGames()))
	for _, pbGame := range gamesResp.GetGames() {
		game := convertPBGameToModelGame(pbGame)

		league, err := f.getModelLeague(ctx, pbGame.GetLeagueId())
		if err != nil {
			return nil, fmt.Errorf("get league by ID error: %w", err)
		}

		game.League = league

		place, err := f.getModelPlace(ctx, pbGame.GetPlaceId())
		if err != nil {
			return nil, fmt.Errorf("get place by ID error: %w", err)
		}

		game.Place = place

		games = append(games, game)
	}

	return games, nil
}

// GetRegisteredGames ...
func (f *Facade) GetRegisteredGames(ctx context.Context, active bool) ([]model.Game, error) {
	gamesResp, err := f.registratorServiceClient.GetRegisteredGames(ctx, &registrator.GetRegisteredGamesRequest{
		Active: active,
	})
	if err != nil {
		return nil, fmt.Errorf("get registered games error: %w", err)
	}

	games := make([]model.Game, 0, len(gamesResp.GetGames()))
	for _, pbGame := range gamesResp.GetGames() {
		game := convertPBGameToModelGame(pbGame)

		league, err := f.getModelLeague(ctx, pbGame.GetLeagueId())
		if err != nil {
			return nil, fmt.Errorf("get league by ID error: %w", err)
		}

		game.League = league

		place, err := f.getModelPlace(ctx, pbGame.GetPlaceId())
		if err != nil {
			return nil, fmt.Errorf("get place by ID error: %w", err)
		}

		game.Place = place

		games = append(games, game)
	}

	return games, nil
}

func (f *Facade) getModelLeague(ctx context.Context, leagueID int32) (model.League, error) {
	if league, ok := f.leagueCache[leagueID]; ok {
		return league, nil
	}

	logger.DebugKV(ctx, "league not found in cache", "league ID", leagueID)

	leagueResp, err := f.registratorServiceClient.GetLeagueByID(ctx, &registrator.GetLeagueByIDRequest{
		Id: leagueID,
	})
	if err != nil {
		return model.League{}, fmt.Errorf("get league error: %w", err)
	}

	league := convertPBLeagueToModelLeague(leagueResp.GetLeague())
	f.leagueCache[leagueID] = league

	return league, nil
}

func (f *Facade) getModelPlace(ctx context.Context, placeID int32) (model.Place, error) {
	if place, ok := f.placeCache[placeID]; ok {
		return place, nil
	}

	logger.DebugKV(ctx, "place not found in cache", "place ID", placeID)

	placeResp, err := f.registratorServiceClient.GetPlaceByID(ctx, &registrator.GetPlaceByIDRequest{
		Id: placeID,
	})
	if err != nil {
		return model.Place{}, fmt.Errorf("get place error: %w", err)
	}

	place := convertPBPlaceToModelPlace(placeResp.GetPlace())
	f.placeCache[placeID] = place

	return place, nil
}

func convertPBLeagueToModelLeague(pbLeague *registrator.League) model.League {
	return model.League{
		ID:        pbLeague.GetId(),
		Name:      pbLeague.GetName(),
		ShortName: pbLeague.GetShortName(),
		LogoLink:  pbLeague.GetLogoLink(),
		WebSite:   pbLeague.GetWebSite(),
	}
}

func convertPBPlaceToModelPlace(pbPlace *registrator.Place) model.Place {
	return model.Place{
		ID:        pbPlace.GetId(),
		Address:   pbPlace.GetAddress(),
		Name:      pbPlace.GetName(),
		ShortName: pbPlace.GetShortName(),
		Longitude: pbPlace.GetLongitude(),
		Latitude:  pbPlace.GetLatitude(),
		MenuLink:  pbPlace.GetMenuLink(),
	}
}

func convertPBGameToModelGame(pbGame *registrator.Game) model.Game {
	return model.Game{
		ID:                  pbGame.GetId(),
		ExternalID:          pbGame.GetExternalId(),
		Type:                int32(pbGame.GetType()),
		Number:              pbGame.GetNumber(),
		Name:                pbGame.GetName(),
		Date:                model.DateTime(pbGame.GetDate().AsTime()),
		Price:               pbGame.GetPrice(),
		PaymentType:         pbGame.GetPaymentType(),
		MaxPlayers:          pbGame.GetMaxPlayers(),
		Payment:             model.PaymentType(pbGame.GetPayment()),
		Registered:          pbGame.GetRegistered(),
		My:                  pbGame.GetMy(),
		NumberOfMyLegioners: pbGame.GetNumberOfMyLegioners(),
		NumberOfLegioners:   pbGame.GetNumberOfLegioners(),
		NumberOfPlayers:     pbGame.GetNumberOfPlayers(),
		ResultPlace:         model.ResultPlace(pbGame.GetResultPlace()),
	}
}
