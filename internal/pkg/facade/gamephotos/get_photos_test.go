package gamephotos

import (
	"errors"
	"testing"

	photomanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/photo_manager"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/games"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestFacade_GetPhotosByGameID(t *testing.T) {
	t.Run("error game not found while get photos by game ID", func(t *testing.T) {
		fx := tearUp(t)

		fx.photographerServiceClient.EXPECT().GetPhotosByGameID(fx.ctx, &photomanagerpb.GetPhotosByGameIDRequest{
			GameId: 1,
		}).Return(nil, status.New(codes.NotFound, "").Err())

		got, err := fx.facade.GetPhotosByGameID(fx.ctx, int32(1))
		assert.Nil(t, got)
		assert.Error(t, err)
		assert.ErrorIs(t, err, games.ErrGameNotFound)
	})

	t.Run("error while get photos by game ID", func(t *testing.T) {
		fx := tearUp(t)

		fx.photographerServiceClient.EXPECT().GetPhotosByGameID(fx.ctx, &photomanagerpb.GetPhotosByGameIDRequest{
			GameId: 1,
		}).Return(nil, errors.New("some error"))

		got, err := fx.facade.GetPhotosByGameID(fx.ctx, int32(1))
		assert.Nil(t, got)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.photographerServiceClient.EXPECT().GetPhotosByGameID(fx.ctx, &photomanagerpb.GetPhotosByGameIDRequest{
			GameId: 1,
		}).Return(&photomanagerpb.GetPhotosByGameIDResponse{
			Urls: []string{
				"url1",
				"url2",
				"url3",
			},
		}, nil)

		got, err := fx.facade.GetPhotosByGameID(fx.ctx, int32(1))
		assert.Equal(t, []string{
			"url1",
			"url2",
			"url3",
		}, got)
		assert.NoError(t, err)
	})
}
