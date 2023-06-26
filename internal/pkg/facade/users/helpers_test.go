package users

import (
	"context"
	"testing"

	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/users/mocks"
)

type fixture struct {
	ctx    context.Context
	facade *Facade

	userManagerServiceClient *mocks.UserManagerServiceClient
}

func tearUp(t *testing.T) *fixture {
	fx := &fixture{
		ctx: context.Background(),

		userManagerServiceClient: mocks.NewUserManagerServiceClient(t),
	}

	fx.facade = &Facade{
		userManagerServiceClient: fx.userManagerServiceClient,
	}

	t.Cleanup(func() {
	})

	return fx
}
