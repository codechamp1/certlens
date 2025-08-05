package service

import (
	"github.com/codechamp1/certlens/internal/domains/cert"
	"github.com/codechamp1/certlens/internal/domains/secret"
)

type mockSecretService struct {
	mockListTLSSecrets      func(namespace string) ([]secret.TLS, error)
	mockListTLSSecret       func(namespace, name string) (secret.TLS, error)
	mockInspectTLSSecret    func(namespace, name string) ([]cert.TLS, error)
	mockRawInspectTLSSecret func(namespace, name string) (string, string, error)
}

func NewMockSecretService(
	mockListTLSSecrets func(namespace string) ([]secret.TLS, error),
	mockListTLSSecret func(namespace, name string) (secret.TLS, error),
	mockInspectTLSSecret func(namespace, name string) ([]cert.TLS, error),
	mockRawInspectTLSSecret func(namespace, name string) (string, string, error)) Manager {
	return mockSecretService{
		mockInspectTLSSecret:    mockInspectTLSSecret,
		mockListTLSSecret:       mockListTLSSecret,
		mockListTLSSecrets:      mockListTLSSecrets,
		mockRawInspectTLSSecret: mockRawInspectTLSSecret,
	}
}

func (m mockSecretService) InspectTLSSecret(namespace, name string) ([]cert.TLS, error) {
	return m.mockInspectTLSSecret(namespace, name)
}

func (m mockSecretService) ListTLSSecrets(namespace string) ([]secret.TLS, error) {
	return m.mockListTLSSecrets(namespace)
}

func (m mockSecretService) ListTLSSecret(namespace, name string) (secret.TLS, error) {
	return m.mockListTLSSecret(namespace, name)
}

func (m mockSecretService) RawInspectTLSSecret(namespace, name string) (string, string, error) {
	return m.mockRawInspectTLSSecret(namespace, name)
}
