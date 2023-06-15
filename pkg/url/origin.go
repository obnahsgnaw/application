package url

import (
	"github.com/obnahsgnaw/application/pkg/utils"
	"strings"
)

const (
	HTTP  Protocol = "http"
	HTTPS Protocol = "https"
)

type Protocol string

func (p Protocol) String() string {
	return string(p)
}

type Origin struct {
	Protocol Protocol
	Host     Host
}

func (o Origin) String() string {
	if o.Protocol == "" {
		return ""
	}
	return utils.ToStr(o.Protocol.String(), "://", o.Host.String())
}

type Url struct {
	Origin Origin
	Path   string
}

func (o Url) String() string {
	if o.Path == "" {
		return o.Origin.String()
	}
	if o.Origin.String() == "" {
		return ""
	}
	return utils.ToStr(o.Origin.String(), "/", strings.TrimPrefix(o.Path, "/"))
}
