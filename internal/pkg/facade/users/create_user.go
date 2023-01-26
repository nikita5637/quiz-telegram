package users

import (
	"context"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
)

// CreateUser ...
func (f *Facade) CreateUser(ctx context.Context, name string, telegramID int64, state int32) (int32, error) {
	resp, err := f.registratorServiceClient.CreateUser(ctx, &registrator.CreateUserRequest{
		Name:       name,
		TelegramId: telegramID,
		State:      registrator.UserState(state),
	})
	if err != nil {
		return 0, err
	}

	return resp.GetId(), nil
}
