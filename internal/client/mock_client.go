package client

type mockSecretClient struct{}

func NewMockSecretsClient() *SecretsFetcher {
	return mockClient{}
}
