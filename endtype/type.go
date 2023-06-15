package endtype

const (
	Backend  EndType = "backend"
	Frontend EndType = "frontend"
)

// EndType end type
type EndType string

func (e EndType) String() string {
	return string(e)
}
