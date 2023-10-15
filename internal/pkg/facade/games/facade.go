//go:generate mockery --case underscore --name GamePlayersFacade --with-expecter
//go:generate mockery --case underscore --name GameServiceClient --with-expecter
//go:generate mockery --case underscore --name GameRegistratorServiceClient --with-expecter

package games

import (
	"context"

	gamepb "github.com/nikita5637/quiz-registrator-api/pkg/pb/game"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GamePlayersFacade ...
type GamePlayersFacade interface {
	GetUserGameIDs(ctx context.Context, userID int32) ([]int32, error)
}

// GameServiceClient ...
type GameServiceClient interface {
	BatchGetGames(ctx context.Context, in *gamepb.BatchGetGamesRequest, opts ...grpc.CallOption) (*gamepb.BatchGetGamesResponse, error)
	GetGame(ctx context.Context, in *gamepb.GetGameRequest, opts ...grpc.CallOption) (*gamepb.Game, error)
	ListGames(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*gamepb.ListGamesResponse, error)
	SearchPassedAndRegisteredGames(ctx context.Context, in *gamepb.SearchPassedAndRegisteredGamesRequest, opts ...grpc.CallOption) (*gamepb.SearchPassedAndRegisteredGamesResponse, error)
}

// GameRegistratorServiceClient ...
type GameRegistratorServiceClient interface {
	RegisterGame(ctx context.Context, in *gamepb.RegisterGameRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UnregisterGame(ctx context.Context, in *gamepb.UnregisterGameRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UpdatePayment(ctx context.Context, in *gamepb.UpdatePaymentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

// Facade ...
type Facade struct {
	gamePlayersFacade GamePlayersFacade

	gameServiceClient            GameServiceClient
	gameRegistratorServiceClient GameRegistratorServiceClient
}

// Config ...
type Config struct {
	GamePlayersFacade GamePlayersFacade

	GameServiceClient            GameServiceClient
	GameRegistratorServiceClient GameRegistratorServiceClient
}

// New ...
func New(cfg Config) *Facade {
	return &Facade{
		gamePlayersFacade: cfg.GamePlayersFacade,

		gameServiceClient:            cfg.GameServiceClient,
		gameRegistratorServiceClient: cfg.GameRegistratorServiceClient,
	}
}
