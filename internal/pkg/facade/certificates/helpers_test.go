package certificates

import (
	"context"
	"testing"

	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/certificates/mocks"
)

type fixture struct {
	ctx context.Context

	certificateManagerServiceClient *mocks.CertificateManagerServiceClient

	facade *Facade
}

func tearUp(t *testing.T) *fixture {
	fx := &fixture{
		ctx: context.Background(),

		certificateManagerServiceClient: mocks.NewCertificateManagerServiceClient(t),
	}

	fx.facade = New(Config{
		CertificateManagerServiceClient: fx.certificateManagerServiceClient,
	})

	t.Cleanup(func() {})

	return fx
}
