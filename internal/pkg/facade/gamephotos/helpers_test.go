package gamephotos

import (
	"context"
	"testing"

	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/gamephotos/mocks"
)

type fixture struct {
	ctx    context.Context
	facade *Facade

	leaguesFacade *mocks.LeaguesFacade
	placesFacade  *mocks.PlacesFacade

	photographerServiceClient *mocks.PhotographerServiceClient
	registratorServiceClient  *mocks.RegistratorServiceClient
}

func tearUp(t *testing.T) *fixture {
	fx := &fixture{
		ctx: context.Background(),

		leaguesFacade: mocks.NewLeaguesFacade(t),
		placesFacade:  mocks.NewPlacesFacade(t),

		photographerServiceClient: mocks.NewPhotographerServiceClient(t),
		registratorServiceClient:  mocks.NewRegistratorServiceClient(t),
	}

	fx.facade = &Facade{
		leaguesFacade: fx.leaguesFacade,
		placesFacade:  fx.placesFacade,

		photographerServiceClient: fx.photographerServiceClient,
		registratorServiceClient:  fx.registratorServiceClient,
	}

	t.Cleanup(func() {
	})

	return fx
}
