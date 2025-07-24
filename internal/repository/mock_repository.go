package repository

import "certlens/internal/domains"

type mockRepository struct {
	mockGetTLSSecrets func(namespace string) ([]domains.SecretInfo, error)
	mockGetTLSSecret  func(namespace string, name string) (domains.SecretInfo, error)
}

func NewMockRepository(
	mockGetTLSSecrets func(namespace string) ([]domains.SecretInfo, error),
	mockGetTLSSecret func(namespace, name string) (domains.SecretInfo, error),
) SecretsRepository {
	return mockRepository{
		mockGetTLSSecrets: mockGetTLSSecrets,
		mockGetTLSSecret:  mockGetTLSSecret,
	}
}

func (m mockRepository) GetTLSSecrets(namespace string) ([]domains.SecretInfo, error) {
	return m.mockGetTLSSecrets(namespace)
}

func (m mockRepository) GetTLSSecret(namespace, name string) (domains.SecretInfo, error) {
	return m.mockGetTLSSecret(namespace, name)
}
