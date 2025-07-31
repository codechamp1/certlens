package domains

type SecretInfo struct {
	Name      string
	Namespace string
	Type      string
	TLSCert   []byte
	TLSKey    []byte
}

type K8SResourceID struct {
	Name      string
	Namespace string
}
