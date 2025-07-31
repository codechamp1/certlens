package service

import "github.com/codechamp1/certlens/internal/domains"

type mockSecretService struct {
	mockListTLSSecrets      func(namespace string) ([]domains.K8SResourceID, error)
	mockListTLSSecret       func(namespace, name string) (domains.K8SResourceID, error)
	mockInspectTLSSecret    func(namespace, name string) ([]CertificateInfo, error)
	mockRawInspectTLSSecret func(namespace, name string) (string, string, error)
}

func NewMockSecretService(
	mockListTLSSecrets func(namespace string) ([]domains.K8SResourceID, error),
	mockListTLSSecret func(namespace, name string) (domains.K8SResourceID, error),
	mockInspectTLSSecret func(namespace, name string) ([]CertificateInfo, error),
	mockRawInspectTLSSecret func(namespace, name string) (string, string, error)) SecretsService {
	return mockSecretService{
		mockInspectTLSSecret:    mockInspectTLSSecret,
		mockListTLSSecret:       mockListTLSSecret,
		mockListTLSSecrets:      mockListTLSSecrets,
		mockRawInspectTLSSecret: mockRawInspectTLSSecret,
	}
}

func (m mockSecretService) InspectTLSSecret(namespace, name string) ([]CertificateInfo, error) {
	return m.mockInspectTLSSecret(namespace, name)
}

func (m mockSecretService) ListTLSSecrets(namespace string) ([]domains.K8SResourceID, error) {
	return m.mockListTLSSecrets(namespace)
}

func (m mockSecretService) ListTLSSecret(namespace, name string) (domains.K8SResourceID, error) {
	return m.mockListTLSSecret(namespace, name)
}

func (m mockSecretService) RawInspectTLSSecret(namespace, name string) (string, string, error) {
	return m.mockRawInspectTLSSecret(namespace, name)
}
