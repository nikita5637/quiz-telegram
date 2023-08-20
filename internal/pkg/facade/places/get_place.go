package places

import (
	"context"
	"fmt"

	placepb "github.com/nikita5637/quiz-registrator-api/pkg/pb/place"
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

	pbPlace, err := f.placeServiceClient.GetPlace(ctx, &placepb.GetPlaceRequest{
		Id: placeID,
	})
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.NotFound {
			return model.Place{}, ErrPlaceNotFound
		}

		return model.Place{}, fmt.Errorf("get place error: %w", err)
	}

	place := convertPBPlaceToModelPlace(pbPlace)
	f.placesCache[placeID] = place

	return place, nil
}

func convertPBPlaceToModelPlace(pbPlace *placepb.Place) model.Place {
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
