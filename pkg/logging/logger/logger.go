package logger

import (
	"errors"
	"github.com/obnahsgnaw/application/pkg/logging"
	"github.com/obnahsgnaw/application/pkg/logging/sinks"
	"github.com/obnahsgnaw/application/pkg/logging/writer"
	"github.com/obnahsgnaw/application/pkg/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	Dir                   string `json:"dir" yaml:"dir" long:"log-dir" description:"Log file dir path." required:"true" default:""`
	MaxSize               int    `json:"max_size" yaml:"max_size" long:"log-maxSize" description:"Log file max size(M)." required:"true" default:"100"`
	MaxBackup             int    `json:"max_backup" yaml:"max_backup" long:"log-maxBackup" description:"Log file max backup." required:"true" default:"5"`
	MaxAge                int    `json:"max_age" yaml:"max_age" long:"log-maxAge" description:"Log file max age (day)." required:"true" default:"5"`
	Level                 string `json:"level" yaml:"level"  long:"log-level" description:"Log level: debug,info, warn,error, ..." required:"true" default:"info"`
	TraceLevel            string `json:"trace_level" yaml:"trace_level" long:"log-trace-level" description:"Log trace level: debug,info, warn,error, ..." required:"true" default:"error"`
	level                 zap.AtomicLevel
	traceLevel            zap.AtomicLevel
	levelInitialized      bool
	traceLevelInitialized bool
	subDir                string
}

func (c *Config) GetDir() string {
	return c.Dir
}
func (c *Config) GetValidDir() (dir string, err error) {
	if c.Dir == "" {
		err = errors.New("dir empty")
		return
	}
	if dir, err = utils.ValidDir(c.GetDir()); err != nil {
		return
	}
	if c.subDir != "" {
		dir = filepath.Join(c.Dir, c.subDir)
		if err = os.MkdirAll(dir, 0777); err != nil {
			return
		}
	}

	return
}
func (c *Config) GetMaxSize() int {
	if c.MaxSize <= 0 {
		return 100
	}

	return c.MaxSize
}
func (c *Config) GetMaxBackup() int {
	if c.MaxBackup <= 0 {
		return 5
	}

	return c.MaxBackup
}
func (c *Config) GetMaxAge() int {
	if c.MaxAge <= 0 {
		return 30
	}

	return c.MaxAge
}
func (c *Config) GetLevelString() string {
	if c.Level == "" {
		return "info"
	}

	return c.Level
}
func (c *Config) InitLevel() error {
	if !c.levelInitialized {
		c.level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		if err := c.SetLevel(c.GetLevelString()); err != nil {
			return err
		}
		c.levelInitialized = true
	}
	return nil
}
func (c *Config) InitTraceLevel() error {
	if !c.traceLevelInitialized {
		c.traceLevel = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		if err := c.SetTraceLevel(c.GetTraceLevelString()); err != nil {
			return err
		}
		c.traceLevelInitialized = true
	}
	return nil
}
func (c *Config) SetLevel(level string) error {
	if l, err := zapcore.ParseLevel(level); err != nil {
		return err
	} else {
		c.level.SetLevel(l)
	}
	return nil
}
func (c *Config) GetLevel() zap.AtomicLevel {
	return c.level
}
func (c *Config) SetTraceLevel(level string) error {
	if l, err := zapcore.ParseLevel(level); err != nil {
		return err
	} else {
		c.traceLevel.SetLevel(l)
	}
	return nil
}
func (c *Config) GetTraceLevel() zap.AtomicLevel {
	return c.traceLevel
}
func (c *Config) GetTraceLevelString() string {
	if c.TraceLevel == "" {
		return "error"
	}

	return c.TraceLevel
}
func (c *Config) AddSubDir(dirname string) {
	if dirname != "" && !strings.HasSuffix(c.subDir, dirname) {
		c.subDir = filepath.Join(c.subDir, dirname)
	}
}

func NewAccessWriter(cnf *Config, debug bool) (w io.Writer, err error) {
	if cnf != nil && cnf.GetDir() != "" {
		var dir string
		if dir, err = cnf.GetValidDir(); err != nil {
			return
		}
		w = writer.NewFileWriter(filepath.Join(dir, "access.log"), cnf.GetMaxSize(), cnf.GetMaxBackup(), cnf.GetMaxAge(), true)
	}
	if w == nil && debug {
		w = writer.NewStdWriter()
	}
	return
}

func NewErrorWriter(cnf *Config, debug bool) (w io.Writer, err error) {
	if cnf != nil && cnf.GetDir() != "" {
		var dir string
		if dir, err = cnf.GetValidDir(); err != nil {
			return
		}
		w = writer.NewFileWriter(filepath.Join(dir, "error.log"), cnf.GetMaxSize(), cnf.GetMaxBackup(), cnf.GetMaxAge(), true)
	}
	if w == nil && debug {
		w = writer.NewStdWriter()
	}
	if w == nil {
		w = writer.NewNullWriter()
	}
	return
}

func loggerError(msg string) error {
	return utils.TitledError("logger error", msg, nil)
}

func NewFileLogger(name string, cnf *Config, develop bool) (l *zap.Logger, err error) {
	var dir string

	if err = sinks.RegisterLumberjackSink(); err != nil {
		return nil, err
	}
	if cnf == nil {
		err = loggerError("log config required")
		return
	}
	if dir, err = cnf.GetValidDir(); err != nil {
		err = loggerError("log dir invalid, err=" + err.Error())
		return
	}
	if err = cnf.InitLevel(); err != nil {
		err = loggerError("level is invalid, err=" + err.Error())
		return
	}
	if err = cnf.InitTraceLevel(); err != nil {
		err = loggerError("trace level is invalid, err=" + err.Error())
		return
	}
	url := utils.ToStr("lumberjack://", filepath.Join(dir, "log.log"), "?max_size=", strconv.Itoa(cnf.GetMaxSize()),
		"&max_age=", strconv.Itoa(cnf.GetMaxAge()), "&max_backup=", strconv.Itoa(cnf.GetMaxBackup()), "&compress=1")
	urlErr := utils.ToStr("lumberjack://", filepath.Join(dir, "error.log"), "?max_size=", strconv.Itoa(cnf.GetMaxSize()),
		"&max_age=", strconv.Itoa(cnf.GetMaxAge()), "&max_backup=", strconv.Itoa(cnf.GetMaxBackup()), "&compress=1")

	if l, err = logging.NewJsonLogger(name, cnf.GetLevel(), []string{url}, []string{urlErr}, develop); err != nil {
		err = loggerError("logger init failed, err=" + err.Error())
		return
	}
	l = l.WithOptions(zap.AddStacktrace(cnf.GetTraceLevel()))

	return
}

func NewCliLogger(name string, level zap.AtomicLevel, develop bool) (l *zap.Logger, err error) {
	if l, err = logging.NewCliLogger(name, level, develop); err == nil {
		l = l.WithOptions(zap.AddStacktrace(zap.FatalLevel))
	}
	return
}

func MergeLogger(l *zap.Logger, l1 ...*zap.Logger) *zap.Logger {
	return logging.NewMultiLogger(l, l1...)
}

func New(name string, cnf *Config, develop bool) (l *zap.Logger, err error) {
	if cnf != nil && cnf.GetDir() != "" {
		if l, err = NewFileLogger(name, cnf, develop); err != nil {
			return
		}
	}
	if l == nil || develop {
		level := zap.NewAtomicLevelAt(zapcore.DebugLevel)
		if cnf != nil {
			if err = cnf.InitLevel(); err != nil {
				return
			}
			level = cnf.GetLevel()
		}
		l1, _ := NewCliLogger(name, level, develop)
		if l == nil {
			l = l1
		} else {
			l = MergeLogger(l, l1)
		}
	}

	return l, nil
}
