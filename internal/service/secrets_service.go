package service

import (
	"fmt"

	"certlens/internal/repository"
)

type SecretsService interface {
	InspectTLSSecret(namespace, name string) (string, error)
	ListTLSSecrets(namespace string) ([]Secret, error)
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

	certInfo := "Name: " + secret.Name + "\n" +
		"Namespace: " + secret.Namespace + "\n" +
		"Type: " + secret.Type + "\n"

	if len(secret.TLSCert) > 0 {
		certInfo += "TLS Certificate: [REDACTED]\n"
	} else {
		certInfo += "TLS Certificate: Not found\n"
	}

	return certInfo, nil
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
