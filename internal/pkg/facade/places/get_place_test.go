package places

import (
	"errors"
	"reflect"
	"testing"

	placepb "github.com/nikita5637/quiz-registrator-api/pkg/pb/place"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestFacade_GetPlaceByID(t *testing.T) {
	t.Run("error place not found while get place", func(t *testing.T) {
		fx := tearUp(t)

		fx.placeServiceClient.EXPECT().GetPlace(fx.ctx, &placepb.GetPlaceRequest{
			Id: 1,
		}).Return(nil, status.New(codes.NotFound, "").Err())

		got, err := fx.facade.GetPlaceByID(fx.ctx, 1)
		assert.Equal(t, model.Place{}, got)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrPlaceNotFound)
	})

	t.Run("error while get place", func(t *testing.T) {
		fx := tearUp(t)

		fx.placeServiceClient.EXPECT().GetPlace(fx.ctx, &placepb.GetPlaceRequest{
			Id: 1,
		}).Return(nil, errors.New("some error"))

		got, err := fx.facade.GetPlaceByID(fx.ctx, 1)
		assert.Equal(t, model.Place{}, got)
		assert.Error(t, err)
	})

	t.Run("ok without cache", func(t *testing.T) {
		fx := tearUp(t)

		fx.placeServiceClient.EXPECT().GetPlace(fx.ctx, &placepb.GetPlaceRequest{
			Id: 1,
		}).Return(&placepb.Place{
			Id:   1,
			Name: "name",
		}, nil)

		got, err := fx.facade.GetPlaceByID(fx.ctx, 1)
		assert.Equal(t, model.Place{
			ID:   1,
			Name: "name",
		}, got)
		assert.NoError(t, err)
	})

	t.Run("ok with cache", func(t *testing.T) {
		fx := tearUp(t)
		fx.facade.placesCache[1] = model.Place{
			ID:   1,
			Name: "name",
		}

		got, err := fx.facade.GetPlaceByID(fx.ctx, 1)
		assert.Equal(t, model.Place{
			ID:   1,
			Name: "name",
		}, got)
		assert.NoError(t, err)
	})
}

func Test_convertPBPlaceToModelPlace(t *testing.T) {
	type args struct {
		pbPlace *placepb.Place
	}
	tests := []struct {
		name string
		args args
		want model.Place
	}{
		{
			name: "test case 1",
			args: args{
				pbPlace: &placepb.Place{
					Id:        1,
					Address:   "address",
					Name:      "name",
					ShortName: "short_name",
					Latitude:  1.1,
					Longitude: 2.2,
					MenuLink:  "menu",
				},
			},
			want: model.Place{
				ID:        1,
				Address:   "address",
				Name:      "name",
				ShortName: "short_name",
				Latitude:  1.1,
				Longitude: 2.2,
				MenuLink:  "menu",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertPBPlaceToModelPlace(tt.args.pbPlace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertPBPlaceToModelPlace() = %v, want %v", got, tt.want)
			}
		})
	}
}
