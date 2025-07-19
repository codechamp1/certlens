package service

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

type CertificateInfo struct {
	CertificateRawInfo      `label:"Certificate Raw Info"`
	CertificateComputedInfo `label:"Certificate Computed Info"`
}

type CertificateRawInfo struct {
	// Raw Info
	Subject            string `label:"Subject"`
	Issuer             string `label:"Issuer"`
	SerialNumber       string `label:"Serial Number"`
	NotBefore          string `label:"Valid From"`
	NotAfter           string `label:"Valid To"`
	SignatureAlgorithm string `label:"Signature Algorithm"`
	PublicKeyAlgorithm string `label:"Public Key Algorithm"`
	IsCA               bool   `label:"Is CA"`

	// Subject Alternative Names
	DNSNames       []string `label:"DNS Names"`
	EmailAddresses []string `label:"Email Addresses"`
	IPAddresses    []string `label:"IP Addresses"`
	URIs           []string `label:"URIs"`

	// Key IDs
	SubjectKeyID   string `label:"Subject Key ID"`
	AuthorityKeyID string `label:"Authority Key ID"`

	// CRL / OCSP
	CRLDistributionPoints []string `label:"CRL Distribution Points"`
	OCSPServers           []string `label:"OCSP Servers"`

	// Usage
	KeyUsage     string   `label:"Key Usage"`
	ExtKeyUsages []string `label:"Extended Key Usage"`
}

type CertificateComputedInfo struct {
	Expired             bool          `label:"Expired"`
	TimeUntilExpiry     time.Duration `label:"Time Until Expiry"`
	TotalValidity       time.Duration `label:"Total Validity Duration"`
	TimeSinceIssued     time.Duration `label:"Time Since Issued"`
	ValidityUsedPercent float64       `label:"Validity Used (%)"`
	RemainingPercent    float64       `label:"Time Remaining (%)"`
	ExpiryStatus        string        `label:"Expiry Status"`
	IsSelfSigned        bool          `label:"Self-Signed"`
	IsCurrentlyValid    bool          `label:"Currently Valid"`
}

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

type Status int

const (
	valid Status = iota
	warning
	critical
	expired
)

var expiryStatusStrings = map[Status]string{
	valid:    "OK",
	warning:  "Warning",
	critical: "Critical",
	expired:  "Expired",
}

func (s Status) String() string {
	if str, ok := expiryStatusStrings[s]; ok {
		return str
	}
	return "Unknown"
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

func parseCertificate(cert x509.Certificate) CertificateInfo {
	percent, status := expiryStatusByPercentage(cert, 25.0, 10.0) // warning at 25%, critical at 10%
	return CertificateInfo{
		CertificateRawInfo: CertificateRawInfo{
			Subject:               cert.Subject.String(),
			Issuer:                cert.Issuer.String(),
			SerialNumber:          cert.SerialNumber.String(),
			NotBefore:             cert.NotBefore.Format(time.RFC1123),
			NotAfter:              cert.NotAfter.Format(time.RFC1123),
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
		},
		CertificateComputedInfo: CertificateComputedInfo{
			Expired:             time.Now().After(cert.NotAfter),
			TimeUntilExpiry:     cert.NotAfter.Sub(time.Now()),
			TotalValidity:       cert.NotAfter.Sub(cert.NotBefore),
			TimeSinceIssued:     time.Now().Sub(cert.NotBefore),
			ValidityUsedPercent: float64(time.Now().Sub(cert.NotBefore)) / float64(cert.NotAfter.Sub(cert.NotBefore)) * 100,
			RemainingPercent:    percent,
			ExpiryStatus:        status.String(),
			IsSelfSigned:        cert.CheckSignatureFrom(&cert) == nil,
			IsCurrentlyValid:    !time.Now().After(cert.NotAfter) && time.Now().After(cert.NotBefore),
		},
	}
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

func expiryStatusByPercentage(cert x509.Certificate, warningThreshold, criticalThreshold float64) (percentRemaining float64, status Status) {
	now := time.Now()
	validityDuration := cert.NotAfter.Sub(cert.NotBefore)
	timeRemaining := cert.NotAfter.Sub(now)

	if validityDuration <= 0 {
		return 0, critical
	}

	percentRemaining = float64(timeRemaining) / float64(validityDuration) * 100

	if now.After(cert.NotAfter) {
		status = expired
	} else if percentRemaining <= criticalThreshold {
		status = critical
	} else if percentRemaining <= warningThreshold {
		status = warning
	} else {
		status = valid
	}

	return percentRemaining, status
}
