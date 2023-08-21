package application

import (
	"context"
	"github.com/obnahsgnaw/application/endtype"
	"github.com/obnahsgnaw/application/pkg/debug"
	"github.com/obnahsgnaw/application/pkg/dynamic"
	"github.com/obnahsgnaw/application/pkg/logging/logger"
	"github.com/obnahsgnaw/application/pkg/signals"
	"github.com/obnahsgnaw/application/pkg/utils"
	"github.com/obnahsgnaw/application/servertype"
	"github.com/obnahsgnaw/application/service/event"
	"github.com/obnahsgnaw/application/service/regCenter"
	"go.uber.org/zap"
)

type Server interface {
	ID() string
	Name() string
	Type() servertype.ServerType
	EndType() endtype.EndType
	Run(func(error))
	Release()
}

func applicationError(msg string, err error) error {
	return utils.TitledError("application error", msg, err)
}

// application -->  server -->  end-type --> service

// Application identify a project
type Application struct {
	id       string
	name     string
	ctx      context.Context
	cancel   context.CancelFunc
	logger   *zap.Logger
	logCnf   *logger.Config
	errs     []error
	debugger debug.Debugger
	servers  map[servertype.ServerType]map[endtype.EndType]map[string]Server
	event    *event.Manger
	register regCenter.Register
	children []*Application
	regTtl   int64
}

// New return a new application
func New(id, name string, options ...Option) *Application {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	s := &Application{
		id:     id,
		name:   name,
		ctx:    ctx,
		cancel: cancel,
		debugger: debug.New(dynamic.NewBool(func() bool {
			return false
		})),
		event:   event.New(),
		servers: make(map[servertype.ServerType]map[endtype.EndType]map[string]Server),
		regTtl:  5,
	}
	s.With(options...)
	s.logger, err = logger.New(utils.ToStr("App[", name, "]"), s.logCnf, s.debugger.Debug())
	s.addErr(err)
	s.logger.Info("init")
	return s
}

func (app *Application) With(options ...Option) {
	for _, o := range options {
		o(app)
	}
}

// ID return application id
func (app *Application) ID() string {
	return app.id
}

// Name return app name
func (app *Application) Name() string {
	return app.name
}

// Context return application context
func (app *Application) Context() context.Context {
	return app.ctx
}

// Debugger return
func (app *Application) Debugger() debug.Debugger {
	return app.debugger
}

// Logger return the logger
func (app *Application) Logger() *zap.Logger {
	return app.logger
}

// LogConfig return
func (app *Application) LogConfig() *logger.Config {
	return app.logCnf
}

// Register return
func (app *Application) Register() regCenter.Register {
	return app.register
}

// RegTtl register ttl
func (app *Application) RegTtl() int64 {
	return app.regTtl
}

// Event return the event manager
func (app *Application) Event() *event.Manger {
	return app.event
}

// AddServer add a server
func (app *Application) AddServer(server Server) {
	if server == nil {
		return
	}
	if _, ok := app.servers[server.Type()]; !ok {
		app.servers[server.Type()] = make(map[endtype.EndType]map[string]Server)
	}
	if _, ok := app.servers[server.Type()][server.EndType()]; !ok {
		app.servers[server.Type()][server.EndType()] = make(map[string]Server)
	}
	app.servers[server.Type()][server.EndType()][server.ID()] = server
	app.debug(utils.ToStr("added ", server.Type().String(), " ", server.EndType().String(), " server:", server.Name()))
}

// GetTypeServers return servers
func (app *Application) GetTypeServers(typ servertype.ServerType) map[endtype.EndType]map[string]Server {
	if typ == "" {
		return make(map[endtype.EndType]map[string]Server)
	}
	if ss, ok := app.servers[typ]; ok {
		return ss
	}

	return make(map[endtype.EndType]map[string]Server)
}

// GetTypeServer return server
func (app *Application) GetTypeServer(typ servertype.ServerType, et endtype.EndType, id string) (Server, bool) {
	if id == "" {
		return nil, false
	}
	if ss, ok := app.servers[typ]; ok {
		if s, ok := ss[et]; ok {
			if s1, ok := s[id]; ok {
				return s1, true
			}
		}

	}

	return nil, false
}

// Run application
func (app *Application) Run(failedCb func(err error)) {
	if app.id == "" || app.name == "" {
		failedCb(applicationError("id or name invalid", nil))
		return
	}
	if len(app.errs) > 0 {
		failedCb(app.errs[0])
		return
	}
	app.debug("init event manager")
	app.initEvent(failedCb)

	for _, typeServers := range app.servers {
		for _, etServers := range typeServers {
			for _, s := range etServers {
				app.debug(utils.ToStr("start run ", s.Type().String(), " ", s.EndType().String(), " server:", s.Name()))
				s.Run(failedCb)
			}
		}
	}

	if len(app.children) > 0 {
		for _, sub := range app.children {
			sub.With(Context(app.ctx))
			if sub.debugger == nil {
				sub.debugger = app.debugger
			}
			if sub.logger == nil {
				sub.logger = app.logger
			}
			if sub.logCnf == nil {
				sub.logCnf = app.logCnf
			}
			if sub.register == nil {
				sub.register = app.register
			}
			if sub.regTtl <= 0 {
				sub.regTtl = app.regTtl
			}
			app.debug("run sub application:" + sub.name)
			go func(subApp *Application) {
				subApp.Run(failedCb)
			}(sub)
		}
	}
	// Listen signal
	signals.Listen(nil)
	app.logger.Info("start and serving...")
	signals.Wait()
	app.cancel()
	app.logger.Info(utils.ToStr("App[", app.name, "] down"))
}

// Release stop and release application
func (app *Application) Release() {
	app.debug("start release application")
	for _, typeServers := range app.servers {
		for _, etServers := range typeServers {
			for _, s := range etServers {
				app.debug(utils.ToStr("release ", s.Type().String(), " ", s.EndType().String(), " server:", s.Name()))
				s.Release()
			}
		}
	}
	app.debug("release application logger")
	_ = app.logger.Sync()
	if app.register != nil {
		if etcd, ok := app.register.(*regCenter.EtcdRegister); ok {
			app.debug("release etcd register")
			etcd.Release()
		}
	}
	if len(app.children) > 0 {
		for _, sub := range app.children {
			app.debug("release sub application:" + sub.name)
			sub.Release()
		}
	}
}

// AddChild add sub application
func (app *Application) AddChild(apps ...*Application) {
	if len(apps) > 0 {
		for _, a := range apps {
			if a != nil {
				app.children = append(app.children, a)
				app.debug("added sub application:" + a.Name())
			}
		}
	}

}

// DoRegister register
func (app *Application) DoRegister(regInfo *regCenter.RegInfo) error {
	if app.register == nil {
		return nil
	}
	for k, v := range regInfo.Kvs() {
		if err := app.register.Register(app.ctx, k, v, regInfo.Ttl); err != nil {
			return err
		}
	}
	return nil
}

// DoUnregister unregister
func (app *Application) DoUnregister(regInfo *regCenter.RegInfo) error {
	if app.register == nil {
		return nil
	}
	for k := range regInfo.Kvs() {
		if err := app.register.Unregister(app.ctx, k); err != nil {
			return err
		}
	}
	return nil
}

func (app *Application) initEvent(failedCb func(err error)) {
	if err := event.Init(app.event); err != nil {
		failedCb(err)
	}
}

func (app *Application) debug(msg string) {
	if app.debugger.Debug() {
		app.logger.Debug(msg)
	}
}

func (app *Application) addErr(err error) {
	if err != nil {
		app.errs = append(app.errs, err)
	}
}

func (app *Application) Errs() []error {
	return app.errs
}
