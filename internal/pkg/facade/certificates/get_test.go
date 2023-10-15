package certificates

import (
	"errors"
	"testing"

	certificatemanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/certificate_manager"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestFacade_GetActiveCertificates(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		fx := tearUp(t)

		fx.certificateManagerServiceClient.EXPECT().ListCertificates(fx.ctx, &emptypb.Empty{}).Return(nil, errors.New("some error"))

		got, err := fx.facade.GetActiveCertificates(fx.ctx)
		assert.Nil(t, got)
		assert.Error(t, err)
	})

	t.Run("empty list", func(t *testing.T) {
		fx := tearUp(t)

		fx.certificateManagerServiceClient.EXPECT().ListCertificates(fx.ctx, &emptypb.Empty{}).Return(&certificatemanagerpb.ListCertificatesResponse{
			Certificates: []*certificatemanagerpb.Certificate{},
		}, nil)

		got, err := fx.facade.GetActiveCertificates(fx.ctx)
		assert.Equal(t, []model.Certificate{}, got)
		assert.NoError(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		fx := tearUp(t)

		fx.certificateManagerServiceClient.EXPECT().ListCertificates(fx.ctx, &emptypb.Empty{}).Return(&certificatemanagerpb.ListCertificatesResponse{
			Certificates: []*certificatemanagerpb.Certificate{
				{
					Id:   1,
					Type: certificatemanagerpb.CertificateType_CERTIFICATE_TYPE_BAR_BILL_PAYMENT,
					Info: &wrapperspb.StringValue{
						Value: "info1",
					},
				},
				{
					Id:   2,
					Type: certificatemanagerpb.CertificateType_CERTIFICATE_TYPE_FREE_PASS,
					SpentOn: &wrapperspb.Int32Value{
						Value: 1,
					},
					Info: &wrapperspb.StringValue{
						Value: "info2",
					},
				},
				{
					Id:   3,
					Type: certificatemanagerpb.CertificateType_CERTIFICATE_TYPE_FREE_PASS,
					Info: &wrapperspb.StringValue{
						Value: "info3",
					},
				},
			},
		}, nil)

		got, err := fx.facade.GetActiveCertificates(fx.ctx)
		assert.ElementsMatch(t, []model.Certificate{
			{
				Type: int32(certificatemanagerpb.CertificateType_CERTIFICATE_TYPE_BAR_BILL_PAYMENT),
				Info: "info1",
			},
			{
				Type: int32(certificatemanagerpb.CertificateType_CERTIFICATE_TYPE_FREE_PASS),
				Info: "info3",
			},
		}, got)
		assert.NoError(t, err)
	})
}
