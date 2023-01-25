//go:generate mockery --case underscore --name RegistratorServiceClient --with-expecter

package places

import (
	"context"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"google.golang.org/grpc"
)

// RegistratorServiceClient ...
type RegistratorServiceClient interface {
	GetPlaceByID(ctx context.Context, in *registrator.GetPlaceByIDRequest, opts ...grpc.CallOption) (*registrator.GetPlaceByIDResponse, error)
}

// Facade ...
type Facade struct {
	placesCache map[int32]model.Place

	registratorServiceClient RegistratorServiceClient
}

// Config ...
type Config struct {
	RegistratorServiceClient RegistratorServiceClient
}

// NewFacade ...
func NewFacade(cfg Config) *Facade {
	return &Facade{
		placesCache: make(map[int32]model.Place, 0),

		registratorServiceClient: cfg.RegistratorServiceClient,
	}
}
