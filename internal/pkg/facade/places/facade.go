//go:generate mockery --case underscore --name PlaceServiceClient --with-expecter

package places

import (
	"context"

	placepb "github.com/nikita5637/quiz-registrator-api/pkg/pb/place"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"google.golang.org/grpc"
)

// PlaceServiceClient ...
type PlaceServiceClient interface {
	GetPlace(ctx context.Context, in *placepb.GetPlaceRequest, opts ...grpc.CallOption) (*placepb.Place, error)
}

// Facade ...
type Facade struct {
	placesCache map[int32]model.Place

	placeServiceClient PlaceServiceClient
}

// Config ...
type Config struct {
	PlaceServiceClient PlaceServiceClient
}

// NewFacade ...
func NewFacade(cfg Config) *Facade {
	return &Facade{
		placesCache: make(map[int32]model.Place, 0),

		placeServiceClient: cfg.PlaceServiceClient,
	}
}
