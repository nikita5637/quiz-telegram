//go:generate mockery --case underscore --name GameResultManagerClient --with-expecter

package gameresults

import (
	"context"

	gameresultmanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/game_result_manager"
	"google.golang.org/grpc"
)

// GameResultManagerClient ...
type GameResultManagerClient interface {
	SearchGameResultByGameID(ctx context.Context, in *gameresultmanagerpb.SearchGameResultByGameIDRequest, opts ...grpc.CallOption) (*gameresultmanagerpb.GameResult, error)
}

// Facade ...
type Facade struct {
	gameResultManagerClient GameResultManagerClient
}

// Config ...
type Config struct {
	GameResultManagerClient GameResultManagerClient
}

// New ...
func New(cfg Config) *Facade {
	return &Facade{
		gameResultManagerClient: cfg.GameResultManagerClient,
	}
}
