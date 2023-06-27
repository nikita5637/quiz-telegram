//go:generate mockery --case underscore --name LeaguesFacade --with-expecter
//go:generate mockery --case underscore --name PlacesFacade --with-expecter
//go:generate mockery --case underscore --name PhotographerServiceClient --with-expecter

package gamephotos

import (
	"context"

	photomanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/photo_manager"
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

// PhotographerServiceClient ...
type PhotographerServiceClient interface {
	photomanagerpb.ServiceClient
}

// Facade ...
type Facade struct {
	leaguesFacade LeaguesFacade
	placesFacade  PlacesFacade

	photographerServiceClient PhotographerServiceClient
}

// Config ...
type Config struct {
	LeaguesFacade LeaguesFacade
	PlacesFacade  PlacesFacade

	PhotographerServiceClient PhotographerServiceClient
}

// NewFacade ...
func NewFacade(cfg Config) *Facade {
	return &Facade{
		leaguesFacade: cfg.LeaguesFacade,
		placesFacade:  cfg.PlacesFacade,

		photographerServiceClient: cfg.PhotographerServiceClient,
	}
}
