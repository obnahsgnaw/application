package traefik

import "strings"

func EtcdKey(k ...string) string {
	return strings.Join(k, "/")
}

func BoolVal(flag bool) string {
	if flag {
		return "true"
	}
	return "false"
}
