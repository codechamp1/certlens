package service

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"time"
)

var keyUsageNames = map[x509.KeyUsage]string{
	x509.KeyUsageDigitalSignature:  "Digital Signature",
	x509.KeyUsageContentCommitment: "Content Commitment",
	x509.KeyUsageKeyEncipherment:   "Key Encipherment",
	x509.KeyUsageDataEncipherment:  "Data Encipherment",
	x509.KeyUsageKeyAgreement:      "Key Agreement",
	x509.KeyUsageCertSign:          "Cert Sign",
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

	// Required fields
	sb.WriteString(fmt.Sprintf("Subject: %s\n", cert.Subject.String()))
	sb.WriteString(fmt.Sprintf("Issuer: %s\n", cert.Issuer.String()))
	sb.WriteString(fmt.Sprintf("Serial Number: %s\n", cert.SerialNumber.String()))
	sb.WriteString(fmt.Sprintf("Valid From: %s\n", cert.NotBefore.Format(time.RFC1123)))
	sb.WriteString(fmt.Sprintf("Valid To: %s\n", cert.NotAfter.Format(time.RFC1123)))
	sb.WriteString(fmt.Sprintf("Signature Algorithm: %s\n", cert.SignatureAlgorithm.String()))
	sb.WriteString(fmt.Sprintf("Public Key Algorithm: %s\n", cert.PublicKeyAlgorithm.String()))
	sb.WriteString(fmt.Sprintf("Is CA: %t\n", cert.IsCA))

	// Optional fields (printed only if non-empty)
	if len(cert.DNSNames) > 0 {
		sb.WriteString(fmt.Sprintf("DNS Names: %s\n", strings.Join(cert.DNSNames, ", ")))
	}
	if len(cert.EmailAddresses) > 0 {
		sb.WriteString(fmt.Sprintf("Email Addresses: %s\n", strings.Join(cert.EmailAddresses, ", ")))
	}
	if len(cert.IPAddresses) > 0 {
		var ips []string
		for _, ip := range cert.IPAddresses {
			ips = append(ips, ip.String())
		}
		sb.WriteString(fmt.Sprintf("IP Addresses: %s\n", strings.Join(ips, ", ")))
	}
	if len(cert.URIs) > 0 {
		var uris []string
		for _, uri := range cert.URIs {
			uris = append(uris, uri.String())
		}
		sb.WriteString(fmt.Sprintf("URIs: %s\n", strings.Join(uris, ", ")))
	}
	if cert.SubjectKeyId != nil {
		sb.WriteString(fmt.Sprintf("Subject Key ID: %X\n", cert.SubjectKeyId))
	}
	if cert.AuthorityKeyId != nil {
		sb.WriteString(fmt.Sprintf("Authority Key ID: %X\n", cert.AuthorityKeyId))
	}
	if len(cert.CRLDistributionPoints) > 0 {
		sb.WriteString(fmt.Sprintf("CRL Distribution Points: %s\n", strings.Join(cert.CRLDistributionPoints, ", ")))
	}
	if len(cert.OCSPServer) > 0 {
		sb.WriteString(fmt.Sprintf("OCSP Servers: %s\n", strings.Join(cert.OCSPServer, ", ")))
	}

	if cert.KeyUsage != 0 {
		sb.WriteString(fmt.Sprintf("Key Usage: %s\n", keyUsageToString(cert.KeyUsage)))
	}

	if len(cert.ExtKeyUsage) > 0 {
		var usages []string
		for _, usage := range cert.ExtKeyUsage {
			usages = append(usages, extKeyUsageToString(usage))
		}
		sb.WriteString(fmt.Sprintf("Extended Key Usage: %s\n", strings.Join(usages, ", ")))
	}
	return sb.String()
}

func extKeyUsageToString(eku x509.ExtKeyUsage) string {
	if s, ok := extKeyUsageMap[eku]; ok {
		return s
	}
	return fmt.Sprintf("Unknown (%d)", eku)
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
