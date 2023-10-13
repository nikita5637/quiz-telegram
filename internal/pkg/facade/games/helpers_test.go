package games

import (
	"testing"

	"github.com/mono83/maybe"
	gamepb "github.com/nikita5637/quiz-registrator-api/pkg/pb/game"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	timeutils "github.com/nikita5637/quiz-telegram/utils/time"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func Test_convertProtoGameToModelGame(t *testing.T) {
	timeNow := timeutils.TimeNow()
	paymentCash := gamepb.Payment_PAYMENT_CASH
	type args struct {
		pbGame *gamepb.Game
	}
	tests := []struct {
		name string
		args args
		want model.Game
	}{
		{
			name: "test case 1",
			args: args{
				pbGame: &gamepb.Game{
					Id:          1,
					ExternalId:  wrapperspb.Int32(2),
					LeagueId:    3,
					Type:        gamepb.GameType_GAME_TYPE_CLASSIC,
					Number:      "1",
					Name:        wrapperspb.String("name"),
					PlaceId:     4,
					Date:        timestamppb.New(timeNow),
					Price:       400,
					PaymentType: wrapperspb.String("cash,card"),
					MaxPlayers:  9,
					Payment:     &paymentCash,
					Registered:  true,
					IsInMaster:  true,
					HasPassed:   true,
				},
			},
			want: model.Game{
				ID:          1,
				ExternalID:  maybe.Just(int32(2)),
				LeagueID:    3,
				Type:        1,
				Number:      "1",
				Name:        maybe.Just("name"),
				PlaceID:     4,
				DateTime:    model.DateTime(timestamppb.New(timeNow).AsTime()),
				Price:       400,
				PaymentType: maybe.Just("cash,card"),
				MaxPlayers:  9,
				Payment:     maybe.Just(int32(1)),
				Registered:  true,
				IsInMaster:  true,
				HasPassed:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertProtoGameToModelGame(tt.args.pbGame)
			assert.Equal(t, tt.want, got)
		})
	}
}
