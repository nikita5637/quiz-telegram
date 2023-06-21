package gamephotos

import (
	"context"
	"fmt"

	commonpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/common"
	photomanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/photo_manager"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	converter_utils "github.com/nikita5637/quiz-telegram/utils/converter"
)

// GetGamesWithPhotos ...
func (f *Facade) GetGamesWithPhotos(ctx context.Context, limit, offset uint32) ([]model.Game, uint32, error) {
	gamesResp, err := f.photographerServiceClient.GetGamesWithPhotos(ctx, &photomanagerpb.GetGamesWithPhotosRequest{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("get games with photos: %w", err)
	}

	games, err := f.getGames(ctx, gamesResp.GetGames())
	if err != nil {
		return nil, 0, err
	}

	return games, gamesResp.GetTotal(), nil
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
