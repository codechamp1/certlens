package service_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/codechamp1/certlens/internal/domains/cert"
	"github.com/codechamp1/certlens/internal/domains/secret"
	"github.com/codechamp1/certlens/internal/service"
)

var errRepo = errors.New("simulated error")

func TestNewSecretsService(t *testing.T) {
	mockRepo := secret.NewMockRepository(nil, nil)
	mockService := cert.NewMockService(nil)
	svc := service.NewDefaultManager(mockRepo, mockService)
	if svc == nil {
		t.Error("secrets service should not be nil")
	}
	//nolint
	if _, ok := svc.(service.Manager); !ok {
		t.Error("secrets service should implement Manager interface")
	}

}

func TestListTLSSecrets(t *testing.T) {
	tests := []struct {
		name              string
		namespace         string
		secrets           []secret.TLS
		expectedSecretIDs []secret.TLS
		expectedRepoErr   error
	}{
		{
			name:              "Should return error if can not fetch secrets",
			namespace:         "",
			secrets:           []secret.TLS{},
			expectedSecretIDs: []secret.TLS{},
			expectedRepoErr:   errRepo,
		},
		{
			name:      "Should transform all TLS secrets secret TLS",
			namespace: "default",
			secrets: []secret.TLS{
				secret.NewTLS(
					"tls-secret-1",
					"default",
					"kubernetes/tls",
					[]byte("cert-data"),
					[]byte("key-data"),
				),
				secret.NewTLS(
					"tls-secret-2",
					"default",
					"kubernetes/tls",
					[]byte("cert-data"),
					[]byte("key-data"),
				),
			},
			expectedSecretIDs: []secret.TLS{
				secret.NewTLS(
					"tls-secret-1",
					"default",
					"kubernetes/tls",
					[]byte("cert-data"),
					[]byte("key-data"),
				),
				secret.NewTLS(
					"tls-secret-2",
					"default",
					"kubernetes/tls",
					[]byte("cert-data"),
					[]byte("key-data"),
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := secret.NewMockRepository(func(namespace string) ([]secret.TLS, error) {
				return tt.secrets, tt.expectedRepoErr
			}, nil)

			svc := service.NewDefaultManager(mockRepo, nil)
			secrets, err := svc.ListTLSSecrets(tt.namespace)

			if !errors.Is(err, tt.expectedRepoErr) {
				t.Errorf("expected error %v, got %v", tt.expectedRepoErr, err)
			}

			if secrets == nil {
				secrets = []secret.TLS{}
			}

			if !reflect.DeepEqual(secrets, tt.expectedSecretIDs) {
				t.Errorf("expected secrets %v, got %v", tt.expectedSecretIDs, secrets)
			}
		})
	}
}

func TestListTLSSecret(t *testing.T) {
	test := []struct {
		name             string
		namespace        string
		secret           secret.TLS
		expectedSecretID secret.TLS
		expectedRepoErr  error
	}{
		{
			name:             "Should return error if can not fetch secret",
			namespace:        "default",
			secret:           secret.TLS{},
			expectedSecretID: secret.TLS{},
			expectedRepoErr:  errRepo,
		},
		{
			name:      "Should secret TLS for a single TLS secret",
			namespace: "default",
			secret: secret.NewTLS(
				"tls-secret-1",
				"default",
				"kubernetes/tls",
				[]byte("cert-data"),
				[]byte("key-data"),
			),
			expectedSecretID: secret.NewTLS(
				"tls-secret-1",
				"default",
				"kubernetes/tls",
				[]byte("cert-data"),
				[]byte("key-data"),
			),
			expectedRepoErr: nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := secret.NewMockRepository(nil, func(namespace, name string) (secret.TLS, error) {
				return tt.secret, tt.expectedRepoErr
			})

			svc := service.NewDefaultManager(mockRepo, nil)
			secretID, err := svc.ListTLSSecret(tt.namespace, tt.secret.Name())

			if !errors.Is(err, tt.expectedRepoErr) {
				t.Errorf("expected error %v, got %v", tt.expectedRepoErr, err)
			}

			if !reflect.DeepEqual(secretID, tt.expectedSecretID) {
				t.Errorf("expected secret ID %v, got %v", tt.expectedSecretID, secretID)
			}
		})
	}
}

func TestRawInspectTLSSecret(t *testing.T) {
	tests := []struct {
		name            string
		namespace       string
		secretName      string
		secret          secret.TLS
		expectedCert    string
		expectedKey     string
		expectedRepoErr error
	}{
		{
			name:            "Should return error if can not fetch secret",
			namespace:       "",
			secretName:      "",
			secret:          secret.TLS{},
			expectedCert:    "",
			expectedKey:     "",
			expectedRepoErr: errRepo,
		},
		{
			name:      "Should return raw TLS cert data",
			namespace: "default",
			secret: secret.NewTLS(
				"tls-secret-1",
				"default",
				"kubernetes/tls",
				[]byte("cert-data"),
				[]byte("key-data"),
			),
			expectedCert:    "cert-data",
			expectedKey:     "key-data",
			expectedRepoErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := secret.NewMockRepository(nil, func(namespace, name string) (secret.TLS, error) {
				return tt.secret, tt.expectedRepoErr
			})

			svc := service.NewDefaultManager(mockRepo, nil)
			cert, key, err := svc.RawInspectTLSSecret(tt.namespace, tt.secretName)

			if !errors.Is(err, tt.expectedRepoErr) {
				t.Errorf("expected error %v, got %v", tt.expectedRepoErr, err)
			}

			if cert != tt.expectedCert {
				t.Errorf("expected cert %s, got %s", tt.expectedCert, cert)
			}

			if key != tt.expectedKey {
				t.Errorf("expected key %s, got %s", tt.expectedCert, cert)
			}
		})
	}
}
