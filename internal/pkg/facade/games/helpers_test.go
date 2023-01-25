package games

import (
	"context"
	"testing"

	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/games/mocks"
)

type fixture struct {
	ctx    context.Context
	facade *Facade

	leaguesFacade *mocks.LeaguesFacade
	placesFacade  *mocks.PlacesFacade

	registratorServiceClient *mocks.RegistratorServiceClient
}

func tearUp(t *testing.T) *fixture {
	fx := &fixture{
		ctx: context.Background(),

		leaguesFacade: mocks.NewLeaguesFacade(t),
		placesFacade:  mocks.NewPlacesFacade(t),

		registratorServiceClient: mocks.NewRegistratorServiceClient(t),
	}

	fx.facade = &Facade{
		leaguesFacade: fx.leaguesFacade,
		placesFacade:  fx.placesFacade,

		registratorServiceClient: fx.registratorServiceClient,
	}

	t.Cleanup(func() {
	})

	return fx
}
