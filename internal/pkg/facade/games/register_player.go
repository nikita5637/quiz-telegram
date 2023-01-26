package games

import (
	"context"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RegisterPlayer ...
func (f *Facade) RegisterPlayer(ctx context.Context, gameID, playerType, degree int32) (int32, error) {
	resp, err := f.registratorServiceClient.RegisterPlayer(ctx, &registrator.RegisterPlayerRequest{
		GameId:     gameID,
		PlayerType: registrator.PlayerType(playerType),
		Degree:     registrator.Degree(degree),
	})
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.NotFound {
			return int32(registrator.RegisterPlayerStatus_REGISTER_PLAYER_STATUS_INVALID), model.ErrGameNotFound
		} else if st.Code() == codes.AlreadyExists {
			return int32(registrator.RegisterPlayerStatus_REGISTER_PLAYER_STATUS_INVALID), model.ErrNoFreeSlot
		}
		return int32(registrator.RegisterPlayerStatus_REGISTER_PLAYER_STATUS_INVALID), err
	}

	return int32(resp.GetStatus()), nil
}
