package games

import (
	"context"
	"fmt"

	"github.com/nikita5637/quiz-telegram/internal/pkg/model"

	commonpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/common"
	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	converter_utils "github.com/nikita5637/quiz-telegram/utils/converter"
)

// GetGameByID ...
func (f *Facade) GetGameByID(ctx context.Context, id int32) (model.Game, error) {
	gameResp, err := f.registratorServiceClient.GetGameByID(ctx, &registrator.GetGameByIDRequest{
		GameId: id,
	})
	if err != nil {
		return model.Game{}, fmt.Errorf("get game by ID error: %w", handleError(err))
	}

	league, err := f.leaguesFacade.GetLeagueByID(ctx, gameResp.GetGame().GetLeagueId())
	if err != nil {
		return model.Game{}, fmt.Errorf("get league by ID error: %w", err)
	}

	place, err := f.placesFacade.GetPlaceByID(ctx, gameResp.GetGame().GetPlaceId())
	if err != nil {
		return model.Game{}, fmt.Errorf("get place by ID error: %w", err)
	}

	game := converter_utils.ConvertPBGameToModelGame(gameResp.GetGame())
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

	return f.getGames(ctx, gamesResp.GetGames())
}

// GetRegisteredGames ...
func (f *Facade) GetRegisteredGames(ctx context.Context, active bool) ([]model.Game, error) {
	gamesResp, err := f.registratorServiceClient.GetRegisteredGames(ctx, &registrator.GetRegisteredGamesRequest{
		Active: active,
	})
	if err != nil {
		return nil, fmt.Errorf("get registered games error: %w", err)
	}

	return f.getGames(ctx, gamesResp.GetGames())
}

// GetUserGames ...
func (f *Facade) GetUserGames(ctx context.Context, active bool, userID int32) ([]model.Game, error) {
	gamesResp, err := f.registratorServiceClient.GetUserGames(ctx, &registrator.GetUserGamesRequest{
		Active: active,
		UserId: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("get user games error: %w", err)
	}

	return f.getGames(ctx, gamesResp.GetGames())
}

func (f *Facade) getGames(ctx context.Context, pbGames []*commonpb.Game) ([]model.Game, error) {
	games := make([]model.Game, 0, len(pbGames))
	for _, pbGame := range pbGames {
		game := converter_utils.ConvertPBGameToModelGame(pbGame)

		league, err := f.leaguesFacade.GetLeagueByID(ctx, pbGame.GetLeagueId())
		if err != nil {
			return nil, fmt.Errorf("get league by ID error: %w", err)
		}

		game.League = league

		place, err := f.placesFacade.GetPlaceByID(ctx, pbGame.GetPlaceId())
		if err != nil {
			return nil, fmt.Errorf("get place by ID error: %w", err)
		}

		game.Place = place

		games = append(games, game)
	}

	return games, nil
}
