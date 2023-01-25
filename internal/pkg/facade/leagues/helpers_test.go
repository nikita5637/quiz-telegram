package leagues

import (
	"context"
	"testing"

	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/games/mocks"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
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
		leaguesCache: make(map[int32]model.League, 0),

		registratorServiceClient: fx.registratorServiceClient,
	}

	t.Cleanup(func() {
	})

	return fx
}
