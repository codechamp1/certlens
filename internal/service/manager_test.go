package service_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/codechamp1/certlens/internal/domains/tls"
	"github.com/codechamp1/certlens/internal/service"
)

var errRepo = errors.New("simulated error")

func TestNewSecretsService(t *testing.T) {
	mockRepo := service.NewMockRepository(nil, nil)
	svc := service.NewDefaultManager(mockRepo)
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
		secrets           []tls.Secret
		expectedSecretIDs []tls.Secret
		expectedRepoErr   error
	}{
		{
			name:              "Should return error if can not fetch secrets",
			namespace:         "",
			secrets:           []tls.Secret{},
			expectedSecretIDs: []tls.Secret{},
			expectedRepoErr:   errRepo,
		},
		{
			name:      "Should transform all Secret secrets secret Secret",
			namespace: "default",
			secrets: []tls.Secret{
				tls.NewTLS(
					"tls-secret-1",
					"default",
					"kubernetes/tls",
					[]byte("cert-data"),
					[]byte("key-data"),
					[]tls.Cert{},
				),
				tls.NewTLS(
					"tls-secret-2",
					"default",
					"kubernetes/tls",
					[]byte("cert-data"),
					[]byte("key-data"),
					[]tls.Cert{},
				),
			},
			expectedSecretIDs: []tls.Secret{
				tls.NewTLS(
					"tls-secret-1",
					"default",
					"kubernetes/tls",
					[]byte("cert-data"),
					[]byte("key-data"),
					[]tls.Cert{},
				),
				tls.NewTLS(
					"tls-secret-2",
					"default",
					"kubernetes/tls",
					[]byte("cert-data"),
					[]byte("key-data"),
					[]tls.Cert{},
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := service.NewMockRepository(func(namespace string) ([]tls.Secret, error) {
				return tt.secrets, tt.expectedRepoErr
			}, nil)

			svc := service.NewDefaultManager(mockRepo)
			secrets, err := svc.ListTLSSecrets(tt.namespace)

			if !errors.Is(err, tt.expectedRepoErr) {
				t.Errorf("expected error %v, got %v", tt.expectedRepoErr, err)
			}

			if secrets == nil {
				secrets = []tls.Secret{}
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
		secret           tls.Secret
		expectedSecretID tls.Secret
		expectedRepoErr  error
	}{
		{
			name:             "Should return error if can not fetch secret",
			namespace:        "default",
			secret:           tls.Secret{},
			expectedSecretID: tls.Secret{},
			expectedRepoErr:  errRepo,
		},
		{
			name:      "Should secret Secret for a single Secret secret",
			namespace: "default",
			secret: tls.NewTLS(
				"tls-secret-1",
				"default",
				"kubernetes/tls",
				[]byte("cert-data"),
				[]byte("key-data"),
				[]tls.Cert{},
			),
			expectedSecretID: tls.NewTLS(
				"tls-secret-1",
				"default",
				"kubernetes/tls",
				[]byte("cert-data"),
				[]byte("key-data"),
				[]tls.Cert{},
			),
			expectedRepoErr: nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := service.NewMockRepository(nil, func(namespace, name string) (tls.Secret, error) {
				return tt.secret, tt.expectedRepoErr
			})

			svc := service.NewDefaultManager(mockRepo)
			secretID, err := svc.LoadTLSSecret(tt.namespace, tt.secret.Name())

			if !errors.Is(err, tt.expectedRepoErr) {
				t.Errorf("expected error %v, got %v", tt.expectedRepoErr, err)
			}

			if !reflect.DeepEqual(secretID, tt.expectedSecretID) {
				t.Errorf("expected secret ID %v, got %v", tt.expectedSecretID, secretID)
			}
		})
	}
}
