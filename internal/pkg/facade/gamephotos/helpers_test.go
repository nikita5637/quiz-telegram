package gamephotos

import (
	"context"
	"testing"

	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/gamephotos/mocks"
)

type fixture struct {
	ctx    context.Context
	facade *Facade

	photographerServiceClient *mocks.PhotographerServiceClient
}

func tearUp(t *testing.T) *fixture {
	fx := &fixture{
		ctx: context.Background(),

		photographerServiceClient: mocks.NewPhotographerServiceClient(t),
	}

	fx.facade = &Facade{
		photographerServiceClient: fx.photographerServiceClient,
	}

	t.Cleanup(func() {
	})

	return fx
}
