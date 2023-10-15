//go:generate mockery --case underscore --name LeagueServiceClient --with-expecter

package leagues

import (
	"context"

	leaguepb "github.com/nikita5637/quiz-registrator-api/pkg/pb/league"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"google.golang.org/grpc"
)

// LeagueServiceClient ...
type LeagueServiceClient interface {
	GetLeague(ctx context.Context, in *leaguepb.GetLeagueRequest, opts ...grpc.CallOption) (*leaguepb.League, error)
}

// Facade ...
type Facade struct {
	leaguesCache map[int32]model.League

	leagueServiceClient LeagueServiceClient
}

// Config ...
type Config struct {
	LeagueServiceClient LeagueServiceClient
}

// NewFacade ...
func NewFacade(cfg Config) *Facade {
	return &Facade{
		leaguesCache: make(map[int32]model.League, 0),

		leagueServiceClient: cfg.LeagueServiceClient,
	}
}
