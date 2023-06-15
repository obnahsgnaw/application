package utils

import "time"

// InitTimezoneE8 设置时区为东八区
func InitTimezoneE8() {
	time.Local = time.FixedZone("CST", 8*3600) // 东8区
}
