package secret

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"github.com/codechamp1/certlens/internal/client"
)

type Repository interface {
	GetTLSSecrets(namespace string) ([]TLS, error)
	GetTLSSecret(namespace, name string) (TLS, error)
}

type defaultRepository struct {
	client client.SecretsFetcher
}

func NewDefaultRepository(client client.SecretsFetcher) Repository {
	return defaultRepository{
		client: client,
	}
}

func (s defaultRepository) GetTLSSecrets(namespace string) ([]TLS, error) {
	secretsList, err := s.client.FetchSecrets(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get secrets in namespace %s: %w", namespace, err)
	}

	var tlsSecrets []TLS
	for _, secret := range secretsList.Items {
		if secret.Type == corev1.SecretTypeTLS {
			tlsSecrets = append(tlsSecrets, mapSecretToModel(secret))
		}
	}

	return tlsSecrets, nil
}

func (s defaultRepository) GetTLSSecret(namespace, name string) (TLS, error) {
	secret, err := s.client.FetchSecret(namespace, name)

	if err != nil {
		return TLS{}, fmt.Errorf("failed to get secret %s in namespace %s: %w", name, namespace, err)
	}

	if secret.Type != corev1.SecretTypeTLS {
		return TLS{}, fmt.Errorf("secret %s in namespace %s is not of type TLS", name, namespace)
	}

	return mapSecretToModel(*secret), nil
}

func mapSecretToModel(secret corev1.Secret) TLS {
	return NewTLS(
		secret.Name,
		secret.Namespace,
		string(secret.Type),
		secret.Data[corev1.TLSCertKey],
		secret.Data[corev1.TLSPrivateKeyKey],
	)
}
