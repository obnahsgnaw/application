package traefik

import "strconv"

const tlsPrefix = "traefik/tls"

/*

traefik/tls/certificates/0/certFile	foobar
traefik/tls/certificates/0/keyFile	foobar
traefik/tls/certificates/0/stores/0	foobar
traefik/tls/certificates/0/stores/1	foobar
traefik/tls/certificates/1/certFile	foobar
traefik/tls/certificates/1/keyFile	foobar
traefik/tls/certificates/1/stores/0	foobar
traefik/tls/certificates/1/stores/1	foobar

traefik/tls/options/Options0/alpnProtocols/0	foobar
traefik/tls/options/Options0/alpnProtocols/1	foobar
traefik/tls/options/Options0/cipherSuites/0	foobar
traefik/tls/options/Options0/cipherSuites/1	foobar
traefik/tls/options/Options0/clientAuth/caFiles/0	foobar
traefik/tls/options/Options0/clientAuth/caFiles/1	foobar
traefik/tls/options/Options0/clientAuth/clientAuthType	foobar
traefik/tls/options/Options0/curvePreferences/0	foobar
traefik/tls/options/Options0/curvePreferences/1	foobar
traefik/tls/options/Options0/maxVersion	foobar
traefik/tls/options/Options0/minVersion	foobar
traefik/tls/options/Options0/preferServerCipherSuites	true
traefik/tls/options/Options0/sniStrict	true
traefik/tls/options/Options1/alpnProtocols/0	foobar
traefik/tls/options/Options1/alpnProtocols/1	foobar
traefik/tls/options/Options1/cipherSuites/0	foobar
traefik/tls/options/Options1/cipherSuites/1	foobar
traefik/tls/options/Options1/clientAuth/caFiles/0	foobar
traefik/tls/options/Options1/clientAuth/caFiles/1	foobar
traefik/tls/options/Options1/clientAuth/clientAuthType	foobar
traefik/tls/options/Options1/curvePreferences/0	foobar
traefik/tls/options/Options1/curvePreferences/1	foobar
traefik/tls/options/Options1/maxVersion	foobar
traefik/tls/options/Options1/minVersion	foobar
traefik/tls/options/Options1/preferServerCipherSuites	true
traefik/tls/options/Options1/sniStrict	true

traefik/tls/stores/Store0/defaultCertificate/certFile	foobar
traefik/tls/stores/Store0/defaultCertificate/keyFile	foobar
traefik/tls/stores/Store1/defaultCertificate/certFile	foobar
traefik/tls/stores/Store1/defaultCertificate/keyFile	foobar

*/

type TlsCertificate struct {
	CertFile string
	KeyFile  string
	Stores   []string
}

type TlsCertificates struct {
	kvs map[string]string
}

func (t *TlsCertificates) GetKvs() map[string]string {
	return t.kvs
}

func NewTlsCertificates(certificates []TlsCertificate) *TlsCertificates {
	kvs := make(map[string]string)
	for i, v := range certificates {
		kvs[EtcdKey(tlsPrefix, "certificates", strconv.Itoa(i), "certFile")] = v.CertFile
		kvs[EtcdKey(tlsPrefix, "certificates", strconv.Itoa(i), "keyFile")] = v.KeyFile
		for i1, v1 := range v.Stores {
			kvs[EtcdKey(tlsPrefix, "certificates", strconv.Itoa(i), "stores", strconv.Itoa(i1))] = v1
		}
	}

	return &TlsCertificates{kvs: kvs}
}

type TlsOption struct {
	AlpnProtocols            []string
	CipherSuites             []string
	ClientAuthCaFiles        []string
	ClientAuthType           string
	CurvePreferences         []string
	MaxVersion               string
	MinVersion               string
	PreferServerCipherSuites bool
	SniStrict                bool
}
type TlsOptions struct {
	name string
	kvs  map[string]string
}

func (t *TlsOptions) GetKvs() map[string]string {
	return t.kvs
}

func NewOption(name string, option TlsOption) *TlsOptions {
	kvs := make(map[string]string)
	for i, v := range option.AlpnProtocols {
		kvs[EtcdKey(tlsPrefix, "options", name, "alpnProtocols", strconv.Itoa(i))] = v
	}
	for i, v := range option.CipherSuites {
		kvs[EtcdKey(tlsPrefix, "options", name, "cipherSuites", strconv.Itoa(i))] = v
	}
	for i, v := range option.ClientAuthCaFiles {
		kvs[EtcdKey(tlsPrefix, "options", name, "clientAuth/caFiles", strconv.Itoa(i))] = v
	}
	kvs[EtcdKey(tlsPrefix, "options", name, "clientAuth/clientAuthType")] = option.ClientAuthType

	for i, v := range option.CurvePreferences {
		kvs[EtcdKey(tlsPrefix, "options", name, "curvePreferences", strconv.Itoa(i))] = v
	}
	kvs[EtcdKey(tlsPrefix, "options", name, "maxVersion")] = option.MaxVersion
	kvs[EtcdKey(tlsPrefix, "options", name, "minVersion")] = option.MinVersion
	kvs[EtcdKey(tlsPrefix, "options", name, "preferServerCipherSuites")] = BoolVal(option.PreferServerCipherSuites)
	kvs[EtcdKey(tlsPrefix, "options", name, "sniStrict")] = BoolVal(option.SniStrict)

	return &TlsOptions{name: name, kvs: kvs}
}

type TlsStore struct {
	name string
	kvs  map[string]string
}

func (t *TlsStore) GetKvs() map[string]string {
	return t.kvs
}

func NewStore(name, certFile, keyFile string) *TlsStore {
	kvs := make(map[string]string)
	kvs[EtcdKey(tlsPrefix, "stores", name, "defaultCertificate/certFile")] = certFile
	kvs[EtcdKey(tlsPrefix, "stores", name, "defaultCertificate/keyFile")] = keyFile

	return &TlsStore{name: name, kvs: kvs}
}
