//go:generate mockery --case underscore --name GamePlayerServiceClient --with-expecter
//go:generate mockery --case underscore --name GamePlayerRegistratorServiceClient --with-expecter
//go:generate mockery --case underscore --name RegistratorServiceClient --with-expecter

package gameplayers

import (
	"context"

	gameplayerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/game_player"
	registratorpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GamePlayerServiceClient ...
type GamePlayerServiceClient interface {
	DeleteGamePlayer(ctx context.Context, in *gameplayerpb.DeleteGamePlayerRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetGamePlayersByGameID(ctx context.Context, in *gameplayerpb.GetGamePlayersByGameIDRequest, opts ...grpc.CallOption) (*gameplayerpb.GetGamePlayersByGameIDResponse, error)
	PatchGamePlayer(ctx context.Context, in *gameplayerpb.PatchGamePlayerRequest, opts ...grpc.CallOption) (*gameplayerpb.GamePlayer, error)
}

// GamePlayerRegistratorServiceClient ...
type GamePlayerRegistratorServiceClient interface {
	RegisterPlayer(ctx context.Context, in *gameplayerpb.RegisterPlayerRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UnregisterPlayer(ctx context.Context, in *gameplayerpb.UnregisterPlayerRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

// RegistratorServiceClient ...
type RegistratorServiceClient interface {
	GetGameByID(ctx context.Context, in *registratorpb.GetGameByIDRequest, opts ...grpc.CallOption) (*registratorpb.GetGameByIDResponse, error)
}

// Facade ...
type Facade struct {
	gamePlayerServiceClient            GamePlayerServiceClient
	gamePlayerRegistratorServiceClient GamePlayerRegistratorServiceClient
	registratorServiceClient           RegistratorServiceClient
}

// Config ...
type Config struct {
	GamePlayerServiceClient            GamePlayerServiceClient
	GamePlayerRegistratorServiceClient GamePlayerRegistratorServiceClient
	RegistratorServiceClient           RegistratorServiceClient
}

// New ...
func New(cfg Config) *Facade {
	return &Facade{
		gamePlayerServiceClient:            cfg.GamePlayerServiceClient,
		gamePlayerRegistratorServiceClient: cfg.GamePlayerRegistratorServiceClient,
		registratorServiceClient:           cfg.RegistratorServiceClient,
	}
}
