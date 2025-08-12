package data

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/codechamp1/certlens/internal/domains/tls"
)

type defaultParser struct{}

type CertParser interface {
	ParseTLSCert(tlsCert []byte) ([]tls.Cert, error)
}

func NewDefaultParser() CertParser {
	return defaultParser{}
}

func (cs defaultParser) ParseTLSCert(tlsCert []byte) ([]tls.Cert, error) {
	x509Certs, err := parseCertsFromString(string(tlsCert))
	if err != nil {
		return []tls.Cert{}, fmt.Errorf("service can not parse tls cert, err: %w", err)
	}

	return parseCertificateChain(x509Certs), nil
}

var keyUsageNames = map[x509.KeyUsage]string{
	x509.KeyUsageDigitalSignature:  "Digital Signature",
	x509.KeyUsageContentCommitment: "Content Commitment",
	x509.KeyUsageKeyEncipherment:   "PemKey Encipherment",
	x509.KeyUsageDataEncipherment:  "Secret Encipherment",
	x509.KeyUsageKeyAgreement:      "PemKey Agreement",
	x509.KeyUsageCertSign:          "Secret Sign",
	x509.KeyUsageCRLSign:           "CRL Sign",
	x509.KeyUsageEncipherOnly:      "Encipher Only",
	x509.KeyUsageDecipherOnly:      "Decipher Only",
}

var extKeyUsageMap = map[x509.ExtKeyUsage]string{
	x509.ExtKeyUsageAny:             "Any",
	x509.ExtKeyUsageServerAuth:      "Server Auth",
	x509.ExtKeyUsageClientAuth:      "Client Auth",
	x509.ExtKeyUsageCodeSigning:     "Code Signing",
	x509.ExtKeyUsageEmailProtection: "Email Protection",
	x509.ExtKeyUsageTimeStamping:    "Timestamping",
	x509.ExtKeyUsageOCSPSigning:     "OCSP Signing",
}

func fromX509(cert x509.Certificate) tls.Cert {
	percent, status := expiryStatusByPercentage(cert, 25, 10)
	return tls.Cert{
		CertRawData: tls.CertRawData{
			Subject:               cert.Subject.String(),
			Issuer:                cert.Issuer.String(),
			SerialNumber:          cert.SerialNumber.String(),
			NotBefore:             cert.NotBefore.Format(time.RFC1123),
			NotAfter:              cert.NotAfter.Format(time.RFC1123),
			Signature:             fmt.Sprintf("%X", cert.Signature),
			SignatureAlgorithm:    cert.SignatureAlgorithm.String(),
			PublicKeyAlgorithm:    cert.PublicKeyAlgorithm.String(),
			IsCA:                  cert.IsCA,
			DNSNames:              cert.DNSNames,
			EmailAddresses:        cert.EmailAddresses,
			IPAddresses:           joinToStringSlice(cert.IPAddresses, func(ip net.IP) string { return ip.String() }),
			URIs:                  joinToStringSlice(cert.URIs, func(uri *url.URL) string { return uri.String() }),
			SubjectKeyID:          fmt.Sprintf("%X", cert.SubjectKeyId),
			AuthorityKeyID:        fmt.Sprintf("%X", cert.AuthorityKeyId),
			CRLDistributionPoints: cert.CRLDistributionPoints,
			OCSPServers:           cert.OCSPServer,
			KeyUsage:              keyUsageToString(cert.KeyUsage),
			ExtKeyUsages:          extractExtendedKeyUsages(cert),
			Version:               cert.Version,
		},
		CertComputedData: tls.CertComputedData{
			Expired:             time.Now().After(cert.NotAfter),
			TimeUntilExpiry:     time.Until(cert.NotAfter),
			TotalValidity:       cert.NotAfter.Sub(cert.NotBefore),
			TimeSinceIssued:     time.Since(cert.NotBefore),
			ValidityUsedPercent: float64(time.Since(cert.NotBefore)) / float64(cert.NotAfter.Sub(cert.NotBefore)) * 100,
			RemainingPercent:    percent,
			ExpiryStatus:        status,
			IsSelfSigned:        cert.CheckSignatureFrom(&cert) == nil,
			IsCurrentlyValid:    !time.Now().After(cert.NotAfter) && time.Now().After(cert.NotBefore),
		},
	}

}

func parseCertsFromString(pemStr string) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate
	data := []byte(pemStr)

	for {
		block, rest := pem.Decode(data)
		if block == nil {
			break // no more blocks
		}

		if block.Type != "CERTIFICATE" {
			data = rest
			continue // skip non-cert blocks
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse x509 cert: %w", err)
		}

		certs = append(certs, cert)
		data = rest
	}

	if len(certs) == 0 {
		return nil, fmt.Errorf("no certificates found in input")
	}

	return certs, nil
}

func parseCertificateChain(certs []*x509.Certificate) []tls.Cert {
	var certData []tls.Cert
	for _, cert := range certs {
		certData = append(certData, fromX509(*cert))
	}
	return certData
}

func extractExtendedKeyUsages(cert x509.Certificate) []string {
	var usages []string
	for _, usage := range cert.ExtKeyUsage {
		if s, ok := extKeyUsageMap[usage]; ok {
			usages = append(usages, s)
		} else {
			usages = append(usages, fmt.Sprintf("Unknown (%d)", usage))
		}
	}
	return usages
}

func joinToStringSlice[T any](items []T, toStr func(T) string) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		result = append(result, toStr(item))
	}
	return result
}

func keyUsageToString(ku x509.KeyUsage) string {
	var usages []string
	for bit, name := range keyUsageNames {
		if ku&bit != 0 {
			usages = append(usages, name)
		}
	}
	return strings.Join(usages, ", ")
}

func expiryStatusByPercentage(cert x509.Certificate, warningThreshold, criticalThreshold float64) (percentRemaining float64, status tls.ExpiryStatus) {
	now := time.Now()
	validityDuration := cert.NotAfter.Sub(cert.NotBefore)
	timeRemaining := cert.NotAfter.Sub(now)

	if validityDuration <= 0 {
		return 0, tls.Critical
	}

	percentRemaining = float64(timeRemaining) / float64(validityDuration) * 100

	if now.After(cert.NotAfter) {
		status = tls.Expired
	} else if percentRemaining <= criticalThreshold {
		status = tls.Critical
	} else if percentRemaining <= warningThreshold {
		status = tls.Warning
	} else {
		status = tls.Valid
	}

	return percentRemaining, status
}
