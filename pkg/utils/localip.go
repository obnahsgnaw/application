package utils

import (
	"net"
)

// GetLocalIp 获取本机ip
func GetLocalIp() (ip string, err error) {
	var conn net.Conn

	if conn, err = net.Dial("udp", "8.8.8.8:80"); err != nil {
		return "", err
	}
	defer func() { _ = conn.Close() }()

	return conn.LocalAddr().(*net.UDPAddr).IP.String(), nil
}
