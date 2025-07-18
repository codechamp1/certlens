package service

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

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

func parseCertFromString(pemStr string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("failed to decode PEM block containing certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse x509 certificate: %w", err)
	}

	return cert, nil
}

func formatCertInfo(cert *x509.Certificate) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Subject: %s\n", cert.Subject.String()))
	sb.WriteString(fmt.Sprintf("Issuer: %s\n", cert.Issuer.String()))
	sb.WriteString(fmt.Sprintf("Serial Number: %s\n", cert.SerialNumber.String()))
	sb.WriteString(fmt.Sprintf("Not Before: %s\n", cert.NotBefore.Format(time.UnixDate)))
	sb.WriteString(fmt.Sprintf("Not After: %s\n", cert.NotAfter.Format(time.UnixDate)))
	sb.WriteString(fmt.Sprintf("Is CA: %t\n", cert.IsCA))
	sb.WriteString(fmt.Sprintf("DNS Names: %v\n", cert.DNSNames))
	sb.WriteString(fmt.Sprintf("Email Addresses: %v\n", cert.EmailAddresses))
	sb.WriteString(fmt.Sprintf("IP Addresses: %v\n", cert.IPAddresses))
	sb.WriteString(fmt.Sprintf("Signature Algorithm: %v\n", cert.SignatureAlgorithm))
	sb.WriteString(fmt.Sprintf("Public Key Algorithm: %v\n", cert.PublicKeyAlgorithm))

	// You can add more fields as needed

	return sb.String()
}
