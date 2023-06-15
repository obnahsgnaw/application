package traefik

import (
	"strconv"
)

/*

traefik/http/routers/Router0/entryPoints/0	foobar
traefik/http/routers/Router0/entryPoints/1	foobar
traefik/http/routers/Router0/middlewares/0	foobar
traefik/http/routers/Router0/middlewares/1	foobar
traefik/http/routers/Router0/priority	42
traefik/http/routers/Router0/rule	foobar
traefik/http/routers/Router0/service	foobar
traefik/http/routers/Router0/tls/certResolver	foobar
traefik/http/routers/Router0/tls/domains/0/main	foobar
traefik/http/routers/Router0/tls/domains/0/sans/0	foobar
traefik/http/routers/Router0/tls/domains/0/sans/1	foobar
traefik/http/routers/Router0/tls/domains/1/main	foobar
traefik/http/routers/Router0/tls/domains/1/sans/0	foobar
traefik/http/routers/Router0/tls/domains/1/sans/1	foobar
traefik/http/routers/Router0/tls/options	foobar

*/

type Router struct {
	name            string
	typ             Typ
	kvs             map[string]string
	middlewares     []string
	tlsSet          bool
	tlsCertResolver string
	tlsOption       string
	tlsDomains      []*TlsDomain
	tlsPassThrough  bool
}

type TlsDomain struct {
	Main string
	Sans []string
}

func newRouter(typ Typ, name, rule, serviceName string, entryPoints []string, priority int) *Router {
	kvs := make(map[string]string)
	for i, entry := range entryPoints {
		kvs[EtcdKey(typRouterPrefix(typ), name, "entryPoints", strconv.Itoa(i))] = entry
	}
	if priority > 0 {
		kvs[EtcdKey(typRouterPrefix(typ), name, "priority")] = strconv.Itoa(priority)
	}
	if rule != "" {
		kvs[EtcdKey(typRouterPrefix(typ), name, "rule")] = rule
	}
	kvs[EtcdKey(typRouterPrefix(typ), name, "service")] = serviceName
	return &Router{
		name: name,
		typ:  typ,
		kvs:  kvs,
	}
}

func NewHttpRouter(name, rule, serviceName string, entryPoints []string, priority int) *Router {
	return newRouter(TypHttp, name, rule, serviceName, entryPoints, priority)
}

func NewTcpRouter(name, rule, serviceName string, entryPoints []string, priority int) *Router {
	return newRouter(TypTcp, name, rule, serviceName, entryPoints, priority)
}

func NewUdpRouter(name, serviceName string, entryPoints []string) *Router {
	return newRouter(TypUdp, name, "", serviceName, entryPoints, 0)
}

func (r *Router) AddMiddlewares(middleware ...string) {
	if r.typ != TypUdp {
		r.middlewares = append(r.middlewares, middleware...)
	}
}

func (r *Router) parseMiddlewares() {
	for i, middleware := range r.middlewares {
		r.kvs[EtcdKey(typRouterPrefix(r.typ), r.name, "middlewares", strconv.Itoa(i))] = middleware
	}
	return
}

func (r *Router) SetTls(certResolver, options string, domains []*TlsDomain, passThrough bool) {
	if r.typ != TypUdp {
		r.tlsSet = true
		r.tlsCertResolver = certResolver
		r.tlsOption = options
		r.tlsDomains = domains
		r.tlsPassThrough = passThrough
	}
}

func (r *Router) parseTls() {
	if !r.tlsSet {
		return
	}
	if r.tlsCertResolver != "" {
		r.kvs[EtcdKey(typRouterPrefix(r.typ), r.name, "tls/certResolver")] = r.tlsCertResolver
	}
	if r.tlsOption != "" {
		r.kvs[EtcdKey(typRouterPrefix(r.typ), r.name, "tls/options")] = r.tlsOption
	}
	for i, domain := range r.tlsDomains {
		r.kvs[EtcdKey(typRouterPrefix(r.typ), r.name, "tls/domains", strconv.Itoa(i), "main")] = domain.Main
		for i1, san := range domain.Sans {
			r.kvs[EtcdKey(typRouterPrefix(r.typ), r.name, "tls/domains", strconv.Itoa(i), "sans", strconv.Itoa(i1))] = san
		}
	}
	r.kvs[EtcdKey(typRouterPrefix(r.typ), r.name, "tls/passthrough")] = BoolVal(r.tlsPassThrough)
	return
}

func (r *Router) GetKvs() map[string]string {
	r.parseMiddlewares()
	r.parseTls()
	return r.kvs
}

func (r *Router) Type() Typ {
	return r.typ
}

func (r *Router) Name() string {
	return r.name
}
