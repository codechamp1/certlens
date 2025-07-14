package repository

type mockRepository struct {
	mockGetTLSSecrets func() ([]SecretInfo, error)
	mockGetTLSSecret  func() (SecretInfo, error)
}

func NewMockRepository() SecretsRepository {
	return mockRepository{}
}

func (m mockRepository) GetTLSSecrets(namespace string) ([]SecretInfo, error) {
	return m.mockGetTLSSecrets()
}

func (m mockRepository) GetTLSSecret(namespace, name string) (SecretInfo, error) {
	return m.mockGetTLSSecret()
}
