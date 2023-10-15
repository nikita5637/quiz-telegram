//go:generate mockery --case underscore --name ICSFileManagerAPIServiceClient --with-expecter

package icsfiles

import (
	"context"

	icsfilemanagerpb "github.com/nikita5637/quiz-ics-manager-api/pkg/pb/ics_file_manager"
	"google.golang.org/grpc"
)

// ICSFileManagerAPIServiceClient ...
type ICSFileManagerAPIServiceClient interface {
	GetICSFileByGameID(ctx context.Context, in *icsfilemanagerpb.GetICSFileByGameIDRequest, opts ...grpc.CallOption) (*icsfilemanagerpb.ICSFile, error)
}

// Facade ...
type Facade struct {
	icsFileManagerAPIClient ICSFileManagerAPIServiceClient
}

// Config ...
type Config struct {
	ICSFileManagerAPIServiceClient ICSFileManagerAPIServiceClient
}

// NewFacade ...
func NewFacade(cfg Config) *Facade {
	return &Facade{
		icsFileManagerAPIClient: cfg.ICSFileManagerAPIServiceClient,
	}
}
