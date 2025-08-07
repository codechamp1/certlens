package service

import (
	"fmt"

	"github.com/codechamp1/certlens/internal/domains/cert"
	"github.com/codechamp1/certlens/internal/domains/secret"
)

type Manager interface {
	InspectTLSSecret(tlsSecret secret.TLS) ([]cert.TLS, error)
	ListTLSSecrets(namespace string) ([]secret.TLS, error)
	ListTLSSecret(namespace, name string) (secret.TLS, error)
	RawInspectTLSSecret(namespace, name string) (string, string, error)
}

type defaultManager struct {
	secret.Repository
	cert.Service
}

func NewDefaultManager(sr secret.Repository, cs cert.Service) Manager {
	return defaultManager{
		Repository: sr,
		Service:    cs,
	}
}

func (s defaultManager) InspectTLSSecret(tlsSecret secret.TLS) ([]cert.TLS, error) {
	certData, err := s.ParseTLSCert(tlsSecret.Cert())

	if err != nil {
		return nil, fmt.Errorf("app service can not parse tls cert: %w", err)
	}

	return certData, nil
}

func (s defaultManager) ListTLSSecrets(namespace string) ([]secret.TLS, error) {
	tlsSecrets, err := s.GetTLSSecrets(namespace)

	if err != nil {
		return nil, fmt.Errorf("can not list TLS secrets: %w", err)
	}

	return tlsSecrets, nil
}

func (s defaultManager) ListTLSSecret(namespace, name string) (secret.TLS, error) {
	tlsSecret, err := s.GetTLSSecret(namespace, name)
	if err != nil {
		return secret.TLS{}, fmt.Errorf("failed to get TLS secret %s in namespace %s: %w", name, namespace, err)
	}

	return tlsSecret, nil
}

func (s defaultManager) RawInspectTLSSecret(namespace, name string) (cert string, key string, err error) {
	tlsSecret, err := s.GetTLSSecret(namespace, name)
	if err != nil {
		return "", "", fmt.Errorf("can not inspect TLS secret: %w", err)
	}

	return tlsSecret.PemCert(), tlsSecret.PemKey(), nil
}
