package model

import (
	certificatemanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/certificate_manager"
)

// CertificateType ...
type CertificateType int32

const (
	// CertificateTypeInvalid ...
	CertificateTypeInvalid CertificateType = iota
	// CertificateTypeFreePass ...
	CertificateTypeFreePass
	// CertificateTypeBarBillPayment ...
	CertificateTypeBarBillPayment
)

// String ...
func (t CertificateType) String() string {
	switch t {
	case CertificateType(certificatemanagerpb.CertificateType_CERTIFICATE_TYPE_FREE_PASS):
		return "Проходка"
	case CertificateType(certificatemanagerpb.CertificateType_CERTIFICATE_TYPE_BAR_BILL_PAYMENT):
		return "Счёт в баре"
	}

	return "unknown certificate type"
}

// Certificate ...
type Certificate struct {
	Type  CertificateType
	WonOn int32
	Info  string
}
