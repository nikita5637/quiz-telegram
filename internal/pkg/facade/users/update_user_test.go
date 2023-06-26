package users

import (
	"errors"
	"testing"

	usermanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/user_manager"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestFacade_UpdateUserEmail(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		fx := tearUp(t)

		fx.userManagerServiceClient.EXPECT().PatchUser(fx.ctx, &usermanagerpb.PatchUserRequest{
			User: &usermanagerpb.User{
				Id:    1,
				Email: "email",
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
				Id:    1,
				Email: "email",
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
				Id:    1,
				Phone: "phone",
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
				Id:    1,
				Phone: "phone",
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
