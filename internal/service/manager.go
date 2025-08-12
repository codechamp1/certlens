package service

import (
	"fmt"

	"github.com/codechamp1/certlens/internal/domains/tls"
)

type Manager interface {
	ListTLSSecrets(namespace string) ([]tls.Secret, error)
	LoadTLSSecret(namespace, name string) (tls.Secret, error)
}

type defaultManager struct {
	tls.Repository
}

func NewDefaultManager(tr tls.Repository) Manager {
	return defaultManager{
		Repository: tr,
	}
}

func (s defaultManager) ListTLSSecrets(namespace string) ([]tls.Secret, error) {
	tlsSecrets, err := s.GetTLSSecrets(namespace)

	if err != nil {
		return nil, fmt.Errorf("can not list TLS secrets: %w", err)
	}

	return tlsSecrets, nil
}

func (s defaultManager) LoadTLSSecret(namespace, name string) (tls.Secret, error) {
	tlsSecret, err := s.GetTLSSecret(namespace, name)
	if err != nil {
		return tls.Secret{}, fmt.Errorf("failed to get TLS  secret %s in namespace %s: %w", name, namespace, err)
	}

	return tlsSecret, nil
}
