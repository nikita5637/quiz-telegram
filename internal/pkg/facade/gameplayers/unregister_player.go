package gameplayers

import (
	"context"

	gameplayerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/game_player"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/games"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	statusutils "github.com/nikita5637/quiz-telegram/internal/pkg/utils/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// UnregisterPlayer ...
func (f *Facade) UnregisterPlayer(ctx context.Context, gamePlayer model.GamePlayer) error {
	req := &gameplayerpb.UnregisterPlayerRequest{
		GamePlayer: &gameplayerpb.GamePlayer{
			GameId:       gamePlayer.GameID,
			RegisteredBy: gamePlayer.RegisteredBy,
			Degree:       gameplayerpb.Degree_DEGREE_LIKELY,
		},
	}
	if userID, ok := gamePlayer.UserID.Get(); ok {
		req.GamePlayer.UserId = wrapperspb.Int32(userID)
	}

	_, err := f.gamePlayerRegistratorServiceClient.UnregisterPlayer(ctx, req)
	if err != nil {
		st := status.Convert(err)
		errorInfo := statusutils.GetErrorInfoFromStatus(st)

		switch st.Code() {
		case codes.FailedPrecondition:
			if errorInfo != nil {
				if errorInfo.Reason == games.ReasonGameHasPassed {
					return games.ErrGameHasPassed
				} else if errorInfo.Reason == games.ReasonGameNotFound {
					return games.ErrGameNotFound
				}
			}
		case codes.NotFound:
			return ErrGamePlayerNotFound
		}

		return st.Err()
	}

	return nil
}
