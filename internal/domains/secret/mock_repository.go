package secret

type mockRepository struct {
	mockGetTLSSecrets func(namespace string) ([]TLS, error)
	mockGetTLSSecret  func(namespace string, name string) (TLS, error)
}

func NewMockRepository(
	mockGetTLSSecrets func(namespace string) ([]TLS, error),
	mockGetTLSSecret func(namespace, name string) (TLS, error),
) Repository {
	return mockRepository{
		mockGetTLSSecrets: mockGetTLSSecrets,
		mockGetTLSSecret:  mockGetTLSSecret,
	}
}

func (m mockRepository) GetTLSSecrets(namespace string) ([]TLS, error) {
	return m.mockGetTLSSecrets(namespace)
}

func (m mockRepository) GetTLSSecret(namespace, name string) (TLS, error) {
	return m.mockGetTLSSecret(namespace, name)
}
