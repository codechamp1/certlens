package cert

type mockService struct {
	mockParseTLSCert func(tlsCert []byte) ([]TLS, error)
}

func NewMockService(mockParseTLSCert func(tlsCert []byte) ([]TLS, error)) Service {
	return mockService{
		mockParseTLSCert,
	}
}

func (ms mockService) ParseTLSCert(tlsCert []byte) ([]TLS, error) {
	return ms.mockParseTLSCert(tlsCert)
}
