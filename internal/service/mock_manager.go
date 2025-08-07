package service

import (
	"github.com/codechamp1/certlens/internal/domains/cert"
	"github.com/codechamp1/certlens/internal/domains/secret"
)

type mockSecretService struct {
	mockListTLSSecrets   func(namespace string) ([]secret.TLS, error)
	mockListTLSSecret    func(namespace, name string) (secret.TLS, error)
	mockInspectTLSSecret func(tlsSecret secret.TLS) ([]cert.TLS, error)
}

func NewMockSecretService(
	mockListTLSSecrets func(namespace string) ([]secret.TLS, error),
	mockListTLSSecret func(namespace, name string) (secret.TLS, error),
	mockInspectTLSSecret func(tlsSecret secret.TLS) ([]cert.TLS, error)) Manager {
	return mockSecretService{
		mockInspectTLSSecret: mockInspectTLSSecret,
		mockListTLSSecret:    mockListTLSSecret,
		mockListTLSSecrets:   mockListTLSSecrets,
	}
}

func (m mockSecretService) InspectTLSSecret(tlsSecret secret.TLS) ([]cert.TLS, error) {
	return m.mockInspectTLSSecret(tlsSecret)
}

func (m mockSecretService) ListTLSSecrets(namespace string) ([]secret.TLS, error) {
	return m.mockListTLSSecrets(namespace)
}

func (m mockSecretService) ListTLSSecret(namespace, name string) (secret.TLS, error) {
	return m.mockListTLSSecret(namespace, name)
}
