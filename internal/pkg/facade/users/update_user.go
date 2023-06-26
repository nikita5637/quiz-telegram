package users

import (
	"context"

	usermanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/user_manager"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// UpdateUserEmail ...
func (f *Facade) UpdateUserEmail(ctx context.Context, userID int32, email string) error {
	_, err := f.userManagerServiceClient.PatchUser(ctx, &usermanagerpb.PatchUserRequest{
		User: &usermanagerpb.User{
			Id:    userID,
			Email: email,
			State: usermanagerpb.UserState_USER_STATE_REGISTERED,
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				"email",
				"state",
			},
		},
	})

	return err
}

// UpdateUserName ...
func (f *Facade) UpdateUserName(ctx context.Context, userID int32, name string) error {
	_, err := f.userManagerServiceClient.PatchUser(ctx, &usermanagerpb.PatchUserRequest{
		User: &usermanagerpb.User{
			Id:    userID,
			Name:  name,
			State: usermanagerpb.UserState_USER_STATE_REGISTERED,
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				"name",
				"state",
			},
		},
	})

	return err
}

// UpdateUserPhone ...
func (f *Facade) UpdateUserPhone(ctx context.Context, userID int32, phone string) error {
	_, err := f.userManagerServiceClient.PatchUser(ctx, &usermanagerpb.PatchUserRequest{
		User: &usermanagerpb.User{
			Id:    userID,
			Phone: phone,
			State: usermanagerpb.UserState_USER_STATE_REGISTERED,
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				"phone",
				"state",
			},
		},
	})

	return err
}

// UpdateUserState ...
func (f *Facade) UpdateUserState(ctx context.Context, userID, state int32) error {
	_, err := f.userManagerServiceClient.PatchUser(ctx, &usermanagerpb.PatchUserRequest{
		User: &usermanagerpb.User{
			Id:    userID,
			State: usermanagerpb.UserState(state),
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				"state",
			},
		},
	})

	return err
}
