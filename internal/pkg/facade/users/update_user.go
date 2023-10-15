package users

import (
	"context"
	"time"

	usermanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/user_manager"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// UpdateUserBirthdate ...
func (f *Facade) UpdateUserBirthdate(ctx context.Context, userID int32, birthdate string) error {
	birthdateTime, err := time.Parse("02.01.2006", birthdate)
	if err != nil {
		return err
	}

	_, err = f.userManagerServiceClient.PatchUser(ctx, &usermanagerpb.PatchUserRequest{
		User: &usermanagerpb.User{
			Id:    userID,
			State: usermanagerpb.UserState_USER_STATE_REGISTERED,
			Birthdate: &wrapperspb.StringValue{
				Value: birthdateTime.Format("2006-01-02"),
			},
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				"state",
				"birthdate",
			},
		},
	})

	return err
}

// UpdateUserEmail ...
func (f *Facade) UpdateUserEmail(ctx context.Context, userID int32, email string) error {
	_, err := f.userManagerServiceClient.PatchUser(ctx, &usermanagerpb.PatchUserRequest{
		User: &usermanagerpb.User{
			Id: userID,
			Email: &wrapperspb.StringValue{
				Value: email,
			},
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
			Id: userID,
			Phone: &wrapperspb.StringValue{
				Value: phone,
			},
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

// UpdateUserSex ...
func (f *Facade) UpdateUserSex(ctx context.Context, userID int32, sex model.Sex) error {
	_, err := f.userManagerServiceClient.PatchUser(ctx, &usermanagerpb.PatchUserRequest{
		User: &usermanagerpb.User{
			Id:    userID,
			State: usermanagerpb.UserState_USER_STATE_REGISTERED,
			Sex:   (*usermanagerpb.Sex)(&sex),
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				"state",
				"sex",
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
