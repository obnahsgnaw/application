package debug

import "github.com/obnahsgnaw/application/pkg/dynamic"

type Debugger interface {
	SetDebug(dynamic.Bool)
	Debug() bool
}

type Debug struct {
	enable dynamic.Bool
}

func New(enable dynamic.Bool) *Debug {
	return &Debug{enable: enable}
}
func (d *Debug) SetDebug(enable dynamic.Bool) {
	d.enable = enable
}
func (d *Debug) Debug() bool {
	return d.enable.Val()
}
