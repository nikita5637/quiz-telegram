package users

import (
	"context"
	"time"

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
	birthdate := ""
	if pbUser.GetBirthdate() != nil {
		birthdateTime, err := time.Parse("2006-01-02", pbUser.GetBirthdate().GetValue())
		if err == nil {
			birthdate = birthdateTime.Format("02.01.2006")
		}
	}

	return model.User{
		ID:        pbUser.GetId(),
		Email:     pbUser.GetEmail().GetValue(),
		Name:      pbUser.GetName(),
		Phone:     pbUser.GetPhone().GetValue(),
		State:     int32(pbUser.GetState()),
		Birthdate: birthdate,
		Sex:       model.Sex(pbUser.GetSex()),
	}
}
