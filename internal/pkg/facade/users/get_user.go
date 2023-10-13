package users

import (
	"context"
	"time"

	"github.com/mono83/maybe"
	usermanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/user_manager"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
)

// GetUser ...
func (f *Facade) GetUser(ctx context.Context, userID int32) (model.User, error) {
	pbUser, err := f.userManagerServiceClient.GetUser(ctx, &usermanagerpb.GetUserRequest{
		Id: userID,
	})
	if err != nil {
		return model.User{}, err
	}

	return convertProtoUserToModelUser(pbUser), nil
}

// GetUserByTelegramID ...
func (f *Facade) GetUserByTelegramID(ctx context.Context, telegramID int64) (model.User, error) {
	pbUser, err := f.userManagerServiceClient.GetUserByTelegramID(ctx, &usermanagerpb.GetUserByTelegramIDRequest{
		TelegramId: telegramID,
	})
	if err != nil {
		return model.User{}, err
	}

	return convertProtoUserToModelUser(pbUser), nil
}

func convertProtoUserToModelUser(pbUser *usermanagerpb.User) model.User {
	modelUser := model.User{
		ID:         pbUser.GetId(),
		Name:       pbUser.GetName(),
		TelegramID: pbUser.GetTelegramId(),
		Email:      maybe.Nothing[string](),
		Phone:      maybe.Nothing[string](),
		State:      int32(pbUser.GetState()),
		Birthdate:  maybe.Nothing[string](),
		Sex:        maybe.Nothing[model.Sex](),
	}

	if email := pbUser.GetEmail(); email != nil {
		modelUser.Email = maybe.Just(email.GetValue())
	}

	if phone := pbUser.GetPhone(); phone != nil {
		modelUser.Phone = maybe.Just(phone.GetValue())
	}

	if birthdate := pbUser.GetBirthdate(); birthdate != nil {
		birthdateTime, err := time.Parse("2006-01-02", birthdate.GetValue())
		if err == nil {
			modelUser.Birthdate = maybe.Just(birthdateTime.Format("02.01.2006"))
		}
	}

	if pbUser != nil && pbUser.Sex != nil {
		modelUser.Sex = maybe.Just(model.Sex(pbUser.GetSex()))
	}

	return modelUser
}
