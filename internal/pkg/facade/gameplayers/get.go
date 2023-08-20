package gameplayers

import (
	"context"
	"fmt"

	"github.com/mono83/maybe"
	gameplayerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/game_player"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
)

// GetGamePlayersByGameID ...
func (f *Facade) GetGamePlayersByGameID(ctx context.Context, gameID int32) ([]model.GamePlayer, error) {
	resp, err := f.gamePlayerServiceClient.GetGamePlayersByGameID(ctx, &gameplayerpb.GetGamePlayersByGameIDRequest{
		GameId: gameID,
	})
	if err != nil {
		return nil, fmt.Errorf("get game players by game ID error: %w", err)
	}

	gamePlayers := make([]model.GamePlayer, 0, len(resp.GetGamePlayers()))
	for _, pbGamePlayer := range resp.GetGamePlayers() {
		modelGamePlayer := model.GamePlayer{
			ID:           pbGamePlayer.GetId(),
			GameID:       pbGamePlayer.GetGameId(),
			UserID:       maybe.Nothing[int32](),
			RegisteredBy: pbGamePlayer.GetRegisteredBy(),
			Degree:       model.Degree(pbGamePlayer.GetDegree()),
		}

		if pbGamePlayer.GetUserId() != nil {
			modelGamePlayer.UserID = maybe.Just(pbGamePlayer.GetUserId().GetValue())
		}

		gamePlayers = append(gamePlayers, modelGamePlayer)
	}

	return gamePlayers, nil
}
