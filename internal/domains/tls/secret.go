package tls

type Secret struct {
	name       string
	namespace  string
	secretType string
	pemCert    []byte
	pemKey     []byte
	certs      []Cert
}

type Repository interface {
	GetTLSSecrets(namespace string) ([]Secret, error)
	GetTLSSecret(namespace, name string) (Secret, error)
}

func NewTLS(name string, namespace string, secretType string, pem []byte, key []byte, certs []Cert) Secret {
	return Secret{
		name:       name,
		namespace:  namespace,
		secretType: secretType,
		pemCert:    pem,
		pemKey:     key,
		certs:      certs,
	}

}

func (t Secret) Name() string {
	return t.name
}

func (t Secret) Namespace() string {
	return t.namespace
}

func (t Secret) PemCert() string {
	return string(t.pemCert)
}

func (t Secret) PemKey() string {
	return string(t.pemKey)
}

func (t Secret) Certs() []Cert {
	return t.certs
}

func (t Secret) Equals(t2 Secret) bool {
	return t.Name() == t2.Name() && t.Namespace() == t2.Namespace()
}
