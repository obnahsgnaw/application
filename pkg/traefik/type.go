package traefik

const httpMidPrefix = "traefik/http/middlewares"
const tcpMidPrefix = "traefik/tcp/middlewares"
const httpRouterPrefix = "traefik/http/routers"
const TcpRouterPrefix = "traefik/tcp/routers"
const udpRouterPrefix = "traefik/udp/routers"
const httpServicePrefix = "traefik/http/services"
const TcpServicePrefix = "traefik/tcp/services"
const udpServicePrefix = "traefik/udp/services"
const httpTransportPrefix = "traefik/http/serversTransports"

func typRouterPrefix(typ Typ) string {
	switch typ {
	case TypHttp:
		return httpRouterPrefix
	case TypTcp:
		return TcpRouterPrefix
	case TypUdp:
		return udpRouterPrefix
	default:
		panic("not support typ")
	}
}

type Typ string

func (t Typ) String() string {
	return string(t)
}
func (t Typ) IsHttp() bool {
	return t == TypHttp
}
func (t Typ) IsTcp() bool {
	return t == TypTcp
}

const (
	TypHttp Typ = "http"
	TypTcp  Typ = "tcp"
	TypUdp  Typ = "udp"
)
