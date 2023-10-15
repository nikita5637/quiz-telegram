package users

import (
	"context"

	usermanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/user_manager"
	telegramutils "github.com/nikita5637/quiz-telegram/utils/telegram"
)

// CreateUser ...
func (f *Facade) CreateUser(ctx context.Context, name string, telegramID int64, state int32) (int32, error) {
	ctx = telegramutils.NewContextWithClientID(ctx, 0)
	resp, err := f.userManagerServiceClient.CreateUser(ctx, &usermanagerpb.CreateUserRequest{
		User: &usermanagerpb.User{
			Name:       name,
			TelegramId: telegramID,
			State:      usermanagerpb.UserState(state),
		},
	})
	if err != nil {
		return 0, err
	}

	return resp.GetId(), nil
}
