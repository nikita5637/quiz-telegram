package gameresults

import (
	"context"
	"fmt"

	"github.com/mono83/maybe"
	gameresultmanager "github.com/nikita5637/quiz-registrator-api/pkg/pb/game_result_manager"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
)

// GetGameResultByGameID ...
func (f *Facade) GetGameResultByGameID(ctx context.Context, gameID int32) (model.GameResult, error) {
	gameResult, err := f.gameResultManagerClient.SearchGameResultByGameID(ctx, &gameresultmanager.SearchGameResultByGameIDRequest{
		Id: gameID,
	})
	if err != nil {
		return model.GameResult{}, fmt.Errorf("SearchGameResultByGameID error: %w", err)
	}

	ret := model.GameResult{
		ID:          gameResult.GetId(),
		GameID:      gameResult.GetGameId(),
		ResultPlace: model.ResultPlace(gameResult.GetResultPlace()),
		RoundPoints: maybe.Nothing[string](),
	}

	if roundPoints := gameResult.GetRoundPoints(); roundPoints != "" {
		ret.RoundPoints = maybe.Just(roundPoints)
	}

	return ret, nil
}
