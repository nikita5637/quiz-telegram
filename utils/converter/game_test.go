package converter

import (
	"testing"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	time_utils "github.com/nikita5637/quiz-registrator-api/utils/time"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestConvertPBGameToModelGame(t *testing.T) {
	timeNow := time_utils.TimeNow()
	type args struct {
		pbGame *registrator.Game
	}
	tests := []struct {
		name string
		args args
		want model.Game
	}{
		{
			name: "test case 1",
			args: args{
				pbGame: &registrator.Game{
					Id:                  1,
					ExternalId:          2,
					LeagueId:            3,
					Type:                registrator.GameType_GAME_TYPE_CLASSIC,
					Number:              "1",
					Name:                "name",
					PlaceId:             4,
					Date:                timestamppb.New(timeNow),
					Price:               400,
					PaymentType:         "cash,card",
					MaxPlayers:          9,
					Payment:             registrator.Payment_PAYMENT_CASH,
					Registered:          true,
					My:                  true,
					NumberOfMyLegioners: 3,
					NumberOfLegioners:   4,
					NumberOfPlayers:     5,
					ResultPlace:         1,
				},
			},
			want: model.Game{
				ID:                  1,
				ExternalID:          2,
				Type:                1,
				Number:              "1",
				Name:                "name",
				Date:                model.DateTime(timestamppb.New(timeNow).AsTime()),
				Price:               400,
				PaymentType:         "cash,card",
				MaxPlayers:          9,
				Payment:             model.PaymentTypeCash,
				Registered:          true,
				My:                  true,
				NumberOfMyLegioners: 3,
				NumberOfLegioners:   4,
				NumberOfPlayers:     5,
				ResultPlace:         1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertPBGameToModelGame(tt.args.pbGame)
			assert.Equal(t, tt.want, got)
		})
	}
}
