package debug

type Debugger interface {
	SetDebug(bool)
	Debug() bool
}

type Debug struct {
	enable bool
}

func New(enable bool) *Debug {
	return &Debug{enable: enable}
}
func (d *Debug) SetDebug(enable bool) {
	d.enable = enable
}
func (d *Debug) Debug() bool {
	return d.enable
}
