package games

import (
	"errors"
	"testing"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestFacade_UpdatePayment(t *testing.T) {
	t.Run("error while update payment", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().UpdatePayment(fx.ctx, &registrator.UpdatePaymentRequest{
			GameId:  1,
			Payment: registrator.Payment_PAYMENT_CERTIFICATE,
		}).Return(nil, errors.New("some error"))

		err := fx.facade.UpdatePayment(fx.ctx, 1, int32(registrator.Payment_PAYMENT_CERTIFICATE))
		assert.Error(t, err)
	})

	t.Run("error game not found while update payment", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().UpdatePayment(fx.ctx, &registrator.UpdatePaymentRequest{
			GameId:  1,
			Payment: registrator.Payment_PAYMENT_CERTIFICATE,
		}).Return(nil, status.New(codes.NotFound, "").Err())

		err := fx.facade.UpdatePayment(fx.ctx, 1, int32(registrator.Payment_PAYMENT_CERTIFICATE))
		assert.Error(t, err)
		assert.ErrorIs(t, err, model.ErrGameNotFound)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.registratorServiceClient.EXPECT().UpdatePayment(fx.ctx, &registrator.UpdatePaymentRequest{
			GameId:  1,
			Payment: registrator.Payment_PAYMENT_CERTIFICATE,
		}).Return(&registrator.UpdatePaymentResponse{}, nil)

		err := fx.facade.UpdatePayment(fx.ctx, 1, int32(registrator.Payment_PAYMENT_CERTIFICATE))
		assert.NoError(t, err)
	})
}
