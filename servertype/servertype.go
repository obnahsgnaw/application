package servertype

type ServerType string

const (
	Api    ServerType = "api"
	Rpc    ServerType = "rpc"
	Tcp    ServerType = "tcp"
	Wss    ServerType = "wss"
	Udp    ServerType = "udp"
	TcpHdl ServerType = "tcp-hdl"
	WssHdl ServerType = "wss-hdl"
	UdpHdl ServerType = "upd-hdl"
)

func (s ServerType) String() string {
	return string(s)
}
