//go:generate mockery --case underscore --name LeaguesFacade --with-expecter
//go:generate mockery --case underscore --name PlacesFacade --with-expecter
//go:generate mockery --case underscore --name RegistratorServiceClient --with-expecter

package games

import (
	"context"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
)

// LeaguesFacade ...
type LeaguesFacade interface {
	GetLeagueByID(ctx context.Context, leagueID int32) (model.League, error)
}

// PlacesFacade ...
type PlacesFacade interface {
	GetPlaceByID(ctx context.Context, placeID int32) (model.Place, error)
}

// RegistratorServiceClient ...
type RegistratorServiceClient interface {
	registrator.RegistratorServiceClient
}

// Facade ...
type Facade struct {
	leaguesFacade LeaguesFacade
	placesFacade  PlacesFacade

	registratorServiceClient RegistratorServiceClient
}

// Config ...
type Config struct {
	LeaguesFacade LeaguesFacade
	PlacesFacade  PlacesFacade

	RegistratorServiceClient RegistratorServiceClient
}

// NewFacade ...
func NewFacade(cfg Config) *Facade {
	return &Facade{
		leaguesFacade: cfg.LeaguesFacade,
		placesFacade:  cfg.PlacesFacade,

		registratorServiceClient: cfg.RegistratorServiceClient,
	}
}
