package service

type mockSecretService struct {
}

func NewMockSecretService() SecretsService {
	return mockSecretService{}
}

func (m mockSecretService) InspectTLSSecret(namespace, name string) (string, error) {
	return "", nil
}

func (m mockSecretService) ListTLSSecrets(namespace string) ([]Secret, error) {
	return []Secret{
		{"Test1", "Test2"},
		{"Test3", "Test4"},
		{"Test3", "Test4"},
		{"Test3", "Test4"},
		{"Test3", "Test4"},
		{"Test3", "Test4"},
		{"Test3", "Test4"},
		{"Test3", "Test4"},
		{"Test3", "Test4"},
	}, nil
}
