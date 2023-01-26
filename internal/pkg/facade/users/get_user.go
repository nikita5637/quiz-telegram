package users

import (
	"context"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
)

// GetUserByID ...
func (f *Facade) GetUserByID(ctx context.Context, userID int32) (model.User, error) {
	resp, err := f.registratorServiceClient.GetUserByID(ctx, &registrator.GetUserByIDRequest{
		Id: userID,
	})
	if err != nil {
		return model.User{}, err
	}

	return convertPBUserToModelUser(resp.GetUser()), nil
}

// GetUserByTelegramID ...
func (f *Facade) GetUserByTelegramID(ctx context.Context, telegramID int64) (model.User, error) {
	resp, err := f.registratorServiceClient.GetUserByTelegramID(ctx, &registrator.GetUserByTelegramIDRequest{
		TelegramId: telegramID,
	})
	if err != nil {
		return model.User{}, err
	}

	return convertPBUserToModelUser(resp.GetUser()), nil
}

func convertPBUserToModelUser(pbUser *registrator.User) model.User {
	return model.User{
		ID:    pbUser.GetId(),
		Email: pbUser.GetEmail(),
		Name:  pbUser.GetName(),
		Phone: pbUser.GetPhone(),
		State: int32(pbUser.GetState()),
	}
}
