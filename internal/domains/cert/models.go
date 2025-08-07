package cert

import "time"

type TLS struct {
	TLSRawData      `label:"Certificate Raw Info"`
	TLSComputedData `label:"Certificate Computed Info"`
}

type TLSRawData struct {
	// Raw Info
	Subject            string `label:"Subject"`
	Issuer             string `label:"Issuer"`
	SerialNumber       string `label:"Serial Number"`
	NotBefore          string `label:"Valid From"`
	NotAfter           string `label:"Valid To"`
	Signature          string `label:"Signature"`
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

	// Certificate Version
	Version int `label:"X.509 Version"`
}

type TLSComputedData struct {
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
