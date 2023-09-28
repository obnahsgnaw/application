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
	"strconv"
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
	name     string
	ctx      context.Context
	cancel   context.CancelFunc
	cluster  *Cluster
	logger   *zap.Logger
	logCnf   *logger.Config
	debugger debug.Debugger
	servers  map[servertype.ServerType]map[endtype.EndType]map[string]Server
	event    *event.Manger
	register regCenter.Register
	errs     []error
	children []*Application
	regTtl   int64
	releases []func()
}

// New return a new application
func New(cluster *Cluster, name string, options ...Option) *Application {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	s := &Application{
		name:     name,
		ctx:      ctx,
		cancel:   cancel,
		cluster:  cluster,
		debugger: debug.New(dynamic.NewBool(func() bool { return false })),
		servers:  make(map[servertype.ServerType]map[endtype.EndType]map[string]Server),
		event:    event.New(),
		regTtl:   5,
	}
	if cluster == nil {
		s.addErr(applicationError("application cluster is required", nil))
	}
	if name == "" {
		s.addErr(applicationError("application name is required", nil))
	}
	s.With(options...)
	if s.logger == nil {
		s.logger, err = logger.New(utils.ToStr("App[", name, "]"), s.logCnf, s.debugger.Debug())
	}
	s.addErr(err)
	return s
}

func (app *Application) With(options ...Option) {
	for _, o := range options {
		o(app)
	}
}

// ID return application cluster id
func (app *Application) ID() string {
	return app.cluster.id
}

// Cluster return application cluster
func (app *Application) Cluster() *Cluster {
	return app.cluster
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
	app.debug(utils.ToStr("added ", server.EndType().String(), " ", server.Type().String(), " server:", server.Name()))
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
	if len(app.errs) > 0 {
		failedCb(app.errs[0])
		return
	}
	app.debug("starting...")
	app.displayConfig()

	app.initEvent(failedCb)

	for _, typeServers := range app.servers {
		for _, etServers := range typeServers {
			for _, s := range etServers {
				app.debug(utils.ToStr(s.EndType().String(), " ", s.Type().String(), " server[", s.Name(), "] start running "))
				s.Run(failedCb)
			}
		}
	}
	app.debug("servers init complete")

	if len(app.children) > 0 {
		for _, sub := range app.children {
			sub.With(Context(app.ctx))
			sub.debugger = app.debugger
			sub.logger = app.logger
			sub.logCnf = app.logCnf
			sub.register = app.register
			sub.regTtl = app.regTtl
			app.debug(utils.ToStr("sub application[", sub.name, "] start init"))
			sub.Run(failedCb)
		}
	}
	app.debug("sub applications init complete")
}

func (app *Application) displayConfig() {
	app.debug("cluster: " + app.cluster.String())
	if app.debugger.Debug() {
		app.debug("debug: true")
	} else {
		app.debug("debug: false")
	}
	if app.register != nil {
		app.debug("register: yes")
	} else {
		app.debug("register: no")
	}
	app.debug("register ttl: " + strconv.FormatInt(app.regTtl, 10))
}

func (app *Application) Wait() {
	signals.Listen(nil)
	app.logger.Info("start and serving...")
	signals.Wait()
	app.cancel()
	app.logger.Info("down")
}

// Release stop and release application
func (app *Application) Release() {
	for _, typeServers := range app.servers {
		for _, etServers := range typeServers {
			for _, s := range etServers {
				s.Release()
			}
		}
	}

	_ = app.logger.Sync()

	if app.register != nil {
		if etcd, ok := app.register.(*regCenter.EtcdRegister); ok {
			etcd.Release()
		}
	}

	if len(app.children) > 0 {
		for _, sub := range app.children {
			sub.Release()
		}
	}

	if len(app.releases) > 0 {
		for _, r := range app.releases {
			r()
		}
	}
	app.debug("application released")
}

func (app *Application) AddRelease(r func()) {
	if r != nil {
		app.releases = append(app.releases, r)
	}
}

// AddChild add sub application
func (app *Application) AddChild(apps ...*Application) {
	if len(apps) > 0 {
		for _, a := range apps {
			if a != nil && a.cluster != nil && a.cluster.id == app.cluster.id && a.name != app.name {
				app.children = append(app.children, a)
				app.debug("added sub application:" + a.Name())
			}
		}
	}

}

// DoRegister register
func (app *Application) DoRegister(regInfo *regCenter.RegInfo) error {
	if app.register == nil {
		app.logger.Warn("do register failed, no register to do")
		return nil
	}
	for k, v := range regInfo.Kvs() {
		if err := app.register.Register(app.ctx, k, v, regInfo.Ttl); err != nil {
			return applicationError("do register failed", err)
		}
		app.debug(utils.ToStr("registered:", k, "=>", v))
	}
	return nil
}

// DoUnregister unregister
func (app *Application) DoUnregister(regInfo *regCenter.RegInfo) error {
	if app.register == nil {
		app.logger.Warn("do unregister failed, no register to do")
		return nil
	}
	for k := range regInfo.Kvs() {
		if err := app.register.Unregister(app.ctx, k); err != nil {
			return applicationError("do unregister failed", err)
		}
		app.debug(utils.ToStr("unregistered:", k))
	}
	return nil
}

func (app *Application) initEvent(failedCb func(err error)) {
	if err := event.Init(app.event); err != nil {
		failedCb(err)
	}
	app.debug("event manager init complete")
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
