//go:generate mockery --case underscore --name RegistratorServiceClient --with-expecter

package games

import (
	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
)

// RegistratorServiceClient ...
type RegistratorServiceClient interface {
	registrator.RegistratorServiceClient
}

// Facade ...
type Facade struct {
	leagueCache              map[int32]model.League
	placeCache               map[int32]model.Place
	registratorServiceClient RegistratorServiceClient
}

// Config ...
type Config struct {
	RegistratorServiceClient RegistratorServiceClient
}

// NewFacade ...
func NewFacade(cfg Config) *Facade {
	return &Facade{
		leagueCache:              make(map[int32]model.League, 0),
		placeCache:               make(map[int32]model.Place, 0),
		registratorServiceClient: cfg.RegistratorServiceClient,
	}
}
