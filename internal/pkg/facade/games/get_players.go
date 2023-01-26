package games

import (
	"context"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
)

// GetPlayersByGameID ...
func (f *Facade) GetPlayersByGameID(ctx context.Context, gameID int32) ([]model.Player, error) {
	resp, err := f.registratorServiceClient.GetPlayersByGameID(ctx, &registrator.GetPlayersByGameIDRequest{
		GameId: gameID,
	})
	if err != nil {
		return nil, handleError(err)
	}

	modelPlayers := make([]model.Player, 0, len(resp.GetPlayers()))
	for _, pbPlayer := range resp.GetPlayers() {
		modelPlayers = append(modelPlayers, model.Player{
			UserID:       pbPlayer.GetUserId(),
			RegisteredBy: pbPlayer.GetRegisteredBy(),
			Degree:       int32(pbPlayer.GetDegree()),
		})
	}

	return modelPlayers, nil
}
