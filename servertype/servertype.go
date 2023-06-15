package servertype

type ServerType string

const (
	Api ServerType = "api"
	Rpc ServerType = "rpc"
	Tcp ServerType = "tcp"
	Wss ServerType = "wss"
	Udp ServerType = "udp"
	Hdl ServerType = "hdl"
)

func (s ServerType) String() string {
	return string(s)
}
