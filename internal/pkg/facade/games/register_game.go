package games

import (
	"context"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
)

// RegisterGame ...
func (f *Facade) RegisterGame(ctx context.Context, gameID int32) (int32, error) {
	resp, err := f.registratorServiceClient.RegisterGame(ctx, &registrator.RegisterGameRequest{
		GameId: gameID,
	})
	if err != nil {
		return int32(registrator.RegisterGameStatus_REGISTER_GAME_STATUS_INVALID), handleError(err)
	}

	return int32(resp.GetStatus()), nil
}
