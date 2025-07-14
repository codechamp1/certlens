package service

type mockSecretService struct {
	mockListTLSSecrets   func() ([]Secret, error)
	mockInspectTLSSecret func() (string, error)
}

func NewMockSecretService() SecretsService {
	return mockSecretService{}
}

func (m mockSecretService) InspectTLSSecret(namespace, name string) (string, error) {
	return m.mockInspectTLSSecret()
}

func (m mockSecretService) ListTLSSecrets(namespace string) ([]Secret, error) {
	return m.mockListTLSSecrets()
}
