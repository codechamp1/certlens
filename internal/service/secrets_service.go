package service

import (
	"fmt"

	"github.com/codechamp1/certlens/internal/domains"
	"github.com/codechamp1/certlens/internal/repository"
)

type SecretsService interface {
	InspectTLSSecret(namespace, name string) ([]CertificateInfo, error)
	ListTLSSecrets(namespace string) ([]domains.SecretInfo, error)
	ListTLSSecret(namespace, name string) (domains.SecretInfo, error)
	RawInspectTLSSecret(namespace, name string) (string, string, error)
}

type secretsService struct {
	repository.SecretsRepository
}

func NewSecretsService(repo repository.SecretsRepository) SecretsService {
	return secretsService{
		SecretsRepository: repo,
	}
}

func (s secretsService) InspectTLSSecret(namespace, name string) ([]CertificateInfo, error) {
	secret, err := s.GetTLSSecret(namespace, name)
	if err != nil {
		return nil, fmt.Errorf("can not inspect TLS secret: %w", err)
	}

	certData, err := ParseCertsFromString(string(secret.TLSCert))

	if err != nil {
		return nil, fmt.Errorf("can not parse TLS secret: %w", err)
	}

	parsedCert := parseCertificates(certData)

	return parsedCert, nil
}

func (s secretsService) ListTLSSecrets(namespace string) ([]domains.SecretInfo, error) {
	secrets, err := s.GetTLSSecrets(namespace)

	if err != nil {
		return nil, fmt.Errorf("can not list TLS secrets: %w", err)
	}
	return secrets, nil
}

func (s secretsService) ListTLSSecret(namespace, name string) (domains.SecretInfo, error) {
	secret, err := s.GetTLSSecret(namespace, name)
	if err != nil {
		return domains.SecretInfo{}, fmt.Errorf("failed to get TLS secret %s in namespace %s: %w", name, namespace, err)
	}
	return secret, nil
}

func (s secretsService) RawInspectTLSSecret(namespace, name string) (cert string, key string, err error) {
	secret, err := s.GetTLSSecret(namespace, name)
	if err != nil {
		return "", "", fmt.Errorf("can not inspect TLS secret: %w", err)
	}

	return string(secret.TLSCert), string(secret.TLSKey), nil
}
