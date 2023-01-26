package users

import (
	"context"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
)

// UpdateUserEmail ...
func (f *Facade) UpdateUserEmail(ctx context.Context, userID int32, email string) error {
	_, err := f.registratorServiceClient.UpdateUserEmail(ctx, &registrator.UpdateUserEmailRequest{
		UserId: userID,
		Email:  email,
	})

	return err
}

//  UpdateUserName ...
func (f *Facade) UpdateUserName(ctx context.Context, userID int32, name string) error {
	_, err := f.registratorServiceClient.UpdateUserName(ctx, &registrator.UpdateUserNameRequest{
		UserId: userID,
		Name:   name,
	})

	return err
}

// UpdateUserPhone ...
func (f *Facade) UpdateUserPhone(ctx context.Context, userID int32, phone string) error {
	_, err := f.registratorServiceClient.UpdateUserPhone(ctx, &registrator.UpdateUserPhoneRequest{
		UserId: userID,
		Phone:  phone,
	})

	return err
}

// UpdateUserState ...
func (f *Facade) UpdateUserState(ctx context.Context, userID, state int32) error {
	_, err := f.registratorServiceClient.UpdateUserState(ctx, &registrator.UpdateUserStateRequest{
		UserId: userID,
		State:  registrator.UserState(state),
	})

	return err
}
