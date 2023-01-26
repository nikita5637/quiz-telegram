package games

import (
	"context"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
)

// UnregisterGame ...
func (f *Facade) UnregisterGame(ctx context.Context, gameID int32) (int32, error) {
	resp, err := f.registratorServiceClient.UnregisterGame(ctx, &registrator.UnregisterGameRequest{
		GameId: gameID,
	})
	if err != nil {
		return int32(registrator.UnregisterGameStatus_UNREGISTER_GAME_STATUS_INVALID), handleError(err)
	}

	return int32(resp.GetStatus()), nil
}
