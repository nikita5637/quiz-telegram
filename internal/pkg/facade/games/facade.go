//go:generate mockery --case underscore --name LeaguesFacade --with-expecter
//go:generate mockery --case underscore --name PlacesFacade --with-expecter
//go:generate mockery --case underscore --name RegistratorServiceClient --with-expecter

package games

import (
	"context"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"google.golang.org/grpc"
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
	GetGameByID(ctx context.Context, in *registrator.GetGameByIDRequest, opts ...grpc.CallOption) (*registrator.GetGameByIDResponse, error)
	GetGames(ctx context.Context, in *registrator.GetGamesRequest, opts ...grpc.CallOption) (*registrator.GetGamesResponse, error)
	GetRegisteredGames(ctx context.Context, in *registrator.GetRegisteredGamesRequest, opts ...grpc.CallOption) (*registrator.GetRegisteredGamesResponse, error)
	GetUserGames(ctx context.Context, in *registrator.GetUserGamesRequest, opts ...grpc.CallOption) (*registrator.GetUserGamesResponse, error)
	RegisterGame(ctx context.Context, in *registrator.RegisterGameRequest, opts ...grpc.CallOption) (*registrator.RegisterGameResponse, error)
	UnregisterGame(ctx context.Context, in *registrator.UnregisterGameRequest, opts ...grpc.CallOption) (*registrator.UnregisterGameResponse, error)
	UpdatePayment(ctx context.Context, in *registrator.UpdatePaymentRequest, opts ...grpc.CallOption) (*registrator.UpdatePaymentResponse, error)
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
