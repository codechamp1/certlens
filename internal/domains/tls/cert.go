package tls

import "time"

type ExpiryStatus int

const (
	Valid ExpiryStatus = iota
	Warning
	Critical
	Expired
)

var expiryStatusStrings = map[ExpiryStatus]string{
	Valid:    "OK",
	Warning:  "Warning",
	Critical: "Critical",
	Expired:  "Expired",
}

func (s ExpiryStatus) String() string {
	if str, ok := expiryStatusStrings[s]; ok {
		return str
	}
	return "Unknown"
}

type Cert struct {
	CertRawData      `label:"Certificate Raw Info"`
	CertComputedData `label:"Certificate Computed Info"`
}

type CertRawData struct {
	// Raw Info
	Subject            string `label:"Subject"`
	Issuer             string `label:"Issuer"`
	SerialNumber       string `label:"Serial Number"`
	NotBefore          string `label:"Valid From"`
	NotAfter           string `label:"Valid To"`
	Signature          string `label:"Signature"`
	SignatureAlgorithm string `label:"Signature Algorithm"`
	PublicKeyAlgorithm string `label:"Public PemKey Algorithm"`
	IsCA               bool   `label:"Is CA"`

	// Subject Alternative Names
	DNSNames       []string `label:"DNS Names"`
	EmailAddresses []string `label:"Email Addresses"`
	IPAddresses    []string `label:"IP Addresses"`
	URIs           []string `label:"URIs"`

	// PemKey IDs
	SubjectKeyID   string `label:"Subject PemKey ID"`
	AuthorityKeyID string `label:"Authority PemKey ID"`

	// CRL / OCSP
	CRLDistributionPoints []string `label:"CRL Distribution Points"`
	OCSPServers           []string `label:"OCSP Servers"`

	// Usage
	KeyUsage     string   `label:"PemKey Usage"`
	ExtKeyUsages []string `label:"Extended PemKey Usage"`

	// Certificate Version
	Version int `label:"X.509 Version"`
}

type CertComputedData struct {
	Expired             bool          `label:"Expired"`
	TimeUntilExpiry     time.Duration `label:"Time Until Expiry"`
	TotalValidity       time.Duration `label:"Total Validity Duration"`
	TimeSinceIssued     time.Duration `label:"Time Since Issued"`
	ValidityUsedPercent float64       `label:"Validity Used (%)"`
	RemainingPercent    float64       `label:"Time Remaining (%)"`
	ExpiryStatus        ExpiryStatus  `label:"Expiry ExpiryStatus"`
	IsSelfSigned        bool          `label:"Self-Signed"`
	IsCurrentlyValid    bool          `label:"Currently Valid"`
}
