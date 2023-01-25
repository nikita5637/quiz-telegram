package leagues

import (
	"context"
	"fmt"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetLeagueByID ...
func (f *Facade) GetLeagueByID(ctx context.Context, leagueID int32) (model.League, error) {
	if league, ok := f.leaguesCache[leagueID]; ok {
		return league, nil
	}

	logger.DebugKV(ctx, "league not found in cache", "league ID", leagueID)

	leagueResp, err := f.registratorServiceClient.GetLeagueByID(ctx, &registrator.GetLeagueByIDRequest{
		Id: leagueID,
	})
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.NotFound {
			return model.League{}, model.ErrLeagueNotFound
		}

		return model.League{}, fmt.Errorf("get league error: %w", err)
	}

	league := convertPBLeagueToModelLeague(leagueResp.GetLeague())
	f.leaguesCache[leagueID] = league

	return league, nil
}

func convertPBLeagueToModelLeague(pbLeague *registrator.League) model.League {
	return model.League{
		ID:        pbLeague.GetId(),
		Name:      pbLeague.GetName(),
		ShortName: pbLeague.GetShortName(),
		LogoLink:  pbLeague.GetLogoLink(),
		WebSite:   pbLeague.GetWebSite(),
	}
}
