package service

import "certlens/internal/domains"

type mockSecretService struct {
	mockListTLSSecrets   func() ([]domains.K8SResourceID, error)
	mockInspectTLSSecret func() (*CertificateInfo, error)
}

func NewMockSecretService(mockListTLSSecrets func() ([]domains.K8SResourceID, error), mockInspectTLSSecret func() (*CertificateInfo, error)) SecretsService {
	return mockSecretService{
		mockInspectTLSSecret: mockInspectTLSSecret,
		mockListTLSSecrets:   mockListTLSSecrets,
	}
}

func (m mockSecretService) InspectTLSSecret(namespace, name string) (*CertificateInfo, error) {
	return m.mockInspectTLSSecret()
}

func (m mockSecretService) ListTLSSecrets(namespace string) ([]domains.K8SResourceID, error) {
	return m.mockListTLSSecrets()
}

func (m mockSecretService) ListTLSSecret(namespace, name string) (domains.K8SResourceID, error) {
	return domains.K8SResourceID{}, nil
}
