//go:generate mockery --case underscore --name CertificateManagerServiceClient --with-expecter

package certificates

import (
	"context"

	certificatemanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/certificate_manager"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CertificateManagerServiceClient ...
type CertificateManagerServiceClient interface {
	ListCertificates(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*certificatemanagerpb.ListCertificatesResponse, error)
}

// Facade ...
type Facade struct {
	certificateManagerServiceClient CertificateManagerServiceClient
}

// Config ...
type Config struct {
	CertificateManagerServiceClient CertificateManagerServiceClient
}

// New ...
func New(cfg Config) *Facade {
	return &Facade{
		certificateManagerServiceClient: cfg.CertificateManagerServiceClient,
	}
}
