package data

import "github.com/codechamp1/certlens/internal/domains/tls"

type mockParser struct {
	mockParseTLSCerts func(tlsCert []byte) ([]tls.Cert, error)
}

func NewMockParser(mockParseTLSCerts func(tlsCert []byte) ([]tls.Cert, error)) CertParser {
	return mockParser{
		mockParseTLSCerts,
	}
}

func (ms mockParser) ParseTLSCerts(tlsCert []byte) ([]tls.Cert, error) {
	return ms.mockParseTLSCerts(tlsCert)
}
