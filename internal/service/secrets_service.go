package service

import (
	"fmt"

	"certlens/internal/domains"
	"certlens/internal/repository"
)

type SecretsService interface {
	InspectTLSSecret(namespace, name string) (*CertificateInfo, error)
	ListTLSSecrets(namespace string) ([]domains.K8SResourceID, error)
	ListTLSSecret(namespace, name string) (domains.K8SResourceID, error)
	RawInspectTLSSecret(namespace, name string) (string, error)
}

type secretsService struct {
	repository.SecretsRepository
}

func NewSecretsService(repo repository.SecretsRepository) SecretsService {
	return secretsService{
		SecretsRepository: repo,
	}
}

func (s secretsService) InspectTLSSecret(namespace, name string) (*CertificateInfo, error) {
	secret, err := s.GetTLSSecret(namespace, name)
	if err != nil {
		return nil, fmt.Errorf("can not inspect TLS secret: %w", err)
	}

	certData, err := parseCertFromString(string(secret.TLSCert))

	if err != nil {
		return nil, fmt.Errorf("can not parse TLS secret: %w", err)
	}

	parsedCert := parseCertificate(*certData)

	return &parsedCert, nil
}

func (s secretsService) ListTLSSecrets(namespace string) ([]domains.K8SResourceID, error) {

	secrets, err := s.GetTLSSecrets(namespace)

	if err != nil {
		return nil, fmt.Errorf("can not list TLS secrets: %w", err)
	}

	var tlsSecretsNames []domains.K8SResourceID
	for _, secret := range secrets {
		tlsSecretsNames = append(tlsSecretsNames, domains.K8SResourceID{secret.Name, secret.Namespace})
	}

	return tlsSecretsNames, nil
}

func (s secretsService) ListTLSSecret(namespace, name string) (domains.K8SResourceID, error) {
	secret, err := s.SecretsRepository.GetTLSSecret(namespace, name)
	if err != nil {
		return domains.K8SResourceID{}, fmt.Errorf("failed to get TLS secret %s in namespace %s: %w", name, namespace, err)
	}

	return domains.K8SResourceID{secret.Name, secret.Namespace}, nil
}

func (s secretsService) RawInspectTLSSecret(namespace, name string) (string, error) {
	secret, err := s.GetTLSSecret(namespace, name)
	if err != nil {
		return "", fmt.Errorf("can not inspect TLS secret: %w", err)
	}

	return string(secret.TLSCert), nil
}
