//go:generate mockery --case underscore --name RegistratorServiceClient --with-expecter

package leagues

import (
	"context"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"google.golang.org/grpc"
)

// RegistratorServiceClient ...
type RegistratorServiceClient interface {
	GetLeagueByID(ctx context.Context, in *registrator.GetLeagueByIDRequest, opts ...grpc.CallOption) (*registrator.GetLeagueByIDResponse, error)
}

// Facade ...
type Facade struct {
	leaguesCache map[int32]model.League

	registratorServiceClient RegistratorServiceClient
}

// Config ...
type Config struct {
	RegistratorServiceClient RegistratorServiceClient
}

// NewFacade ...
func NewFacade(cfg Config) *Facade {
	return &Facade{
		leaguesCache: make(map[int32]model.League, 0),

		registratorServiceClient: cfg.RegistratorServiceClient,
	}
}
