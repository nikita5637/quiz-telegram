//go:generate mockery --case underscore --name GamePlayerServiceClient --with-expecter
//go:generate mockery --case underscore --name GamePlayerRegistratorServiceClient --with-expecter

package gameplayers

import (
	"context"

	gameplayerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/game_player"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GamePlayerServiceClient ...
type GamePlayerServiceClient interface {
	GetGamePlayersByGameID(ctx context.Context, in *gameplayerpb.GetGamePlayersByGameIDRequest, opts ...grpc.CallOption) (*gameplayerpb.GetGamePlayersByGameIDResponse, error)
	GetUserGameIDs(ctx context.Context, in *gameplayerpb.GetUserGameIDsRequest, opts ...grpc.CallOption) (*gameplayerpb.GetUserGameIDsResponse, error)
}

// GamePlayerRegistratorServiceClient ...
type GamePlayerRegistratorServiceClient interface {
	RegisterPlayer(ctx context.Context, in *gameplayerpb.RegisterPlayerRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UnregisterPlayer(ctx context.Context, in *gameplayerpb.UnregisterPlayerRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UpdatePlayerDegree(ctx context.Context, in *gameplayerpb.UpdatePlayerDegreeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

// Facade ...
type Facade struct {
	gamePlayerServiceClient            GamePlayerServiceClient
	gamePlayerRegistratorServiceClient GamePlayerRegistratorServiceClient
}

// Config ...
type Config struct {
	GamePlayerServiceClient            GamePlayerServiceClient
	GamePlayerRegistratorServiceClient GamePlayerRegistratorServiceClient
}

// New ...
func New(cfg Config) *Facade {
	return &Facade{
		gamePlayerServiceClient:            cfg.GamePlayerServiceClient,
		gamePlayerRegistratorServiceClient: cfg.GamePlayerRegistratorServiceClient,
	}
}
