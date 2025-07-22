package repository_test

import (
	"testing"

	"certlens/internal/client"
	"certlens/internal/repository"
)

func TestNewSecretsRepository(t *testing.T) {
	t.Run("Should create a repository with the given client", func(t *testing.T) {
		mockClient := &client.Client{}
		repo := repository.NewSecretsRepository(mockClient)
		if repo == nil {
			t.Error("expected a repository, got nil")
		}

	})
}
