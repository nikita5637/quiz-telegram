package games

import (
	"context"
	"testing"

	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/games/mocks"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
)

type fixture struct {
	ctx    context.Context
	facade *Facade

	photographerServiceClient *mocks.PhotographerServiceClient
	registratorServiceClient  *mocks.RegistratorServiceClient
}

func tearUp(t *testing.T) *fixture {
	fx := &fixture{
		ctx: context.Background(),

		photographerServiceClient: mocks.NewPhotographerServiceClient(t),
		registratorServiceClient:  mocks.NewRegistratorServiceClient(t),
	}

	fx.facade = &Facade{
		leagueCache: make(map[int32]model.League, 0),
		placeCache:  make(map[int32]model.Place, 0),

		photographerServiceClient: fx.photographerServiceClient,
		registratorServiceClient:  fx.registratorServiceClient,
	}

	t.Cleanup(func() {
	})

	return fx
}
