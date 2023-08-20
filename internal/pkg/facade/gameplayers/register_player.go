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

const ()

// RegisterPlayer ...
func (f *Facade) RegisterPlayer(ctx context.Context, gamePlayer model.GamePlayer) error {
	_, err := f.gamePlayerRegistratorServiceClient.RegisterPlayer(ctx, &gameplayerpb.RegisterPlayerRequest{
		GamePlayer: convertModelGamePlayerToProtoGamePlayer(gamePlayer),
	})

	return handleRegisterPlayerError(err)
}

func convertModelGamePlayerToProtoGamePlayer(gamePlayer model.GamePlayer) *gameplayerpb.GamePlayer {
	ret := &gameplayerpb.GamePlayer{
		Id:           gamePlayer.ID,
		GameId:       gamePlayer.GameID,
		RegisteredBy: gamePlayer.RegisteredBy,
		Degree:       gameplayerpb.Degree(gamePlayer.Degree),
	}

	if userID, ok := gamePlayer.UserID.Get(); ok {
		ret.UserId = &wrapperspb.Int32Value{
			Value: userID,
		}
	}

	return ret
}

func handleRegisterPlayerError(err error) error {
	if err == nil {
		return nil
	}

	st := status.Convert(err)
	errorInfo := statusutils.GetErrorInfoFromStatus(st)
	switch st.Code() {
	case codes.FailedPrecondition:
		switch errorInfo.Reason {
		case games.ReasonGameHasPassed:
			return games.ErrGameHasPassed
		case games.ReasonGameNotFound:
			return games.ErrGameNotFound
		case ReasonNoFreeSlot:
			return ErrNoFreeSlot
		}
	case codes.AlreadyExists:
		return ErrGamePlayerAlreadyRegistered
	}

	return err
}
