package secret

type TLS struct {
	name       string
	namespace  string
	secretType string
	tlsCert    []byte
	tlsKey     []byte
}

func NewTLS(name string, namespace string, secretType string, tlsCert []byte, tlsKey []byte) TLS {
	return TLS{
		name:       name,
		namespace:  namespace,
		secretType: secretType,
		tlsCert:    tlsCert,
		tlsKey:     tlsKey,
	}

}

func (t TLS) Name() string {
	return t.name
}

func (t TLS) Namespace() string {
	return t.namespace
}

func (t TLS) PemCert() string {
	return string(t.tlsCert)
}

func (t TLS) PemKey() string {
	return string(t.tlsKey)
}

func (t TLS) Cert() []byte {
	return t.tlsCert
}

func (t TLS) Key() []byte {
	return t.tlsKey
}

func (t TLS) Equals(t2 TLS) bool {
	return t.Name() == t2.Name() && t.Namespace() == t2.Namespace()
}
