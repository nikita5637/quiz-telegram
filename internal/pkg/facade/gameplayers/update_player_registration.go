package gameplayers

import (
	"context"
	"errors"
	"fmt"

	gameplayer "github.com/nikita5637/quiz-registrator-api/pkg/pb/game_player"
	gameplayerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/game_player"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/games"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	statusutils "github.com/nikita5637/quiz-telegram/internal/pkg/utils/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdatePlayerRegistration ...
func (f *Facade) UpdatePlayerRegistration(ctx context.Context, gamePlayer model.GamePlayer) error {
	getGamePlayersByGameIDResp, err := f.gamePlayerServiceClient.GetGamePlayersByGameID(ctx, &gameplayer.GetGamePlayersByGameIDRequest{
		GameId: gamePlayer.GameID,
	})
	if err != nil {
		return fmt.Errorf("get game players by game ID: %w", err)
	}

	id := int32(0)
	for _, pbGamePlayer := range getGamePlayersByGameIDResp.GetGamePlayers() {
		if pbGamePlayer.GetGameId() == gamePlayer.GameID &&
			pbGamePlayer.GetUserId().GetValue() == gamePlayer.UserID.Value() &&
			pbGamePlayer.GetRegisteredBy() == gamePlayer.RegisteredBy {
			id = pbGamePlayer.GetId()
			break
		}
	}
	if id == 0 {
		return errors.New("game player not found")
	}

	if _, err := f.gamePlayerRegistratorServiceClient.UpdatePlayerDegree(ctx, &gameplayerpb.UpdatePlayerDegreeRequest{
		Id:     id,
		Degree: gameplayerpb.Degree(gamePlayer.Degree),
	}); err != nil {
		st := status.Convert(err)
		errorInfo := statusutils.GetErrorInfoFromStatus(st)
		switch st.Code() {
		case codes.NotFound:
			return ErrGamePlayerNotFound
		case codes.FailedPrecondition:
			if errorInfo.Reason == games.ReasonGameNotFound {
				return games.ErrGameNotFound
			} else if errorInfo.Reason == ReasonThereAreNoRegistrationForTheGame {
				return ErrThereAreNoRegistrationForTheGame
			}

			return games.ErrGameHasPassed
		}

		return fmt.Errorf("patch game player error: %w", err)
	}

	return nil
}
