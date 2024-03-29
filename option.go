package application

import (
	"context"
	"github.com/obnahsgnaw/application/pkg/debug"
	"github.com/obnahsgnaw/application/pkg/dynamic"
	"github.com/obnahsgnaw/application/pkg/logging/logger"
	"github.com/obnahsgnaw/application/service/regCenter"
)

type Option func(s *Application)

func Context(ctx context.Context) Option {
	return func(s *Application) {
		if ctx != nil {
			s.ctx, s.cancel = context.WithCancel(ctx)
		}
	}
}

func CusCluster(c *Cluster) Option {
	return func(s *Application) {
		if c != nil {
			s.cluster = c
		}
	}
}

func Register(register regCenter.Register, ttl int64) Option {
	return func(s *Application) {
		if register != nil {
			s.register = register
			s.cusRegister = true
		}
		if ttl > 0 {
			s.regTtl = ttl
		}
	}
}

// Logger 需要在 CusCluster 和 debug之后
func Logger(config *logger.Config) Option {
	return func(s *Application) {
		if config != nil {
			s.initLogger(config)
		}
	}
}

func Debug(cb func() bool) Option {
	return func(s *Application) {
		if cb != nil {
			s.debugger = debug.New(dynamic.NewBool(cb))
		}
	}
}
