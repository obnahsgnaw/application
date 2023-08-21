package logger

import (
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
)

type Config struct {
	Dir        string `json:"dir" yaml:"dir" long:"log-dir" description:"Log file dir path." required:"true" default:""`
	MaxSize    int    `json:"max_size" yaml:"max_size" long:"log-maxSize" description:"Log file max size(M)." required:"true" default:"100"`
	MaxBackup  int    `json:"max_backup" yaml:"max_backup" long:"log-maxBackup" description:"Log file max backup." required:"true" default:"5"`
	MaxAge     int    `json:"max_age" yaml:"max_age" long:"log-maxAge" description:"Log file max age (day)." required:"true" default:"5"`
	Level      string `json:"level" yaml:"level"  long:"log-level" description:"Log level: debug,info, warn,error, ..." required:"true" default:"info"`
	TraceLevel string `json:"trace_level" yaml:"trace_level" long:"log-trace-level" description:"Log trace level: debug,info, warn,error, ..." required:"true" default:"error"`
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
		return "info"
	}

	return c.Level
}
func (c *Config) GetTraceLevel() string {
	if c.TraceLevel == "" {
		return "error"
	}

	return c.TraceLevel
}

func NewAccessWriter(cnf *Config, debug bool) (w io.Writer) {
	if cnf != nil && cnf.GetDir() != "" {
		w = writer.NewFileWriter(filepath.Join(cnf.GetDir(), "access.log"), cnf.GetMaxSize(), cnf.GetMaxBackup(), cnf.GetMaxAge(), true)
	}
	if w == nil && debug {
		w = writer.NewStdWriter()
	}
	return
}

func NewErrorWriter(cnf *Config, debug bool) (w io.Writer) {
	if cnf != nil && cnf.GetDir() != "" {
		w = writer.NewFileWriter(filepath.Join(cnf.Dir, "error.log"), cnf.GetMaxSize(), cnf.GetMaxBackup(), cnf.GetMaxAge(), true)
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
	if err = sinks.RegisterLumberjackSink(); err != nil {
		return nil, err
	}
	if cnf == nil || cnf.GetDir() == "" {
		err = loggerError("dir not set")
		return
	}

	f, err1 := os.Stat(cnf.GetDir())
	if err1 != nil {
		err = loggerError("dir invalid, err=" + err1.Error())
		return
	}
	if !f.IsDir() {
		err = loggerError("dir is not a directory")
		return
	}

	var level zapcore.Level
	if level, err = zapcore.ParseLevel(cnf.GetLevel()); err != nil {
		err = loggerError("level is invalid, err=" + err.Error())
		return
	}
	if name == "" {
		name = "log"
	}
	url := utils.ToStr("lumberjack://", filepath.Join(cnf.GetDir(), name+".log"), "?max_size=", strconv.Itoa(cnf.GetMaxSize()),
		"&max_age=", strconv.Itoa(cnf.GetMaxAge()), "&max_backup=", strconv.Itoa(cnf.GetMaxBackup()), "&compress=1")
	urlErr := utils.ToStr("lumberjack://", filepath.Join(cnf.GetDir(), "error.log"), "?max_size=", strconv.Itoa(cnf.GetMaxSize()),
		"&max_age=", strconv.Itoa(cnf.GetMaxAge()), "&max_backup=", strconv.Itoa(cnf.GetMaxBackup()), "&compress=1")
	l, err = logging.NewJsonLogger(name, level, []string{url}, []string{urlErr}, develop)
	if err != nil {
		err = loggerError("init failed, err=" + err.Error())
		return
	}
	l = l.WithOptions(zap.AddStacktrace(zap.ErrorLevel))
	return
}

func NewCliLogger(name, level string, develop bool) (l *zap.Logger, err error) {
	var levelZ zapcore.Level
	if levelZ, err = zapcore.ParseLevel(level); err != nil {
		err = loggerError("level is invalid, err=" + err.Error())
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
