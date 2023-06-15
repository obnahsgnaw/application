package utils

import "strings"

func ToStr(s ...string) string {
	return strings.Join(s, "")
}
