package certificates

import (
	"context"

	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GetActiveCertificates ...
func (f *Facade) GetActiveCertificates(ctx context.Context) ([]model.Certificate, error) {
	resp, err := f.certificateManagerServiceClient.ListCertificates(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	modelCertificates := make([]model.Certificate, 0, len(resp.GetCertificates()))
	for _, pbCertificate := range resp.GetCertificates() {
		if pbCertificate.GetSpentOn() != nil {
			continue
		}

		modelCertificates = append(modelCertificates, model.Certificate{
			Type:  model.CertificateType(pbCertificate.GetType()),
			WonOn: pbCertificate.GetWonOn(),
			Info:  pbCertificate.GetInfo().GetValue(),
		})
	}

	return modelCertificates, nil
}
