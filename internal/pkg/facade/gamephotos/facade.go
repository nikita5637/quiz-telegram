//go:generate mockery --case underscore --name PhotographerServiceClient --with-expecter

package gamephotos

import (
	"context"

	photomanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/photo_manager"
	"google.golang.org/grpc"
)

// PhotographerServiceClient ...
type PhotographerServiceClient interface {
	GetPhotosByGameID(ctx context.Context, in *photomanagerpb.GetPhotosByGameIDRequest, opts ...grpc.CallOption) (*photomanagerpb.GetPhotosByGameIDResponse, error)
}

// Facade ...
type Facade struct {
	photographerServiceClient PhotographerServiceClient
}

// Config ...
type Config struct {
	PhotographerServiceClient PhotographerServiceClient
}

// New ...
func New(cfg Config) *Facade {
	return &Facade{
		photographerServiceClient: cfg.PhotographerServiceClient,
	}
}
