package repository

import "certlens/internal/domains"

type mockRepository struct {
	mockGetTLSSecrets func() ([]domains.SecretInfo, error)
	mockGetTLSSecret  func() (domains.SecretInfo, error)
}

func NewMockRepository(
	mockGetTLSSecrets func() ([]domains.SecretInfo, error),
	mockGetTLSSecret func() (domains.SecretInfo, error),
) SecretsRepository {
	return mockRepository{
		mockGetTLSSecrets: mockGetTLSSecrets,
		mockGetTLSSecret:  mockGetTLSSecret,
	}
}

func (m mockRepository) GetTLSSecrets(namespace string) ([]domains.SecretInfo, error) {
	return m.mockGetTLSSecrets()
}

func (m mockRepository) GetTLSSecret(namespace, name string) (domains.SecretInfo, error) {
	return m.mockGetTLSSecret()
}
