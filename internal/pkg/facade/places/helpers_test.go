package places

import (
	"context"
	"testing"

	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/places/mocks"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
)

type fixture struct {
	ctx    context.Context
	facade *Facade

	placeServiceClient *mocks.PlaceServiceClient
}

func tearUp(t *testing.T) *fixture {
	fx := &fixture{
		ctx: context.Background(),

		placeServiceClient: mocks.NewPlaceServiceClient(t),
	}

	fx.facade = &Facade{
		placesCache: make(map[int32]model.Place, 0),

		placeServiceClient: fx.placeServiceClient,
	}

	t.Cleanup(func() {})

	return fx
}
