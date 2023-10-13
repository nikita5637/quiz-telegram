package games

import (
	"context"
	"fmt"

	gamepb "github.com/nikita5637/quiz-registrator-api/pkg/pb/game"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdatePayment ...
func (f *Facade) UpdatePayment(ctx context.Context, gameID, payment int32) error {
	_, err := f.gameRegistratorServiceClient.UpdatePayment(ctx, &gamepb.UpdatePaymentRequest{
		Id:      gameID,
		Payment: gamepb.Payment(payment),
	})
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.NotFound {
			return ErrGameNotFound
		}

		return fmt.Errorf("updating payment error: %w", err)
	}

	return nil
}
