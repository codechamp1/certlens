package secret_test

import (
	"errors"
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/codechamp1/certlens/internal/client"
	"github.com/codechamp1/certlens/internal/domains/secret"
)

var errTest = errors.New("simulated error")

func TestNewSecretsRepository(t *testing.T) {
	t.Run("Should create a domainservice with the given client", func(t *testing.T) {
		mockClient := client.NewMockSecretsFetcher(nil, nil)
		repo := secret.NewDefaultRepository(mockClient)
		if repo == nil {
			t.Error("Expected domainservice to be created, but got nil")
		}
		//nolint
		if _, ok := repo.(secret.Repository); !ok {
			t.Error("Expected domainservice to implement Repository interface, but it does not")
		}
	})
}

func TestGetTLSSecrets(t *testing.T) {
	test := []struct {
		name            string
		namespace       string
		secrets         v1.SecretList
		expectedSecrets []secret.TLS
		expectedErr     error
	}{
		{
			name:            "Should return error if can not fetch the secret from the client",
			namespace:       "",
			secrets:         v1.SecretList{},
			expectedSecrets: []secret.TLS{},
			expectedErr:     errTest,
		},
		{
			name:      "Should return secrets if they are fetched successfully",
			namespace: "default",
			secrets: v1.SecretList{
				Items: []v1.Secret{
					{
						Type: v1.SecretTypeTLS,
						ObjectMeta: metav1.ObjectMeta{
							Name:      "tls-secret-1",
							Namespace: "default",
						},
						Data: map[string][]byte{
							v1.TLSCertKey:       []byte("cert-data"),
							v1.TLSPrivateKeyKey: []byte("key-data"),
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "simple-secret-1",
							Namespace: "default",
						},
					},
				},
			},
			expectedSecrets: []secret.TLS{
				secret.NewTLS(
					"tls-secret-1",
					"default",
					"kubernetes.io/tls",
					[]byte("cert-data"),
					[]byte("key-data"),
				),
			},
			expectedErr: nil,
		},
		{
			name:      "Should return no secrets if there are no TLS secrets",
			namespace: "default",
			secrets: v1.SecretList{
				Items: []v1.Secret{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "simple-secret-1",
							Namespace: "default",
						},
					},
				},
			},
			expectedSecrets: []secret.TLS{},
			expectedErr:     nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := client.NewMockSecretsFetcher(
				func(namespace string) (*v1.SecretList, error) {
					return &tt.secrets, tt.expectedErr
				},
				nil,
			)

			repo := secret.NewDefaultRepository(mockClient)

			secrets, err := repo.GetTLSSecrets(tt.namespace)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}

			// normalize secrets for comparison
			if secrets == nil {
				secrets = []secret.TLS{}
			}

			if !reflect.DeepEqual(secrets, tt.expectedSecrets) {
				t.Errorf("Expected secrets %+v, got %+v", tt.expectedSecrets, secrets)
			}
		})
	}
}

func TestGetTLSSecret(t *testing.T) {
	test := []struct {
		name        string
		namespace   string
		secret      v1.Secret
		expected    secret.TLS
		expectedErr error
	}{
		{
			name:        "Should return error if can not fetch the secret from the client",
			namespace:   "default",
			secret:      v1.Secret{},
			expected:    secret.TLS{},
			expectedErr: errTest,
		},
		{
			name:      "Should return tls secret if it is fetched successfully",
			namespace: "default",
			secret: v1.Secret{
				Type: v1.SecretTypeTLS,
				ObjectMeta: metav1.ObjectMeta{
					Name:      "tls-secret-1",
					Namespace: "default",
				},
				Data: map[string][]byte{
					v1.TLSCertKey:       []byte("cert-data"),
					v1.TLSPrivateKeyKey: []byte("key-data"),
				},
			},
			expected: secret.NewTLS(
				"tls-secret-1",
				"default",
				"kubernetes.io/tls",
				[]byte("cert-data"),
				[]byte("key-data"),
			),
		},
		{
			name:      "Should return error if the secret is not of type TLS",
			namespace: "default",
			secret: v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "no-tls-secret",
					Namespace: "default",
				},
			},
			expected:    secret.TLS{},
			expectedErr: errTest,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := client.NewMockSecretsFetcher(
				nil,
				func(namespace, name string) (*v1.Secret, error) {
					return &tt.secret, tt.expectedErr
				},
			)

			repo := secret.NewDefaultRepository(mockClient)

			tlsSecret, err := repo.GetTLSSecret(tt.namespace, tt.secret.Name)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(tlsSecret, tt.expected) {
				t.Errorf("Expected secret %+v, got %+v", tt.expected, tlsSecret)
			}
		})
	}
}
