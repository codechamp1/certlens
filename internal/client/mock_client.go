package client

import corev1 "k8s.io/api/core/v1"

type mockSecretsFetcher struct {
	mockFetchSecrets func(namespace string) (*corev1.SecretList, error)
	mockFetchSecret  func(namespace, name string) (*corev1.Secret, error)
}

func NewSecretsFetcher(
	mockGetTLSSecrets func(namespace string) (*corev1.SecretList, error),
	mockGetTLSSecret func(namespace, name string) (*corev1.Secret, error),
) SecretsFetcher {
	return mockSecretsFetcher{
		mockFetchSecrets: mockGetTLSSecrets,
		mockFetchSecret:  mockGetTLSSecret,
	}
}

func (m mockSecretsFetcher) FetchSecrets(namespace string) (*corev1.SecretList, error) {
	return m.mockFetchSecrets(namespace)
}

func (m mockSecretsFetcher) FetchSecret(namespace, name string) (*corev1.Secret, error) {
	return m.mockFetchSecret(namespace, name)
}
