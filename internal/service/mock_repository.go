package service

import "github.com/codechamp1/certlens/internal/domains/tls"

type mockRepository struct {
	mockGetTLSSecrets func(namespace string) ([]tls.Secret, error)
	mockGetTLSSecret  func(namespace string, name string) (tls.Secret, error)
}

func NewMockRepository(
	mockGetTLSSecrets func(namespace string) ([]tls.Secret, error),
	mockGetTLSSecret func(namespace, name string) (tls.Secret, error),
) tls.Repository {
	return mockRepository{
		mockGetTLSSecrets: mockGetTLSSecrets,
		mockGetTLSSecret:  mockGetTLSSecret,
	}
}

func (m mockRepository) GetTLSSecrets(namespace string) ([]tls.Secret, error) {
	return m.mockGetTLSSecrets(namespace)
}

func (m mockRepository) GetTLSSecret(namespace, name string) (tls.Secret, error) {
	return m.mockGetTLSSecret(namespace, name)
}
