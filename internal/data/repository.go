package data

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"github.com/codechamp1/certlens/internal/domains/tls"
)

type defaultRepository struct {
	SecretsFetcher
	CertParser
}

func NewDefaultRepository(c SecretsFetcher, cp CertParser) tls.Repository {
	return defaultRepository{
		SecretsFetcher: c,
		CertParser:     cp,
	}
}

func (s defaultRepository) GetTLSSecrets(namespace string) ([]tls.Secret, error) {
	secretsList, err := s.FetchSecrets(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get secrets in namespace %s: %w", namespace, err)
	}

	var tlsSecrets []tls.Secret
	for _, secret := range secretsList.Items {
		if secret.Type == corev1.SecretTypeTLS {
			certs, err := s.ParseTLSCert(secret.Data[corev1.TLSCertKey])
			if err != nil {
				return nil, fmt.Errorf("failed to parse tls cert in the repository %s/%s: %w", secret.Namespace, secret.Name, err)
			}
			tlsSecrets = append(tlsSecrets, mapToModel(secret, certs))
		}
	}

	return tlsSecrets, nil
}

func (s defaultRepository) GetTLSSecret(namespace, name string) (tls.Secret, error) {
	secret, err := s.FetchSecret(namespace, name)

	if err != nil {
		return tls.Secret{}, fmt.Errorf("failed to get secret %s in namespace %s: %w", name, namespace, err)
	}

	if secret.Type != corev1.SecretTypeTLS {
		return tls.Secret{}, fmt.Errorf("secret %s in namespace %s is not of type Secret", name, namespace)
	}

	certs, err := s.ParseTLSCert(secret.Data[corev1.TLSCertKey])
	if err != nil {
		return tls.Secret{}, fmt.Errorf("failed to parse tls cert in the repository %s/%s: %w", secret.Namespace, secret.Name, err)
	}

	return mapToModel(*secret, certs), nil
}

func mapToModel(secret corev1.Secret, certs []tls.Cert) tls.Secret {
	return tls.NewTLS(
		secret.Name,
		secret.Namespace,
		string(secret.Type),
		secret.Data[corev1.TLSCertKey],
		secret.Data[corev1.TLSPrivateKeyKey],
		certs,
	)
}
