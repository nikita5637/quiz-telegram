package users

import (
	"context"
	"testing"

	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/users/mocks"
)

type fixture struct {
	ctx    context.Context
	facade *Facade

	registratorServiceClient *mocks.RegistratorServiceClient
}

func tearUp(t *testing.T) *fixture {
	fx := &fixture{
		ctx: context.Background(),

		registratorServiceClient: mocks.NewRegistratorServiceClient(t),
	}

	fx.facade = &Facade{
		registratorServiceClient: fx.registratorServiceClient,
	}

	t.Cleanup(func() {
	})

	return fx
}
