package data_test

import (
	"errors"
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/codechamp1/certlens/internal/data"
	"github.com/codechamp1/certlens/internal/domains/tls"
)

var errTest = errors.New("simulated error")

func TestNewSecretsRepository(t *testing.T) {
	t.Run("Should create a domainservice with the given data", func(t *testing.T) {
		mockClient := data.NewMockSecretsFetcher(nil, nil)
		mockParser := data.NewMockParser(func(tlsCert []byte) ([]tls.Cert, error) { return []tls.Cert{}, nil })
		repo := data.NewDefaultRepository(mockClient, mockParser)
		if repo == nil {
			t.Error("Expected domainservice to be created, but got nil")
		}
		//nolint
		if _, ok := repo.(tls.Repository); !ok {
			t.Error("Expected domainservice to implement Repository interface, but it does not")
		}
	})
}

func TestGetTLSSecrets(t *testing.T) {
	test := []struct {
		name            string
		namespace       string
		secrets         v1.SecretList
		expectedSecrets []tls.Secret
		expectedErr     error
	}{
		{
			name:            "Should return error if can not fetch the secret from the data",
			namespace:       "",
			secrets:         v1.SecretList{},
			expectedSecrets: []tls.Secret{},
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
			expectedSecrets: []tls.Secret{
				tls.NewTLS(
					"tls-secret-1",
					"default",
					"kubernetes.io/tls",
					[]byte("cert-data"),
					[]byte("key-data"),
					[]tls.Cert{},
				),
			},
			expectedErr: nil,
		},
		{
			name:      "Should return no secrets if there are no Secret secrets",
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
			expectedSecrets: []tls.Secret{},
			expectedErr:     nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := data.NewMockSecretsFetcher(
				func(namespace string) (*v1.SecretList, error) {
					return &tt.secrets, tt.expectedErr
				},
				nil,
			)
			mockParser := data.NewMockParser(func(tlsCert []byte) ([]tls.Cert, error) { return []tls.Cert{}, nil })

			repo := data.NewDefaultRepository(mockClient, mockParser)

			secrets, err := repo.GetTLSSecrets(tt.namespace)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}

			// normalize secrets for comparison
			if secrets == nil {
				secrets = []tls.Secret{}
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
		expected    tls.Secret
		expectedErr error
	}{
		{
			name:        "Should return error if can not fetch the secret from the data",
			namespace:   "default",
			secret:      v1.Secret{},
			expected:    tls.Secret{},
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
			expected: tls.NewTLS(
				"tls-secret-1",
				"default",
				"kubernetes.io/tls",
				[]byte("cert-data"),
				[]byte("key-data"),
				[]tls.Cert{},
			),
		},
		{
			name:      "Should return error if the secret is not of type Secret",
			namespace: "default",
			secret: v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "no-tls-secret",
					Namespace: "default",
				},
			},
			expected:    tls.Secret{},
			expectedErr: errTest,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := data.NewMockSecretsFetcher(
				nil,
				func(namespace, name string) (*v1.Secret, error) {
					return &tt.secret, tt.expectedErr
				},
			)
			mockParser := data.NewMockParser(func(tlsCert []byte) ([]tls.Cert, error) { return []tls.Cert{}, nil })

			repo := data.NewDefaultRepository(mockClient, mockParser)

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
