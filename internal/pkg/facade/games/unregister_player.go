package games

import (
	"context"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
)

// UnregisterPlayer ...
func (f *Facade) UnregisterPlayer(ctx context.Context, gameID, playerType int32) (int32, error) {
	resp, err := f.registratorServiceClient.UnregisterPlayer(ctx, &registrator.UnregisterPlayerRequest{
		GameId:     gameID,
		PlayerType: registrator.PlayerType(playerType),
	})
	if err != nil {
		return int32(registrator.UnregisterPlayerStatus_UNREGISTER_PLAYER_STATUS_INVALID), handleError(err)
	}

	return int32(resp.GetStatus()), nil
}
