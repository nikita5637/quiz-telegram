package games

import (
	"context"
	"fmt"

	gamepb "github.com/nikita5637/quiz-registrator-api/pkg/pb/game"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RegisterGame ...
func (f *Facade) RegisterGame(ctx context.Context, gameID int32) error {
	if _, err := f.gameRegistratorServiceClient.RegisterGame(ctx, &gamepb.RegisterGameRequest{
		Id: gameID,
	}); err != nil {
		st := status.Convert(err)
		if st.Code() == codes.NotFound {
			return ErrGameNotFound
		}

		return fmt.Errorf("registering game error: %w", err)
	}

	return nil
}
