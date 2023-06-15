package logger

import (
	"errors"
	"github.com/obnahsgnaw/application/pkg/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path/filepath"
)

type Config struct {
	Dir        string `ini:"dir" long:"log-dir" description:"Log file dir path." required:"true" default:""`
	MaxSize    int    `ini:"max-size" long:"log-maxSize" description:"Log file max size(M)." required:"true" default:"100"`
	MaxBackup  int    `ini:"max-backup" long:"log-maxBackup" description:"Log file max backup." required:"true" default:"5"`
	MaxAge     int    `ini:"max-age" long:"log-maxAge" description:"Log file max age (day)." required:"true" default:"5"`
	Level      string `ini:"level" long:"log-level" description:"Log level: debug,info, warn,error, ..." required:"true" default:"Info"`
	TraceLevel string `ini:"trace-level" long:"log-level" description:"Log level: debug,info, warn,error, ..." required:"true" default:"Error"`
}

func (c *Config) GetDir() string {
	return c.Dir
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
func (c *Config) GetLevel() string {
	if c.Level == "" {
		return "Info"
	}

	return c.Level
}
func (c *Config) GetTraceLevel() string {
	if c.TraceLevel == "" {
		return "Error"
	}

	return c.TraceLevel
}

func NewAccessWriter(cnf *Config, debug bool) (w io.Writer) {
	if cnf != nil && cnf.GetDir() != "" {
		w = logging.NewFileWriter(filepath.Join(cnf.GetDir(), "access.log"), cnf.GetMaxSize(), cnf.GetMaxBackup(), cnf.GetMaxAge(), true)
	}
	if w == nil && debug {
		w = logging.NewStdWriter()
	}
	return
}

func NewErrorWriter(cnf *Config, debug bool) (w io.Writer) {
	if cnf != nil && cnf.GetDir() != "" {
		w = logging.NewFileWriter(filepath.Join(cnf.Dir, "error.log"), cnf.GetMaxSize(), cnf.GetMaxBackup(), cnf.GetMaxAge(), true)
	}
	if w == nil && debug {
		w = logging.NewStdWriter()
	}
	return
}

func NewFileLogger(name string, cnf *Config, develop bool) (l *zap.Logger, err error) {
	if cnf == nil || cnf.GetDir() == "" {
		err = errors.New("file log dir required")
		return
	}

	f, err1 := os.Stat(cnf.GetDir())
	if err1 != nil {
		err = errors.New("log dir err, err=" + err1.Error())
		return
	}
	if !f.IsDir() {
		err = errors.New("log dir is not a directory")
		return
	}

	var level zapcore.Level
	if level, err = zapcore.ParseLevel(cnf.GetLevel()); err != nil {
		err = errors.New("logger level is invalid, err=" + err.Error())
		return
	}
	if name == "" {
		name = "log"
	}
	l, err = logging.NewJsonLogger(name, level, []string{filepath.Join(cnf.GetDir(), name+".log")}, []string{filepath.Join(cnf.GetDir(), "error.log")}, develop)
	if err != nil {
		err = errors.New("logger init failed, err=" + err.Error())
		return
	}
	l = l.WithOptions(zap.AddStacktrace(zap.ErrorLevel))
	return
}

func NewCliLogger(name, level string, develop bool) (l *zap.Logger, err error) {
	var levelZ zapcore.Level
	if levelZ, err = zapcore.ParseLevel(level); err != nil {
		err = errors.New("logger level is invalid, err=" + err.Error())
		return
	}
	if name == "" {
		name = "log"
	}
	l, err = logging.NewCliLogger(name, levelZ, develop)
	if err == nil {
		l = l.WithOptions(zap.AddStacktrace(zap.ErrorLevel))
	}
	return
}

func MergeLogger(l *zap.Logger, l1 ...*zap.Logger) *zap.Logger {
	return logging.NewMultiLogger(l, l1...)
}

func New(name string, cnf *Config, develop bool) (*zap.Logger, error) {
	l, err := NewFileLogger(name, cnf, develop)
	if err != nil || develop {
		err = nil
		l1, _ := NewCliLogger(name, "debug", develop)
		if l == nil {
			return l1, err
		}
		l = MergeLogger(l, l1)
	}

	return l, err
}
