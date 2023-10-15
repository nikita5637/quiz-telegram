package gameplayers

import (
	"context"
	"testing"

	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/gameplayers/mocks"
)

type fixture struct {
	ctx context.Context

	gamePlayerServiceClient            *mocks.GamePlayerServiceClient
	gamePlayerRegistratorServiceClient *mocks.GamePlayerRegistratorServiceClient

	facade *Facade
}

func tearUp(t *testing.T) *fixture {
	fx := &fixture{
		ctx: context.Background(),

		gamePlayerServiceClient:            mocks.NewGamePlayerServiceClient(t),
		gamePlayerRegistratorServiceClient: mocks.NewGamePlayerRegistratorServiceClient(t),
	}

	fx.facade = New(Config{
		GamePlayerServiceClient:            fx.gamePlayerServiceClient,
		GamePlayerRegistratorServiceClient: fx.gamePlayerRegistratorServiceClient,
	})

	t.Cleanup(func() {})

	return fx
}
