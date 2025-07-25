package service_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/codechamp1/certlens/internal/domains"
	"github.com/codechamp1/certlens/internal/repository"
	"github.com/codechamp1/certlens/internal/service"
)

var errRepo = errors.New("simulated error")

func TestNewSecretsService(t *testing.T) {
	mockRepo := repository.NewMockRepository(nil, nil)
	svc := service.NewSecretsService(mockRepo)
	if svc == nil {
		t.Error("secrets service should not be nil")
	}
	if _, ok := svc.(service.SecretsService); !ok {
		t.Error("secrets service should implement SecretsService interface")
	}
}

func TestListTLSSecrets(t *testing.T) {
	tests := []struct {
		name              string
		namespace         string
		secrets           []domains.SecretInfo
		expectedSecretIDs []domains.K8SResourceID
		expectedRepoErr   error
	}{
		{
			name:              "Should return error if can not fetch secrets",
			namespace:         "",
			secrets:           []domains.SecretInfo{},
			expectedSecretIDs: []domains.K8SResourceID{},
			expectedRepoErr:   errRepo,
		},
		{
			name:      "Should transform all TLS secrets to K8SResourceID",
			namespace: "default",
			secrets: []domains.SecretInfo{
				{
					Name:      "tls-secret-1",
					Namespace: "default",
					Type:      "kubernetes/tls",
					TLSCert:   []byte("cert-data"),
					TLSKey:    []byte("key-data"),
				},
				{
					Name:      "tls-secret-2",
					Namespace: "default",
					Type:      "kubernetes/tls",
					TLSCert:   []byte("cert-data"),
					TLSKey:    []byte("key-data"),
				},
			},
			expectedSecretIDs: []domains.K8SResourceID{
				{Name: "tls-secret-1", Namespace: "default"},
				{Name: "tls-secret-2", Namespace: "default"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := repository.NewMockRepository(func(namespace string) ([]domains.SecretInfo, error) {
				return tt.secrets, tt.expectedRepoErr
			}, nil)

			svc := service.NewSecretsService(mockRepo)
			secrets, err := svc.ListTLSSecrets(tt.namespace)

			if !errors.Is(err, tt.expectedRepoErr) {
				t.Errorf("expected error %v, got %v", tt.expectedRepoErr, err)
			}

			if secrets == nil {
				secrets = []domains.K8SResourceID{}
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
		secret           domains.SecretInfo
		expectedSecretID domains.K8SResourceID
		expectedRepoErr  error
	}{
		{
			name:             "Should return error if can not fetch secret",
			namespace:        "default",
			secret:           domains.SecretInfo{},
			expectedSecretID: domains.K8SResourceID{},
			expectedRepoErr:  errRepo,
		},
		{
			name:      "Should return K8SResourceID for a single TLS secret",
			namespace: "default",
			secret: domains.SecretInfo{
				Name:      "tls-secret-1",
				Namespace: "default",
				Type:      "kubernetes/tls",
				TLSCert:   []byte("cert-data"),
				TLSKey:    []byte("key-data"),
			},
			expectedSecretID: domains.K8SResourceID{Name: "tls-secret-1", Namespace: "default"},
			expectedRepoErr:  nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := repository.NewMockRepository(nil, func(namespace, name string) (domains.SecretInfo, error) {
				return tt.secret, tt.expectedRepoErr
			})

			svc := service.NewSecretsService(mockRepo)
			secretID, err := svc.ListTLSSecret(tt.namespace, tt.secret.Name)

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
		secret          domains.SecretInfo
		expectedCert    string
		expectedKey     string
		expectedRepoErr error
	}{
		{
			name:            "Should return error if can not fetch secret",
			namespace:       "",
			secretName:      "",
			secret:          domains.SecretInfo{},
			expectedCert:    "",
			expectedKey:     "",
			expectedRepoErr: errRepo,
		},
		{
			name:      "Should return raw TLS certificate data",
			namespace: "default",
			secret: domains.SecretInfo{
				Name:      "tls-secret-1",
				Namespace: "default",
				Type:      "kubernetes/tls",
				TLSCert:   []byte("cert-data"),
				TLSKey:    []byte("key-data"),
			},
			expectedCert:    "cert-data",
			expectedKey:     "key-data",
			expectedRepoErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := repository.NewMockRepository(nil, func(namespace, name string) (domains.SecretInfo, error) {
				return tt.secret, tt.expectedRepoErr
			})

			svc := service.NewSecretsService(mockRepo)
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
