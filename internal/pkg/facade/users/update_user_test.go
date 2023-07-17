package users

import (
	"errors"
	"testing"

	usermanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/user_manager"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestFacade_UpdateUserBirthdate(t *testing.T) {
	t.Run("parse error", func(t *testing.T) {
		fx := tearUp(t)

		err := fx.facade.UpdateUserBirthdate(fx.ctx, 1, "invalid value")
		assert.Error(t, err)
	})

	t.Run("error", func(t *testing.T) {
		fx := tearUp(t)

		fx.userManagerServiceClient.EXPECT().PatchUser(fx.ctx, &usermanagerpb.PatchUserRequest{
			User: &usermanagerpb.User{
				Id: 1,
				Birthdate: &wrapperspb.StringValue{
					Value: "1990-01-30",
				},
				State: usermanagerpb.UserState_USER_STATE_REGISTERED,
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{
					"state",
					"birthdate",
				},
			},
		}).Return(nil, errors.New("some error"))

		err := fx.facade.UpdateUserBirthdate(fx.ctx, 1, "30.01.1990")
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.userManagerServiceClient.EXPECT().PatchUser(fx.ctx, &usermanagerpb.PatchUserRequest{
			User: &usermanagerpb.User{
				Id: 1,
				Birthdate: &wrapperspb.StringValue{
					Value: "1990-01-30",
				},
				State: usermanagerpb.UserState_USER_STATE_REGISTERED,
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{
					"state",
					"birthdate",
				},
			},
		}).Return(nil, nil)

		err := fx.facade.UpdateUserBirthdate(fx.ctx, 1, "30.01.1990")
		assert.NoError(t, err)
	})
}
func TestFacade_UpdateUserEmail(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		fx := tearUp(t)

		fx.userManagerServiceClient.EXPECT().PatchUser(fx.ctx, &usermanagerpb.PatchUserRequest{
			User: &usermanagerpb.User{
				Id: 1,
				Email: &wrapperspb.StringValue{
					Value: "email",
				},
				State: usermanagerpb.UserState_USER_STATE_REGISTERED,
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{
					"email",
					"state",
				},
			},
		}).Return(nil, errors.New("some error"))

		err := fx.facade.UpdateUserEmail(fx.ctx, 1, "email")
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.userManagerServiceClient.EXPECT().PatchUser(fx.ctx, &usermanagerpb.PatchUserRequest{
			User: &usermanagerpb.User{
				Id: 1,
				Email: &wrapperspb.StringValue{
					Value: "email",
				},
				State: usermanagerpb.UserState_USER_STATE_REGISTERED,
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{
					"email",
					"state",
				},
			},
		}).Return(nil, nil)

		err := fx.facade.UpdateUserEmail(fx.ctx, 1, "email")
		assert.NoError(t, err)
	})
}

func TestFacade_UpdateUserName(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		fx := tearUp(t)

		fx.userManagerServiceClient.EXPECT().PatchUser(fx.ctx, &usermanagerpb.PatchUserRequest{
			User: &usermanagerpb.User{
				Id:    1,
				Name:  "name",
				State: usermanagerpb.UserState_USER_STATE_REGISTERED,
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{
					"name",
					"state",
				},
			},
		}).Return(nil, errors.New("some error"))

		err := fx.facade.UpdateUserName(fx.ctx, 1, "name")
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.userManagerServiceClient.EXPECT().PatchUser(fx.ctx, &usermanagerpb.PatchUserRequest{
			User: &usermanagerpb.User{
				Id:    1,
				Name:  "name",
				State: usermanagerpb.UserState_USER_STATE_REGISTERED,
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{
					"name",
					"state",
				},
			},
		}).Return(nil, nil)

		err := fx.facade.UpdateUserName(fx.ctx, 1, "name")
		assert.NoError(t, err)
	})
}

func TestFacade_UpdateUserPhone(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		fx := tearUp(t)

		fx.userManagerServiceClient.EXPECT().PatchUser(fx.ctx, &usermanagerpb.PatchUserRequest{
			User: &usermanagerpb.User{
				Id: 1,
				Phone: &wrapperspb.StringValue{
					Value: "phone",
				},
				State: usermanagerpb.UserState_USER_STATE_REGISTERED,
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{
					"phone",
					"state",
				},
			},
		}).Return(nil, errors.New("some error"))

		err := fx.facade.UpdateUserPhone(fx.ctx, 1, "phone")
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.userManagerServiceClient.EXPECT().PatchUser(fx.ctx, &usermanagerpb.PatchUserRequest{
			User: &usermanagerpb.User{
				Id: 1,
				Phone: &wrapperspb.StringValue{
					Value: "phone",
				},
				State: usermanagerpb.UserState_USER_STATE_REGISTERED,
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{
					"phone",
					"state",
				},
			},
		}).Return(nil, nil)

		err := fx.facade.UpdateUserPhone(fx.ctx, 1, "phone")
		assert.NoError(t, err)
	})
}

func TestFacade_UpdateUserSex(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		fx := tearUp(t)

		pbSex := usermanagerpb.Sex_SEX_MALE

		fx.userManagerServiceClient.EXPECT().PatchUser(fx.ctx, &usermanagerpb.PatchUserRequest{
			User: &usermanagerpb.User{
				Id:    1,
				State: usermanagerpb.UserState_USER_STATE_REGISTERED,
				Sex:   &pbSex,
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{
					"state",
					"sex",
				},
			},
		}).Return(nil, errors.New("some error"))

		err := fx.facade.UpdateUserSex(fx.ctx, 1, model.Sex(1))
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		pbSex := usermanagerpb.Sex_SEX_MALE

		fx.userManagerServiceClient.EXPECT().PatchUser(fx.ctx, &usermanagerpb.PatchUserRequest{
			User: &usermanagerpb.User{
				Id:    1,
				State: usermanagerpb.UserState_USER_STATE_REGISTERED,
				Sex:   &pbSex,
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{
					"state",
					"sex",
				},
			},
		}).Return(nil, nil)

		err := fx.facade.UpdateUserSex(fx.ctx, 1, model.Sex(1))
		assert.NoError(t, err)
	})
}

func TestFacade_UpdateUserState(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		fx := tearUp(t)

		fx.userManagerServiceClient.EXPECT().PatchUser(fx.ctx, &usermanagerpb.PatchUserRequest{
			User: &usermanagerpb.User{
				Id:    1,
				State: usermanagerpb.UserState_USER_STATE_CHANGING_NAME,
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{
					"state",
				},
			},
		}).Return(nil, errors.New("some error"))

		err := fx.facade.UpdateUserState(fx.ctx, 1, int32(usermanagerpb.UserState_USER_STATE_CHANGING_NAME))
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.userManagerServiceClient.EXPECT().PatchUser(fx.ctx, &usermanagerpb.PatchUserRequest{
			User: &usermanagerpb.User{
				Id:    1,
				State: usermanagerpb.UserState_USER_STATE_CHANGING_NAME,
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{
					"state",
				},
			},
		}).Return(nil, nil)

		err := fx.facade.UpdateUserState(fx.ctx, 1, int32(usermanagerpb.UserState_USER_STATE_CHANGING_NAME))
		assert.NoError(t, err)
	})
}
