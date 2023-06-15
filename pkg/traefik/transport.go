package traefik

import "strconv"

type Certificate struct {
	CertFile string
	KeyFile  string
}
type ForwardingTimeouts struct {
	DialTimeout           int
	IdleConnTimeout       int
	PingTimeout           int
	ReadIdleTimeout       int
	ResponseHeaderTimeout int
}

type Transport struct {
	name string
	kvs  map[string]string
}

func (s *Transport) Name() string {
	return s.name
}

/*

traefik/http/serversTransports/ServersTransport0/certificates/0/certFile	foobar
traefik/http/serversTransports/ServersTransport0/certificates/0/keyFile	foobar
traefik/http/serversTransports/ServersTransport0/certificates/1/certFile	foobar
traefik/http/serversTransports/ServersTransport0/certificates/1/keyFile	foobar
traefik/http/serversTransports/ServersTransport0/disableHTTP2	true
traefik/http/serversTransports/ServersTransport0/forwardingTimeouts/dialTimeout	42s
traefik/http/serversTransports/ServersTransport0/forwardingTimeouts/idleConnTimeout	42s
traefik/http/serversTransports/ServersTransport0/forwardingTimeouts/pingTimeout	42s
traefik/http/serversTransports/ServersTransport0/forwardingTimeouts/readIdleTimeout	42s
traefik/http/serversTransports/ServersTransport0/forwardingTimeouts/responseHeaderTimeout	42s
traefik/http/serversTransports/ServersTransport0/insecureSkipVerify	true
traefik/http/serversTransports/ServersTransport0/maxIdleConnsPerHost	42
traefik/http/serversTransports/ServersTransport0/peerCertURI	foobar
traefik/http/serversTransports/ServersTransport0/rootCAs/0	foobar
traefik/http/serversTransports/ServersTransport0/rootCAs/1	foobar
traefik/http/serversTransports/ServersTransport0/serverName	foobar

*/

func NewTransport(transportName, serverName string, disableHTTP2, insecureSkipVerify bool, maxIdleConnsPerHost int, peerCertURI string) *Transport {
	kvs := make(map[string]string)
	if serverName != "" {
		kvs[EtcdKey(httpTransportPrefix, transportName, "serverName")] = serverName
	}
	kvs[EtcdKey(httpTransportPrefix, transportName, "disableHTTP2")] = BoolVal(disableHTTP2)
	kvs[EtcdKey(httpTransportPrefix, transportName, "insecureSkipVerify")] = BoolVal(insecureSkipVerify)
	if maxIdleConnsPerHost > 0 {
		kvs[EtcdKey(httpTransportPrefix, transportName, "maxIdleConnsPerHost")] = strconv.Itoa(maxIdleConnsPerHost)
	}
	if peerCertURI != "" {
		kvs[EtcdKey(httpTransportPrefix, transportName, "peerCertURI")] = peerCertURI
	}

	return &Transport{
		name: transportName,
		kvs:  kvs,
	}
}

func (s *Transport) SetRootCAs(rootCAs []string) {
	for i, v := range rootCAs {
		s.kvs[EtcdKey(httpTransportPrefix, s.name, "rootCAs", strconv.Itoa(i))] = v
	}
}
func (s *Transport) SetCertificate(certificates []*Certificate) {
	for i, v := range certificates {
		s.kvs[EtcdKey(httpTransportPrefix, s.name, "certificates", strconv.Itoa(i), "certFile")] = v.CertFile
		s.kvs[EtcdKey(httpTransportPrefix, s.name, "certificates", strconv.Itoa(i), "keyFile")] = v.KeyFile
	}
}
func (s *Transport) SetForwardingTimeouts(timeouts ForwardingTimeouts) {
	s.kvs[EtcdKey(httpTransportPrefix, s.name, "forwardingTimeouts/dialTimeout")] = strconv.Itoa(timeouts.DialTimeout) + "s"
	s.kvs[EtcdKey(httpTransportPrefix, s.name, "forwardingTimeouts/idleConnTimeout")] = strconv.Itoa(timeouts.IdleConnTimeout) + "s"
	s.kvs[EtcdKey(httpTransportPrefix, s.name, "forwardingTimeouts/pingTimeout")] = strconv.Itoa(timeouts.PingTimeout) + "s"
	s.kvs[EtcdKey(httpTransportPrefix, s.name, "forwardingTimeouts/readIdleTimeout")] = strconv.Itoa(timeouts.ReadIdleTimeout) + "s"
	s.kvs[EtcdKey(httpTransportPrefix, s.name, "forwardingTimeouts/responseHeaderTimeout")] = strconv.Itoa(timeouts.ResponseHeaderTimeout) + "s"
}
