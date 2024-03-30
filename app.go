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

// application -->  server -->  end-type --> service

// Application identify a project
type Application struct {
	ctx         context.Context
	cancel      context.CancelFunc
	name        string
	cluster     *Cluster
	logger      *zap.Logger
	logCnf      *logger.Config
	logCus      bool
	debugger    debug.Debugger
	event       *event.Manger
	register    regCenter.Register
	cusRegister bool
	servers     map[servertype.ServerType]map[endtype.EndType]map[string]Server
	errs        []error
	releases    []func()
	callbacks   []func()
	children    []*Application
	regTtl      int64
}

// New return a new application
func New(name string, options ...Option) *Application {
	ctx, cancel := context.WithCancel(context.Background())
	if name == "" {
		name = "default"
	}
	s := &Application{
		name:     name,
		ctx:      ctx,
		cancel:   cancel,
		cluster:  NewCluster("dev", "Dev"),
		debugger: debug.New(dynamic.NewBool(func() bool { return true })),
		event:    event.New(),
		register: regCenter.NewNone(),
		servers:  make(map[servertype.ServerType]map[endtype.EndType]map[string]Server),
		regTtl:   5,
		logCnf:   &logger.Config{},
	}
	s.With(options...)
	if s.logger == nil {
		_ = s.initLogger()
	}
	return s
}

func (app *Application) With(options ...Option) {
	for _, o := range options {
		if o != nil {
			o(app)
		}
	}
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
}

// DelServer del added server
func (app *Application) DelServer(server Server) {
	if server == nil {
		return
	}
	if _, ok := app.servers[server.Type()]; ok {
		if _, ok = app.servers[server.Type()][server.EndType()]; !ok {
			delete(app.servers[server.Type()][server.EndType()], server.ID())
		}
	}
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
	defer app.handleCallback()
	if app.logCus {
		if err := app.initLogger(); err != nil {
			failedCb(err)
		}
	}
	app.logger.Info(app.prefixedMsg("init starting..."))
	app.displayConfig()
	if !app.initEvent(func(err error) {
		failedCb(app.error("init event failed", err))
	}) {
		return
	}
	app.logger.Info(app.prefixedMsg("event manager initialized"))
	hadServer := false
	for _, typeServers := range app.servers {
		for _, etServers := range typeServers {
			for _, s := range etServers {
				app.logger.Debug(app.prefixedMsg(s.EndType().String(), " ", s.Type().String(), " server[", s.Name(), "] init starting..."))
				s.Run(failedCb)
				app.logger.Debug(app.prefixedMsg(s.EndType().String(), " ", s.Type().String(), " server[", s.Name(), "] initialized"))
				hadServer = true
			}
		}
	}
	if !hadServer {
		app.logger.Warn(app.prefixedMsg("services initialized, but no services registered"))
	} else {
		app.logger.Info(app.prefixedMsg("services initialized"))
	}
	if len(app.children) > 0 {
		for _, sub := range app.children {
			sub.With(Context(app.ctx))
			sub.debugger = app.debugger
			sub.logger = app.logger.Named(sub.name)
			sub.logCnf = app.logCnf
			if sub.register == nil {
				sub.register = app.register
				sub.regTtl = app.regTtl
			}
			app.logger.Debug(app.prefixedMsg("sub-application[", sub.name, "] init starting..."))
			sub.Run(failedCb)
			app.logger.Debug(app.prefixedMsg("sub-application[", sub.name, "] initialized"))
		}
		app.logger.Info(app.prefixedMsg("sub-applications initialized"))
	} else {
		app.logger.Info(app.prefixedMsg("sub-applications initialized, but no sub-applications registered"))
	}
	if !app.cusRegister {
		app.logger.Warn(app.prefixedMsg("register, no server-register registered"))
	}
	app.logger.Info(app.prefixedMsg("initialized"))
}

func (app *Application) Wait() {
	signals.Listen(nil)
	app.logger.Info(app.prefixedMsg("started and serving..."))
	signals.Wait()
	app.cancel()
	app.logger.Info(app.prefixedMsg("down"))
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

	if app.logger != nil {
		app.logger.Info(app.prefixedMsg("released"))
		_ = app.logger.Sync()
	}
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
			}
		}
	}
}

// DoRegister register
func (app *Application) DoRegister(regInfo *regCenter.RegInfo, cb func(string)) error {
	for k, v := range regInfo.Kvs() {
		if err := app.register.Register(app.ctx, k, v, regInfo.Ttl); err != nil {
			return app.error("register failed", err)
		}
		if cb != nil {
			cb(utils.ToStr("registered:", k, "=>", v))
		}
	}
	return nil
}

// DoUnregister unregister
func (app *Application) DoUnregister(regInfo *regCenter.RegInfo, cb func(string)) error {
	for k := range regInfo.Kvs() {
		if err := app.register.Unregister(app.ctx, k); err != nil {
			return app.error("unregister failed", err)
		}
		if cb != nil {
			cb(utils.ToStr("unregistered:", k))
		}
	}
	return nil
}

func (app *Application) RegisterCallback(cb func()) {
	if cb != nil {
		app.callbacks = append(app.callbacks, cb)
	}
}

func (app *Application) initEvent(failedCb func(err error)) bool {
	if err := event.Init(app.event); err != nil {
		failedCb(err)
		return false
	}
	return true
}

func (app *Application) handleCallback() {
	for _, cb := range app.callbacks {
		cb()
	}
	if len(app.children) > 0 {
		for _, sub := range app.children {
			sub.handleCallback()
		}
	}
}

func (app *Application) displayConfig() {
	debugDesc := "debug=false"
	if app.debugger.Debug() {
		debugDesc = "debug=true"
	}
	regDesc := "register-set=false"
	if app.register != nil {
		regDesc = "register-set=true"
	}
	app.logger.Debug(app.prefixedMsg(utils.ToStr(
		"config: cluster=", app.cluster.String(),
		",", debugDesc,
		",", regDesc,
		",register-ttl=", strconv.FormatInt(app.regTtl, 10),
	)))
}

func (app *Application) prefixedMsg(msg ...string) string {
	return utils.ToStr(msg...)
}

func (app *Application) initLogger() (err error) {
	app.logCnf.SetFilename(app.cluster.id)
	app.logger, err = logger.NewLogger(app.logCnf, app.debugger.Debug())
	if err != nil {
		return
	}
	app.logger = app.logger.Named(app.name)
	return
}

func (app *Application) error(msg string, err error) error {
	return utils.TitledError("application["+app.name+"] error", msg, err)
}
