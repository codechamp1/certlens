package service

import (
	"fmt"

	"certlens/internal/repository"
)

type SecretsService interface {
	InspectTLSSecret(namespace, name string) (string, error)
	ListTLSSecrets(namespace string) ([]Secret, error)
	ListTLSSecret(namespace, name string) (Secret, error)
}

type Secret struct {
	Name      string
	Namespace string
}

type secretsService struct {
	repository.SecretsRepository
}

func NewSecretsService(repo repository.SecretsRepository) SecretsService {
	return secretsService{
		SecretsRepository: repo,
	}
}

func (s secretsService) InspectTLSSecret(namespace, name string) (string, error) {
	secret, err := s.GetTLSSecret(namespace, name)
	if err != nil {
		return "", fmt.Errorf("can not inspect TLS secret: %w", err)
	}

	certData, err := parseCertFromString(string(secret.TLSCert))

	if err != nil {
		return "", fmt.Errorf("can not parse TLS secret: %w", err)
	}

	return formatCertInfo(certData), nil
}

func (s secretsService) ListTLSSecrets(namespace string) ([]Secret, error) {

	secrets, err := s.GetTLSSecrets(namespace)

	if err != nil {
		return nil, fmt.Errorf("can not list TLS secrets: %w", err)
	}

	var tlsSecretsNames []Secret
	for _, secret := range secrets {
		tlsSecretsNames = append(tlsSecretsNames, Secret{secret.Name, secret.Namespace})
	}

	return tlsSecretsNames, nil
}

func (s secretsService) ListTLSSecret(namespace, name string) (Secret, error) {
	secret, err := s.SecretsRepository.GetTLSSecret(namespace, name)
	if err != nil {
		return Secret{}, fmt.Errorf("failed to get TLS secret %s in namespace %s: %w", name, namespace, err)
	}

	return Secret{secret.Name, secret.Namespace}, nil
}
