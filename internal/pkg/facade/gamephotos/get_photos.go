package gamephotos

import (
	"context"

	photomanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/photo_manager"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/games"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetPhotosByGameID ...
func (f *Facade) GetPhotosByGameID(ctx context.Context, gameID int32) ([]string, error) {
	resp, err := f.photographerServiceClient.GetPhotosByGameID(ctx, &photomanagerpb.GetPhotosByGameIDRequest{
		GameId: gameID,
	})
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.NotFound {
			return nil, games.ErrGameNotFound
		}

		return nil, err
	}

	return resp.GetUrls(), nil
}
