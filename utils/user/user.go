package user

import (
	"context"

	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
)

type userKeyType struct{}

var (
	userKey = userKeyType{}
)

// NewContextWithUser ...
func NewContextWithUser(ctx context.Context, user model.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// GetUserFromContext ...
func GetUserFromContext(ctx context.Context) model.User {
	val, ok := ctx.Value(userKey).(model.User)
	if !ok {
		return model.User{}
	}

	return val
}
