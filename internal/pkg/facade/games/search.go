package games

import (
	"context"
	"fmt"

	gamepb "github.com/nikita5637/quiz-registrator-api/pkg/pb/game"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
)

// SearchPassedAndRegisteredGames ...
func (f *Facade) SearchPassedAndRegisteredGames(ctx context.Context, page, pageSize uint64) ([]model.Game, uint64, error) {
	passedGamesResp, err := f.gameServiceClient.SearchPassedAndRegisteredGames(ctx, &gamepb.SearchPassedAndRegisteredGamesRequest{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("searching passed and registered games error: %w", err)
	}

	modelGames := make([]model.Game, 0, len(passedGamesResp.GetGames()))
	for _, pbGame := range passedGamesResp.GetGames() {
		modelGames = append(modelGames, convertProtoGameToModelGame(pbGame))
	}

	return modelGames, passedGamesResp.GetTotal(), nil
}
