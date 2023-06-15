package regtype

type RegType string

const (
	Http RegType = "http"
	Rpc  RegType = "rpc"
	Doc  RegType = "doc"
	Tcp  RegType = "tcp"
	Wss  RegType = "wss"
	Udp  RegType = "udp"
)

func (s RegType) String() string {
	return string(s)
}
