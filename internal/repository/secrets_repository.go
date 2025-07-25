package repository

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"github.com/codechamp1/certlens/internal/client"
	"github.com/codechamp1/certlens/internal/domains"
)

type SecretsRepository interface {
	GetTLSSecrets(namespace string) ([]domains.SecretInfo, error)
	GetTLSSecret(namespace, name string) (domains.SecretInfo, error)
}

type secretsRepository struct {
	client client.SecretsFetcher
}

func NewSecretsRepository(client client.SecretsFetcher) SecretsRepository {
	return secretsRepository{
		client: client,
	}
}

func (s secretsRepository) GetTLSSecrets(namespace string) ([]domains.SecretInfo, error) {
	secretsList, err := s.client.FetchSecrets(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get secrets in namespace %s: %w", namespace, err)
	}

	var tlsSecrets []domains.SecretInfo
	for _, secret := range secretsList.Items {
		if secret.Type == corev1.SecretTypeTLS {
			tlsSecrets = append(tlsSecrets, mapSecretToModel(secret))
		}
	}

	return tlsSecrets, nil
}

func (s secretsRepository) GetTLSSecret(namespace, name string) (domains.SecretInfo, error) {
	secret, err := s.client.FetchSecret(namespace, name)

	if err != nil {
		return domains.SecretInfo{}, fmt.Errorf("failed to get secret %s in namespace %s: %w", name, namespace, err)
	}

	if secret.Type != corev1.SecretTypeTLS {
		return domains.SecretInfo{}, fmt.Errorf("secret %s in namespace %s is not of type TLS", name, namespace)
	}

	return mapSecretToModel(*secret), nil
}

func mapSecretToModel(secret corev1.Secret) domains.SecretInfo {
	return domains.SecretInfo{
		Name:      secret.Name,
		Namespace: secret.Namespace,
		Type:      string(secret.Type),
		TLSCert:   secret.Data[corev1.TLSCertKey],
		TLSKey:    secret.Data[corev1.TLSPrivateKeyKey],
	}
}
