//go:generate mockery --case underscore --name RegistratorServiceClient --with-expecter

package users

import (
	"context"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"google.golang.org/grpc"
)

// RegistratorServiceClient ...
type RegistratorServiceClient interface {
	CreateUser(ctx context.Context, in *registrator.CreateUserRequest, opts ...grpc.CallOption) (*registrator.CreateUserResponse, error)
	GetUserByID(ctx context.Context, in *registrator.GetUserByIDRequest, opts ...grpc.CallOption) (*registrator.GetUserByIDResponse, error)
	GetUserByTelegramID(ctx context.Context, in *registrator.GetUserByTelegramIDRequest, opts ...grpc.CallOption) (*registrator.GetUserByTelegramIDResponse, error)
	UpdateUserEmail(ctx context.Context, in *registrator.UpdateUserEmailRequest, opts ...grpc.CallOption) (*registrator.UpdateUserEmailResponse, error)
	UpdateUserName(ctx context.Context, in *registrator.UpdateUserNameRequest, opts ...grpc.CallOption) (*registrator.UpdateUserNameResponse, error)
	UpdateUserPhone(ctx context.Context, in *registrator.UpdateUserPhoneRequest, opts ...grpc.CallOption) (*registrator.UpdateUserPhoneResponse, error)
	UpdateUserState(ctx context.Context, in *registrator.UpdateUserStateRequest, opts ...grpc.CallOption) (*registrator.UpdateUserStateResponse, error)
}

// Facade ...
type Facade struct {
	registratorServiceClient RegistratorServiceClient
}

// Config ...
type Config struct {
	RegistratorServiceClient RegistratorServiceClient
}

// NewFacade ...
func NewFacade(cfg Config) *Facade {
	return &Facade{
		registratorServiceClient: cfg.RegistratorServiceClient,
	}
}
