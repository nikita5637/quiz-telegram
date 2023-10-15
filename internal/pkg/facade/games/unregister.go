package games

import (
	"context"
	"fmt"

	gamepb "github.com/nikita5637/quiz-registrator-api/pkg/pb/game"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnregisterGame ...
func (f *Facade) UnregisterGame(ctx context.Context, gameID int32) error {
	if _, err := f.gameRegistratorServiceClient.UnregisterGame(ctx, &gamepb.UnregisterGameRequest{
		Id: gameID,
	}); err != nil {
		st := status.Convert(err)
		if st.Code() == codes.NotFound {
			return ErrGameNotFound
		}

		return fmt.Errorf("unregistering game error: %w", err)
	}

	return nil
}
