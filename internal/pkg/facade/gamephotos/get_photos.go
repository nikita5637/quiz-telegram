package gamephotos

import (
	"context"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetPhotosByGameID ...
func (f *Facade) GetPhotosByGameID(ctx context.Context, gameID int32) ([]string, error) {
	resp, err := f.photographerServiceClient.GetPhotosByGameID(ctx, &registrator.GetPhotosByGameIDRequest{
		GameId: gameID,
	})
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.NotFound {
			return nil, model.ErrGameNotFound
		}

		return nil, err
	}

	return resp.GetUrls(), nil
}
