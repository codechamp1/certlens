package data

import "github.com/codechamp1/certlens/internal/domains/tls"

type mockParser struct {
	mockParseTLSCert func(tlsCert []byte) ([]tls.Cert, error)
}

func NewMockParser(mockParseTLSCert func(tlsCert []byte) ([]tls.Cert, error)) CertParser {
	return mockParser{
		mockParseTLSCert,
	}
}

func (ms mockParser) ParseTLSCert(tlsCert []byte) ([]tls.Cert, error) {
	return ms.mockParseTLSCert(tlsCert)
}
