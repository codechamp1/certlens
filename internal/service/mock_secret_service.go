package service

type mockSecretService struct {
	mockListTLSSecrets   func() ([]Secret, error)
	mockInspectTLSSecret func() (string, error)
}

func NewMockSecretService(mockListTLSSecrets func() ([]Secret, error), mockInspectTLSSecret func() (string, error)) SecretsService {
	return mockSecretService{
		mockInspectTLSSecret: mockInspectTLSSecret,
		mockListTLSSecrets:   mockListTLSSecrets,
	}
}

func (m mockSecretService) InspectTLSSecret(namespace, name string) (string, error) {
	return m.mockInspectTLSSecret()
}

func (m mockSecretService) ListTLSSecrets(namespace string) ([]Secret, error) {
	return m.mockListTLSSecrets()
}

func (m mockSecretService) ListTLSSecret(namespace, name string) (Secret, error) {
	return Secret{}, nil
}
