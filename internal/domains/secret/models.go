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

func (s TLS) Name() string {
	return s.name
}

func (s TLS) Namespace() string {
	return s.namespace
}

func (s TLS) PemCert() string {
	return string(s.tlsCert)
}

func (s TLS) PemKey() string {
	return string(s.tlsKey)
}

func (s TLS) Cert() []byte {
	return s.tlsCert
}

func (s TLS) Key() []byte {
	return s.tlsKey
}
