package leagues

import (
	"context"
	"testing"

	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/leagues/mocks"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
)

type fixture struct {
	ctx    context.Context
	facade *Facade

	leagueServiceClient *mocks.LeagueServiceClient
}

func tearUp(t *testing.T) *fixture {
	fx := &fixture{
		ctx: context.Background(),

		leagueServiceClient: mocks.NewLeagueServiceClient(t),
	}

	fx.facade = &Facade{
		leaguesCache: make(map[int32]model.League, 0),

		leagueServiceClient: fx.leagueServiceClient,
	}

	t.Cleanup(func() {
	})

	return fx
}
