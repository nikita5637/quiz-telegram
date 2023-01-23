//go:generate mockery --case underscore --name PhotographerServiceClient --with-expecter
//go:generate mockery --case underscore --name RegistratorServiceClient --with-expecter

package games

import (
	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
)

// PhotographerServiceClient ...
type PhotographerServiceClient interface {
	registrator.PhotographerServiceClient
}

// RegistratorServiceClient ...
type RegistratorServiceClient interface {
	registrator.RegistratorServiceClient
}

// Facade ...
type Facade struct {
	leagueCache map[int32]model.League
	placeCache  map[int32]model.Place

	photographerServiceClient PhotographerServiceClient
	registratorServiceClient  RegistratorServiceClient
}

// Config ...
type Config struct {
	PhotographerServiceClient PhotographerServiceClient
	RegistratorServiceClient  RegistratorServiceClient
}

// NewFacade ...
func NewFacade(cfg Config) *Facade {
	return &Facade{
		leagueCache: make(map[int32]model.League, 0),
		placeCache:  make(map[int32]model.Place, 0),

		photographerServiceClient: cfg.PhotographerServiceClient,
		registratorServiceClient:  cfg.RegistratorServiceClient,
	}
}
