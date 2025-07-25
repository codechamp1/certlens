package client

import (
	"testing"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8sTesting "k8s.io/client-go/testing"
)

var errTest = errors.New("simulated error")

func TestFetchSecrets(t *testing.T) {
	tests := []struct {
		name        string
		namespace   string
		secrets     []runtime.Object
		expectedErr error
	}{
		{
			name:        "Should return error if the client fails to fetch secrets",
			namespace:   "",
			secrets:     []runtime.Object{},
			expectedErr: errTest,
		},
		{
			name:      "Should return error if the client fails to fetch secrets",
			namespace: "default",
			secrets: []runtime.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "secret",
						Namespace: "default",
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k8sClient := fake.NewClientset(tt.secrets...)
			if tt.expectedErr != nil {
				k8sClient.PrependReactor("list", "secrets", func(action k8sTesting.Action) (bool, runtime.Object, error) {
					return true, nil, tt.expectedErr
				})
			}

			client := &Client{k8sClient}
			secrets, err := client.FetchSecrets(tt.namespace)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if tt.expectedErr == nil && len(secrets.Items) != len(tt.secrets) {
				t.Errorf("expected %d secrets, got %d", len(tt.secrets), len(secrets.Items))
			}
		})
	}
}

func TestClient_FetchSecret(t *testing.T) {
	tests := []struct {
		name        string
		namespace   string
		secret      runtime.Object
		secretName  string
		expectedErr error
	}{
		{
			name:        "Should return error if the client fails to fetch secret",
			namespace:   "default",
			secret:      &corev1.Secret{},
			secretName:  "",
			expectedErr: errTest,
		},
		{
			name:      "Should the secret",
			namespace: "default",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "default",
				},
			},
			secretName:  "secret",
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k8sClient := fake.NewClientset(tt.secret)
			if tt.expectedErr != nil {
				k8sClient.PrependReactor("get", "secrets", func(action k8sTesting.Action) (bool, runtime.Object, error) {
					return true, nil, tt.expectedErr
				})
			}

			client := &Client{k8sClient}
			secret, err := client.FetchSecret(tt.namespace, tt.secretName)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if tt.expectedErr == nil && secret == nil {
				t.Error("expected a secret, got nil")
			}
		})
	}
}
