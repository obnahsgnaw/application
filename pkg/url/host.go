package url

import (
	"github.com/obnahsgnaw/application/pkg/utils"
	"strconv"
)

type Host struct {
	Ip   string
	Port int
}

func New(ip string, port int) Host {
	return Host{
		Ip:   ip,
		Port: port,
	}
}

func (h Host) String() string {
	if h.Ip == "" && h.Port == 0 {
		return ""
	}
	return utils.ToStr(h.Ip, ":", strconv.Itoa(h.Port))
}
