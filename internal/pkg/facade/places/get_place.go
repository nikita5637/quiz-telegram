package places

import (
	"context"
	"fmt"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetPlaceByID ...
func (f *Facade) GetPlaceByID(ctx context.Context, placeID int32) (model.Place, error) {
	if place, ok := f.placesCache[placeID]; ok {
		return place, nil
	}

	logger.DebugKV(ctx, "place not found in cache", "place ID", placeID)

	placeResp, err := f.registratorServiceClient.GetPlaceByID(ctx, &registrator.GetPlaceByIDRequest{
		Id: placeID,
	})
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.NotFound {
			return model.Place{}, model.ErrPlaceNotFound
		}

		return model.Place{}, fmt.Errorf("get place error: %w", err)
	}

	place := convertPBPlaceToModelPlace(placeResp.GetPlace())
	f.placesCache[placeID] = place

	return place, nil
}

func convertPBPlaceToModelPlace(pbPlace *registrator.Place) model.Place {
	return model.Place{
		ID:        pbPlace.GetId(),
		Address:   pbPlace.GetAddress(),
		Name:      pbPlace.GetName(),
		ShortName: pbPlace.GetShortName(),
		Longitude: pbPlace.GetLongitude(),
		Latitude:  pbPlace.GetLatitude(),
		MenuLink:  pbPlace.GetMenuLink(),
	}
}
