package service

import (
	"github.com/codechamp1/certlens/internal/domains/tls"
)

type mockSecretService struct {
	mockListTLSSecrets func(namespace string) ([]tls.Secret, error)
	mockLoadTLSSecret  func(namespace, name string) (tls.Secret, error)
}

func NewMockSecretService(
	mockListTLSSecrets func(namespace string) ([]tls.Secret, error),
	mockLoadTLSSecret func(namespace, name string) (tls.Secret, error)) Manager {
	return mockSecretService{
		mockLoadTLSSecret:  mockLoadTLSSecret,
		mockListTLSSecrets: mockListTLSSecrets,
	}
}

func (m mockSecretService) ListTLSSecrets(namespace string) ([]tls.Secret, error) {
	return m.mockListTLSSecrets(namespace)
}

func (m mockSecretService) LoadTLSSecret(namespace, name string) (tls.Secret, error) {
	return m.mockLoadTLSSecret(namespace, name)
}
