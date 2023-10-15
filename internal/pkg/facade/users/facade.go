//go:generate mockery --case underscore --name UserManagerServiceClient --with-expecter

package users

import (
	"context"

	usermanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/user_manager"
	"google.golang.org/grpc"
)

// UserManagerServiceClient ...
type UserManagerServiceClient interface {
	CreateUser(ctx context.Context, in *usermanagerpb.CreateUserRequest, opts ...grpc.CallOption) (*usermanagerpb.User, error)
	GetUser(ctx context.Context, in *usermanagerpb.GetUserRequest, opts ...grpc.CallOption) (*usermanagerpb.User, error)
	GetUserByTelegramID(ctx context.Context, in *usermanagerpb.GetUserByTelegramIDRequest, opts ...grpc.CallOption) (*usermanagerpb.User, error)
	PatchUser(ctx context.Context, in *usermanagerpb.PatchUserRequest, opts ...grpc.CallOption) (*usermanagerpb.User, error)
}

// Facade ...
type Facade struct {
	userManagerServiceClient UserManagerServiceClient
}

// Config ...
type Config struct {
	UserManagerServiceClient UserManagerServiceClient
}

// NewFacade ...
func NewFacade(cfg Config) *Facade {
	return &Facade{
		userManagerServiceClient: cfg.UserManagerServiceClient,
	}
}
