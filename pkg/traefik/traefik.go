package traefik

type Traefik struct {
	typ             Typ
	kvs             map[string]string
	services        map[string]*Service
	routers         map[string]*Router
	httpMiddlewares map[string]HttpMiddleware
	tcpMiddlewares  map[string]TcpMiddleware
	tlsStore        map[string]*TlsStore
	tlsOption       map[string]*TlsOptions
	tlsCert         *TlsCertificates
}

func NewTraefik(typ Typ) *Traefik {
	return &Traefik{
		typ:             typ,
		kvs:             make(map[string]string),
		services:        make(map[string]*Service),
		routers:         make(map[string]*Router),
		httpMiddlewares: make(map[string]HttpMiddleware),
		tcpMiddlewares:  make(map[string]TcpMiddleware),
		tlsStore:        make(map[string]*TlsStore),
		tlsOption:       make(map[string]*TlsOptions),
	}
}

func (t *Traefik) DefineService(service *Service) {
	if service.Type() == t.typ {
		t.services[service.name] = service
	}
}
func (t *Traefik) DefineRouter(router *Router) {
	if t.typ == router.typ {
		t.routers[router.name] = router
	}
}
func (t *Traefik) DefineHttpMiddleware(middleware ...HttpMiddleware) {
	if t.typ == TypHttp {
		for _, mid := range middleware {
			t.httpMiddlewares[mid.Name()] = mid
		}
	}
}
func (t *Traefik) DefineTcpMiddleware(middleware ...TcpMiddleware) {
	if t.typ == TypTcp {
		for _, mid := range middleware {
			t.tcpMiddlewares[mid.Name()] = mid
		}
	}
}
func (t *Traefik) DefineTlsStore(store *TlsStore) {
	t.tlsStore[store.name] = store
}
func (t *Traefik) DefineTlsOption(option *TlsOptions) {
	t.tlsOption[option.name] = option
}
func (t *Traefik) SetTlsCertificate(cert *TlsCertificates) {
	t.tlsCert = cert
}
func (t *Traefik) GetKvs() map[string]string {
	kvs := make(map[string]string)
	kvs = mergeMap(kvs, t.kvs)
	for _, s := range t.services {
		kvs = mergeMap(kvs, s.GetKvs())
	}
	for _, s := range t.routers {
		kvs = mergeMap(kvs, s.GetKvs())
	}
	if t.typ == TypHttp {
		for _, s := range t.httpMiddlewares {
			kvs = mergeMap(kvs, s.GetKvs())
		}
	}
	if t.typ == TypTcp {
		for _, s := range t.tcpMiddlewares {
			kvs = mergeMap(kvs, s.GetKvs())
		}
	}
	if t.tlsCert != nil {
		kvs = mergeMap(kvs, t.tlsCert.GetKvs())
	}
	for _, s := range t.tlsStore {
		kvs = mergeMap(kvs, s.GetKvs())
	}
	for _, s := range t.tlsOption {
		kvs = mergeMap(kvs, s.GetKvs())
	}

	return kvs
}

func mergeMap(m map[string]string, m1 ...map[string]string) map[string]string {
	for _, mp := range m1 {
		for k, v := range mp {
			m[k] = v
		}
	}

	return m
}
