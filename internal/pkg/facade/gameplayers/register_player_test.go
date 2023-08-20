package gameplayers

import (
	"errors"
	"reflect"
	"testing"

	"github.com/mono83/maybe"
	gameplayerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/game_player"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/games"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestFacade_RegisterPlayer(t *testing.T) {
	t.Run("error while register player", func(t *testing.T) {
		fx := tearUp(t)

		fx.gamePlayerRegistratorServiceClient.EXPECT().RegisterPlayer(fx.ctx, &gameplayerpb.RegisterPlayerRequest{
			GamePlayer: &gameplayerpb.GamePlayer{
				GameId: 1,
				UserId: &wrapperspb.Int32Value{
					Value: 1,
				},
				RegisteredBy: 1,
				Degree:       gameplayerpb.Degree_DEGREE_LIKELY,
			},
		}).Return(nil, errors.New("some error"))

		err := fx.facade.RegisterPlayer(fx.ctx, model.GamePlayer{
			GameID:       1,
			UserID:       maybe.Just(int32(1)),
			RegisteredBy: 1,
			Degree:       model.DegreeLikely,
		})
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.gamePlayerRegistratorServiceClient.EXPECT().RegisterPlayer(fx.ctx, &gameplayerpb.RegisterPlayerRequest{
			GamePlayer: &gameplayerpb.GamePlayer{
				GameId: 1,
				UserId: &wrapperspb.Int32Value{
					Value: 1,
				},
				RegisteredBy: 1,
				Degree:       gameplayerpb.Degree_DEGREE_LIKELY,
			},
		}).Return(&emptypb.Empty{}, nil)

		err := fx.facade.RegisterPlayer(fx.ctx, model.GamePlayer{
			GameID:       1,
			UserID:       maybe.Just(int32(1)),
			RegisteredBy: 1,
			Degree:       model.DegreeLikely,
		})
		assert.NoError(t, err)
	})
}

func Test_convertModelGamePlayerToProtoGamePlayer(t *testing.T) {
	type args struct {
		gamePlayer model.GamePlayer
	}
	tests := []struct {
		name string
		args args
		want *gameplayerpb.GamePlayer
	}{
		{
			name: "tc1",
			args: args{
				gamePlayer: model.GamePlayer{
					ID:           1,
					GameID:       1,
					UserID:       maybe.Just(int32(1)),
					RegisteredBy: 1,
					Degree:       model.DegreeLikely,
				},
			},
			want: &gameplayerpb.GamePlayer{
				Id:     1,
				GameId: 1,
				UserId: &wrapperspb.Int32Value{
					Value: 1,
				},
				RegisteredBy: 1,
				Degree:       gameplayerpb.Degree_DEGREE_LIKELY,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertModelGamePlayerToProtoGamePlayer(tt.args.gamePlayer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertModelGamePlayerToProtoGamePlayer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleRegisterPlayerError(t *testing.T) {
	t.Run("error is nil", func(t *testing.T) {
		err := handleRegisterPlayerError(nil)
		assert.NoError(t, err)
	})

	t.Run("internal error", func(t *testing.T) {
		st := status.New(codes.Internal, "some error")
		err := handleRegisterPlayerError(st.Err())
		assert.Error(t, err)
	})

	t.Run("failed precondition error. reasonGameHasPassed", func(t *testing.T) {
		st := status.New(codes.FailedPrecondition, "some error")
		errorInfo := &errdetails.ErrorInfo{
			Reason: games.ReasonGameHasPassed,
		}
		st, err := st.WithDetails(errorInfo)
		assert.NoError(t, err)

		err = handleRegisterPlayerError(st.Err())
		assert.Error(t, err)
		assert.ErrorIs(t, err, games.ErrGameHasPassed)
	})

	t.Run("failed precondition error. reasonGameNotFound", func(t *testing.T) {
		st := status.New(codes.FailedPrecondition, "some error")
		errorInfo := &errdetails.ErrorInfo{
			Reason: games.ReasonGameNotFound,
		}
		st, err := st.WithDetails(errorInfo)
		assert.NoError(t, err)

		err = handleRegisterPlayerError(st.Err())
		assert.Error(t, err)
		assert.ErrorIs(t, err, games.ErrGameNotFound)
	})

	t.Run("failed precondition error. reasonNoFreeSlot", func(t *testing.T) {
		st := status.New(codes.FailedPrecondition, "some error")
		errorInfo := &errdetails.ErrorInfo{
			Reason: ReasonNoFreeSlot,
		}
		st, err := st.WithDetails(errorInfo)
		assert.NoError(t, err)

		err = handleRegisterPlayerError(st.Err())
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNoFreeSlot)
	})

	t.Run("already exists error", func(t *testing.T) {
		st := status.New(codes.AlreadyExists, "some error")
		err := handleRegisterPlayerError(st.Err())
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrGamePlayerAlreadyRegistered)
	})
}
