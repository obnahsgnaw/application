package utils

import (
	"fmt"
	"runtime/debug"
)

func RecoverHandler(desc string, handler func(err, stack string)) {
	if err := recover(); err != nil {
		s := string(debug.Stack())
		e := fmt.Sprintf("%v", err)
		if handler != nil {
			handler(e, s)
		}
		fmt.Println(desc, e, s)
	}
}
