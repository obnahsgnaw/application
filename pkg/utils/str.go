package utils

import "strings"

func ToStr(s ...string) string {
	if len(s) == 0 {
		return ""
	}
	if len(s) == 1 {
		return s[0]
	}
	return strings.Join(s, "")
}
