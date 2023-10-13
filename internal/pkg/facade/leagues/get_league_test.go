package leagues

import (
	"errors"
	"reflect"
	"testing"

	leaguepb "github.com/nikita5637/quiz-registrator-api/pkg/pb/league"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestFacade_GetLeague(t *testing.T) {
	t.Run("error league not found while get league", func(t *testing.T) {
		fx := tearUp(t)

		fx.leagueServiceClient.EXPECT().GetLeague(fx.ctx, &leaguepb.GetLeagueRequest{
			Id: 1,
		}).Return(nil, status.New(codes.NotFound, "").Err())

		got, err := fx.facade.GetLeague(fx.ctx, 1)
		assert.Equal(t, model.League{}, got)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrLeagueNotFound)
	})

	t.Run("error while get league", func(t *testing.T) {
		fx := tearUp(t)

		fx.leagueServiceClient.EXPECT().GetLeague(fx.ctx, &leaguepb.GetLeagueRequest{
			Id: 1,
		}).Return(nil, errors.New("some error"))

		got, err := fx.facade.GetLeague(fx.ctx, 1)
		assert.Equal(t, model.League{}, got)
		assert.Error(t, err)
	})

	t.Run("ok without cache", func(t *testing.T) {
		fx := tearUp(t)

		fx.leagueServiceClient.EXPECT().GetLeague(fx.ctx, &leaguepb.GetLeagueRequest{
			Id: 1,
		}).Return(&leaguepb.League{
			Id:   1,
			Name: "name",
		}, nil)

		got, err := fx.facade.GetLeague(fx.ctx, 1)
		assert.Equal(t, model.League{
			ID:   1,
			Name: "name",
		}, got)
		assert.NoError(t, err)
	})

	t.Run("ok with cache", func(t *testing.T) {
		fx := tearUp(t)
		fx.facade.leaguesCache[1] = model.League{
			ID:   1,
			Name: "name",
		}

		got, err := fx.facade.GetLeague(fx.ctx, 1)
		assert.Equal(t, model.League{
			ID:   1,
			Name: "name",
		}, got)
		assert.NoError(t, err)
	})
}

func Test_convertPBLeagueToModelLeague(t *testing.T) {
	type args struct {
		pbLeague *leaguepb.League
	}
	tests := []struct {
		name string
		args args
		want model.League
	}{
		{
			name: "test case 1",
			args: args{
				pbLeague: &leaguepb.League{
					Id:        1,
					Name:      "name",
					ShortName: "short_name",
					LogoLink:  "link",
					WebSite:   "site",
				},
			},
			want: model.League{
				ID:        1,
				Name:      "name",
				ShortName: "short_name",
				LogoLink:  "link",
				WebSite:   "site",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertPBLeagueToModelLeague(tt.args.pbLeague); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertPBLeagueToModelLeague() = %v, want %v", got, tt.want)
			}
		})
	}
}
