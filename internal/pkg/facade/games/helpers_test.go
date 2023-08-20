package games

import (
	"context"
	"errors"
	"testing"

	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/games/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func Test_handleError(t *testing.T) {
	t.Run("error is nil", func(t *testing.T) {
		err := handleError(nil)
		assert.Nil(t, err)
	})

	t.Run("error is not found", func(t *testing.T) {
		err := handleError(status.New(codes.NotFound, "").Err())
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrGameNotFound)
	})

	t.Run("otherwise case", func(t *testing.T) {
		err := handleError(errors.New("some error"))
		assert.Error(t, err)
	})
}
