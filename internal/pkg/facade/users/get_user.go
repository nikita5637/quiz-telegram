package users

import (
	"context"

	usermanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/user_manager"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
)

// GetUserByID ...
func (f *Facade) GetUserByID(ctx context.Context, userID int32) (model.User, error) {
	pbUser, err := f.userManagerServiceClient.GetUser(ctx, &usermanagerpb.GetUserRequest{
		Id: userID,
	})
	if err != nil {
		return model.User{}, err
	}

	return convertPBUserToModelUser(pbUser), nil
}

// GetUserByTelegramID ...
func (f *Facade) GetUserByTelegramID(ctx context.Context, telegramID int64) (model.User, error) {
	pbUser, err := f.userManagerServiceClient.GetUserByTelegramID(ctx, &usermanagerpb.GetUserByTelegramIDRequest{
		TelegramId: telegramID,
	})
	if err != nil {
		return model.User{}, err
	}

	return convertPBUserToModelUser(pbUser), nil
}

func convertPBUserToModelUser(pbUser *usermanagerpb.User) model.User {
	return model.User{
		ID:    pbUser.GetId(),
		Email: pbUser.GetEmail(),
		Name:  pbUser.GetName(),
		Phone: pbUser.GetPhone(),
		State: int32(pbUser.GetState()),
	}
}
