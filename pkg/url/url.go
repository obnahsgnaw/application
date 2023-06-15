package url

import (
	"strconv"
	"strings"
)

func ParseUrl(url string) (protocol Protocol, host Host, path string, ok bool) {
	if strings.Index(url, "://") == -1 {
		ok = false
		return
	}
	pathSeg := strings.Split(url, "://")
	protocol = Protocol(pathSeg[0])
	hostEnd := strings.Index(pathSeg[1], "/")
	var hostStr string
	if hostEnd == -1 {
		hostStr = pathSeg[1]
		path = ""
	} else {
		hostStr = pathSeg[1][0:hostEnd]
		path = pathSeg[1][hostEnd:]
	}
	if strings.Index(hostStr, ":") == -1 {
		host = Host{
			Ip:   hostStr,
			Port: 80,
		}
	} else {
		hostSeg := strings.Split(hostStr, ":")
		port, _ := strconv.Atoi(hostSeg[1])
		host = Host{
			Ip:   hostSeg[0],
			Port: port,
		}
	}
	ok = true
	return
}
