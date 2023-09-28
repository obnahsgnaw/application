package application

import (
	"context"
	"github.com/obnahsgnaw/application/pkg/debug"
	"github.com/obnahsgnaw/application/pkg/dynamic"
	"github.com/obnahsgnaw/application/pkg/logging/logger"
	"github.com/obnahsgnaw/application/pkg/utils"
	"github.com/obnahsgnaw/application/service/regCenter"
	"time"
)

type Option func(s *Application)

func RegTtl(ttl int64) Option {
	return func(s *Application) {
		s.regTtl = ttl
	}
}

func Debugger(debugger debug.Debugger) Option {
	return func(s *Application) {
		if debugger != nil {
			s.debugger = debugger
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

func Logger(config *logger.Config) Option {
	return func(s *Application) {
		if config != nil {
			s.logCnf = config
			var err error
			s.logger, err = logger.New(utils.ToStr("App[", s.name, "]"), s.logCnf, s.debugger.Debug())
			s.addErr(err)
		}
	}
}

func Context(ctx context.Context) Option {
	return func(s *Application) {
		if ctx != nil {
			s.ctx, s.cancel = context.WithCancel(ctx)
		}
	}
}

func Register(register regCenter.Register) Option {
	return func(s *Application) {
		if register != nil {
			s.register = register
		}
	}
}

func EtcdRegister(endpoints []string, opeTimeout time.Duration) Option {
	return func(s *Application) {
		if opeTimeout == 0 {
			opeTimeout = 5 * time.Second
		}
		if len(endpoints) == 0 {
			s.addErr(applicationError("with EtcdRegister failed, etcd endpoint required", nil))
			return
		}
		etcdReg, err := regCenter.NewEtcdRegister(endpoints, opeTimeout)
		if err != nil {
			s.addErr(applicationError("with EtcdRegister new etcd register failed", err))
			return
		}
		s.With(Register(etcdReg))
	}
}
