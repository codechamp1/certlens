package repository

type mockRepository struct {
	mockGetTLSSecrets func() ([]SecretInfo, error)
	mockGetTLSSecret  func() (SecretInfo, error)
}

func NewMockRepository(
	mockGetTLSSecrets func() ([]SecretInfo, error),
	mockGetTLSSecret func() (SecretInfo, error),
) SecretsRepository {
	return mockRepository{
		mockGetTLSSecrets: mockGetTLSSecrets,
		mockGetTLSSecret:  mockGetTLSSecret,
	}
}

func (m mockRepository) GetTLSSecrets(namespace string) ([]SecretInfo, error) {
	return m.mockGetTLSSecrets()
}

func (m mockRepository) GetTLSSecret(namespace, name string) (SecretInfo, error) {
	return m.mockGetTLSSecret()
}
