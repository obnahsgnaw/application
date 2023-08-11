package traefik

import "strconv"

type Service struct {
	name string
	typ  Typ
	kvs  map[string]string
}

func (s *Service) Name() string {
	return s.name
}
func (s *Service) Type() Typ {
	return s.typ
}

/*

traefik/http/services/Service01/loadBalancer/healthCheck/followRedirects	true
traefik/http/services/Service01/loadBalancer/healthCheck/headers/name0	foobar
traefik/http/services/Service01/loadBalancer/healthCheck/headers/name1	foobar
traefik/http/services/Service01/loadBalancer/healthCheck/hostname	foobar
traefik/http/services/Service01/loadBalancer/healthCheck/interval	foobar
traefik/http/services/Service01/loadBalancer/healthCheck/path	foobar
traefik/http/services/Service01/loadBalancer/healthCheck/port	42
traefik/http/services/Service01/loadBalancer/healthCheck/scheme	foobar
traefik/http/services/Service01/loadBalancer/healthCheck/timeout	foobar
traefik/http/services/Service01/loadBalancer/passHostHeader	true
traefik/http/services/Service01/loadBalancer/responseForwarding/flushInterval	foobar
traefik/http/services/Service01/loadBalancer/servers/0/url	foobar
traefik/http/services/Service01/loadBalancer/servers/1/url	foobar
traefik/http/services/Service01/loadBalancer/serversTransport	foobar
traefik/http/services/Service01/loadBalancer/sticky/cookie/httpOnly	true
traefik/http/services/Service01/loadBalancer/sticky/cookie/name	foobar
traefik/http/services/Service01/loadBalancer/sticky/cookie/sameSite	foobar
traefik/http/services/Service01/loadBalancer/sticky/cookie/secure	true

*/

func NewHttpService(serviceName string, serverHosts []string, passHostHeader bool, flushInterval int, serversTransport string) *Service {
	kvs := make(map[string]string)
	kvs[EtcdKey(httpServicePrefix, serviceName, "loadBalancer/passHostHeader")] = BoolVal(passHostHeader)
	if flushInterval > 0 {
		kvs[EtcdKey(httpServicePrefix, serviceName, "loadBalancer/responseForwarding/flushInterval")] = strconv.Itoa(flushInterval)
	}
	if serversTransport != "" {
		kvs[EtcdKey(httpServicePrefix, serviceName, "loadBalancer/serversTransport")] = serversTransport
	}
	for i, v := range serverHosts {
		if v != "" {
			kvs[EtcdKey(httpServicePrefix, serviceName, "loadBalancer/servers", strconv.Itoa(i), "url")] = v
		}
	}

	return &Service{
		name: serviceName,
		typ:  TypHttp,
		kvs:  kvs,
	}
}

func NewHttpServiceServerKey(serviceName string, index int) string {
	return EtcdKey(httpServicePrefix, serviceName, "loadBalancer/servers", strconv.Itoa(index), "url")
}

/*
traefik/tcp/services/TCPService01/loadBalancer/proxyProtocol/version	42
traefik/tcp/services/TCPService01/loadBalancer/servers/0/address	foobar
traefik/tcp/services/TCPService01/loadBalancer/servers/1/address	foobar
traefik/tcp/services/TCPService01/loadBalancer/terminationDelay	42
traefik/tcp/services/TCPService02/weighted/services/0/name	foobar
traefik/tcp/services/TCPService02/weighted/services/0/weight	42
traefik/tcp/services/TCPService02/weighted/services/1/name	foobar
traefik/tcp/services/TCPService02/weighted/services/1/weight	42
*/

func NewTcpService(serviceName string, serverHosts []string, proxyProtocolVersion, terminationDelay int, weightService []WeightService) *Service {
	kvs := make(map[string]string)
	if proxyProtocolVersion > 0 {
		kvs[EtcdKey(TcpServicePrefix, serviceName, "loadBalancer/proxyProtocol/version")] = strconv.Itoa(proxyProtocolVersion)
	}
	if terminationDelay > 0 {
		kvs[EtcdKey(TcpServicePrefix, serviceName, "loadBalancer/terminationDelay")] = strconv.Itoa(terminationDelay)
	}
	for i, v := range serverHosts {
		if v != "" {
			kvs[EtcdKey(TcpServicePrefix, serviceName, "loadBalancer/servers", strconv.Itoa(i), "address")] = v
		}
	}
	if len(weightService) > 0 {
		for i, v := range weightService {
			if v.Name != "" {
				kvs[EtcdKey(TcpServicePrefix, serviceName, "weighted/servers", strconv.Itoa(i), "name")] = v.Name
				kvs[EtcdKey(TcpServicePrefix, serviceName, "weighted/servers", strconv.Itoa(i), "weight")] = strconv.Itoa(v.weight)
			}
		}
	}
	return &Service{
		name: serviceName,
		typ:  TypTcp,
		kvs:  kvs,
	}
}
func NewTcpServiceServerKey(serviceName string, index int) string {
	return EtcdKey(TcpServicePrefix, serviceName, "loadBalancer/servers", strconv.Itoa(index), "address")
}

/*
traefik/udp/services/UDPService01/loadBalancer/servers/0/address	foobar
traefik/udp/services/UDPService01/loadBalancer/servers/1/address	foobar
traefik/udp/services/UDPService02/weighted/services/0/name	foobar
traefik/udp/services/UDPService02/weighted/services/0/weight	42
traefik/udp/services/UDPService02/weighted/services/1/name	foobar
traefik/udp/services/UDPService02/weighted/services/1/weight
*/

func NewUdpService(serviceName string, serverHosts []string, weightService []WeightService) *Service {
	kvs := make(map[string]string)
	for i, v := range serverHosts {
		kvs[EtcdKey(udpServicePrefix, serviceName, "loadBalancer/servers", strconv.Itoa(i), "address")] = v
	}
	if len(weightService) > 0 {
		for i, v := range weightService {
			if v.Name != "" {
				kvs[EtcdKey(udpServicePrefix, serviceName, "weighted/servers", strconv.Itoa(i), "name")] = v.Name
				kvs[EtcdKey(udpServicePrefix, serviceName, "weighted/servers", strconv.Itoa(i), "weight")] = strconv.Itoa(v.weight)
			}
		}
	}
	return &Service{
		name: serviceName,
		typ:  TypTcp,
		kvs:  kvs,
	}
}

func NewUdpServiceServerKey(serviceName string, index int) string {
	return EtcdKey(udpServicePrefix, serviceName, "loadBalancer/servers", strconv.Itoa(index), "address")
}

func (s *Service) SetHealthCheck(followRedirects bool, headers map[string]string, hostname string, interval int, path string, port int, scheme string, timeout int) {
	if s.typ != TypHttp {
		return
	}
	if hostname != "" && path != "" && port > 0 && scheme != "" {
		s.kvs[EtcdKey(httpServicePrefix, s.name, "loadBalancer/healthCheck/followRedirects")] = BoolVal(followRedirects)
		s.kvs[EtcdKey(httpServicePrefix, s.name, "loadBalancer/healthCheck/hostname")] = hostname
		s.kvs[EtcdKey(httpServicePrefix, s.name, "loadBalancer/healthCheck/interval")] = strconv.Itoa(interval)
		s.kvs[EtcdKey(httpServicePrefix, s.name, "loadBalancer/healthCheck/path")] = path
		s.kvs[EtcdKey(httpServicePrefix, s.name, "loadBalancer/healthCheck/port")] = strconv.Itoa(port)
		s.kvs[EtcdKey(httpServicePrefix, s.name, "loadBalancer/healthCheck/scheme")] = scheme
		if timeout > 0 {
			s.kvs[EtcdKey(httpServicePrefix, s.name, "loadBalancer/healthCheck/timeout")] = strconv.Itoa(timeout)
		}
		for i, ss := range headers {
			s.kvs[EtcdKey(httpServicePrefix, s.name, "loadBalancer/healthCheck/headers", i)] = ss
		}
	}
}

type Cookie struct {
	HttpOnly bool
	Name     string
	SameSite string
	Secure   bool
}

func (s *Service) SetCookie(cookie Cookie) {
	if s.typ != TypHttp {
		return
	}
	s.kvs[EtcdKey(httpServicePrefix, s.name, "loadBalancer/sticky/cookie/httpOnly")] = BoolVal(cookie.HttpOnly)
	s.kvs[EtcdKey(httpServicePrefix, s.name, "loadBalancer/sticky/cookie/name")] = cookie.Name
	s.kvs[EtcdKey(httpServicePrefix, s.name, "loadBalancer/sticky/cookie/sameSite")] = cookie.SameSite
	s.kvs[EtcdKey(httpServicePrefix, s.name, "loadBalancer/sticky/cookie/secure")] = BoolVal(cookie.Secure)
}

/*
traefik/http/services/Service02/mirroring/healthCheck	``
traefik/http/services/Service02/mirroring/maxBodySize	42
traefik/http/services/Service02/mirroring/mirrors/0/name	foobar
traefik/http/services/Service02/mirroring/mirrors/0/percent	42
traefik/http/services/Service02/mirroring/mirrors/1/name	foobar
traefik/http/services/Service02/mirroring/mirrors/1/percent	42
traefik/http/services/Service02/mirroring/service	foobar
*/

type Mirror struct {
	Name    string
	Percent int
}

func (s *Service) SetMirroring(healthCheck, service string, maxBodySize int, mirrors []Mirror) {
	if s.typ != TypHttp {
		return
	}
	s.kvs[EtcdKey(httpServicePrefix, s.name, "mirroring/healthCheck")] = healthCheck
	s.kvs[EtcdKey(httpServicePrefix, s.name, "mirroring/maxBodySize")] = strconv.Itoa(maxBodySize)
	s.kvs[EtcdKey(httpServicePrefix, s.name, "mirroring/service")] = service
	for i, mirror := range mirrors {
		s.kvs[EtcdKey(httpServicePrefix, s.name, "mirroring/mirrors", strconv.Itoa(i), "name")] = mirror.Name
		s.kvs[EtcdKey(httpServicePrefix, s.name, "mirroring/mirrors", strconv.Itoa(i), "percent")] = strconv.Itoa(mirror.Percent)
	}
}

/*
traefik/http/services/Service04/failover/fallback	foobar
traefik/http/services/Service04/failover/healthCheck	``
traefik/http/services/Service04/failover/service	foobar
*/

func (s *Service) SetFailover(fallback, healthCheck, service string) {
	if s.typ != TypHttp {
		return
	}
	s.kvs[EtcdKey(httpServicePrefix, s.name, "failover/healthCheck")] = healthCheck
	s.kvs[EtcdKey(httpServicePrefix, s.name, "failover/fallback")] = fallback
	s.kvs[EtcdKey(httpServicePrefix, s.name, "failover/service")] = service
}

/*
traefik/http/services/Service03/weighted/healthCheck	``
traefik/http/services/Service03/weighted/services/0/name	foobar
traefik/http/services/Service03/weighted/services/0/weight	42
traefik/http/services/Service03/weighted/services/1/name	foobar
traefik/http/services/Service03/weighted/services/1/weight	42
traefik/http/services/Service03/weighted/sticky/cookie/httpOnly	true
traefik/http/services/Service03/weighted/sticky/cookie/name	foobar
traefik/http/services/Service03/weighted/sticky/cookie/sameSite	foobar
traefik/http/services/Service03/weighted/sticky/cookie/secure	true
*/

type WeightService struct {
	Name   string
	weight int
}

func (s *Service) SetWeightService(healthCheck string, services []WeightService, cookie *Cookie) {
	if s.typ != TypHttp {
		return
	}
	s.kvs[EtcdKey(httpServicePrefix, s.name, "weighted/healthCheck")] = healthCheck
	for i, ss := range services {
		s.kvs[EtcdKey(httpServicePrefix, s.name, "weighted/services", strconv.Itoa(i), "name")] = ss.Name
		s.kvs[EtcdKey(httpServicePrefix, s.name, "weighted/services", strconv.Itoa(i), "weight")] = strconv.Itoa(ss.weight)
	}
	if cookie != nil {
		s.kvs[EtcdKey(httpServicePrefix, s.name, "weighted/sticky/cookie/httpOnly")] = BoolVal(cookie.HttpOnly)
		s.kvs[EtcdKey(httpServicePrefix, s.name, "weighted/sticky/cookie/name")] = cookie.Name
		s.kvs[EtcdKey(httpServicePrefix, s.name, "weighted/sticky/cookie/sameSite")] = cookie.SameSite
		s.kvs[EtcdKey(httpServicePrefix, s.name, "weighted/sticky/cookie/secure")] = BoolVal(cookie.Secure)
	}
}

func (s *Service) GetKvs() map[string]string {
	return s.kvs
}
