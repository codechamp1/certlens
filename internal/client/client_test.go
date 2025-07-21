package client_test

import (
	"os"
	"testing"

	"certlens/internal/client"
)

const kubeconfigContent = `
apiVersion: v1
kind: Config
contexts:
- name: valid-context
  context:
    cluster: my-cluster
    user: my-user
clusters:
- name: my-cluster
  cluster:
    server: https://localhost
users:
- name: my-user
  user:
    token: dummytoken
current-context: valid-context
`

func TestNewClient(t *testing.T) {
	t.Run("Should return an error if somethinf goes wrong", func(t *testing.T) {
		t.Run("Should return an error if kubeconfig is not found", func(t *testing.T) {
			_, err := client.NewClient("non_existent_kubeconfig", "test_context")
			if err == nil {
				t.Error("expected an error, got nil")
			}
		})
		t.Run("Should return an error if context is not found", func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "kubeconfig")
			if err != nil {
				t.Fatalf("failed to create temp kubeconfig file: %v", err)
			}
			defer os.Remove(tmpFile.Name())
			if _, err := tmpFile.WriteString(kubeconfigContent); err != nil {
				t.Fatalf("failed to write kubeconfig: %v", err)
			}
			tmpFile.Close()
			_, err = client.NewClient(tmpFile.Name(), "non_existent_context")
			if err == nil {
				t.Error("expected an error, got nil")
			}
		})
	})
}
