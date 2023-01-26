package games

import (
	"context"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
)

// UpdatePayment ...
func (f *Facade) UpdatePayment(ctx context.Context, gameID, payment int32) error {
	_, err := f.registratorServiceClient.UpdatePayment(ctx, &registrator.UpdatePaymentRequest{
		GameId:  gameID,
		Payment: registrator.Payment(payment),
	})
	if err != nil {
		return handleError(err)
	}

	return nil
}
