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
	url2 "net/url"
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
	Format                string `json:"format" yaml:"format" long:"log-format" description:"Log format: json 、 console" required:"false" default:"json"`
	level                 zap.AtomicLevel
	traceLevel            zap.AtomicLevel
	levelInitialized      bool
	traceLevelInitialized bool
	subDir                string
	fileName              string
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
		c.traceLevel = zap.NewAtomicLevelAt(zapcore.WarnLevel)
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
func (c *Config) ReplaceTraceLevel(l zap.AtomicLevel) {
	c.traceLevel = l
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
func (c *Config) AddSubDir(dirname ...string) {
	for _, n := range dirname {
		if n != "" && !strings.HasSuffix(c.subDir, n) {
			c.subDir = filepath.Join(c.subDir, n)
		}
	}
}
func (c *Config) GetFormat() string {
	if c.Format == "" || c.Format == "json" {
		return "json"
	} else {
		return "console"
	}
}
func (c *Config) SetFormat(f string) {
	if f == "" || f == "json" {
		c.Format = "json"
	} else {
		c.Format = "console"
	}
}
func (c *Config) GetFilename() string {
	if c.fileName == "" {
		return "log"
	} else {
		return c.fileName
	}
}
func (c *Config) SetFilename(f string) {
	c.fileName = f
}

func CopyCnfWithLevel(cnf *Config) *Config {
	c := CopyCnf(cnf)
	if c != nil {
		c.level = cnf.level
		c.traceLevel = cnf.traceLevel
	}

	return c
}

func CopyCnf(cnf *Config) *Config {
	if cnf == nil {
		return nil
	}
	c := *cnf

	return &c
}

func NewAccessWriter(cnf *Config) (w io.Writer, err error) {
	if cnf != nil && cnf.GetDir() != "" {
		var dir string
		if dir, err = cnf.GetValidDir(); err != nil {
			return
		}
		name := "access.log"
		if cnf.fileName != "" {
			name = cnf.fileName + "-" + name
		}
		w = writer.NewFileWriter(filepath.Join(dir, name), cnf.GetMaxSize(), cnf.GetMaxBackup(), cnf.GetMaxAge(), true)
	} else {
		err = errors.New("log config invalid")
	}
	return
}

func NewErrorWriter(cnf *Config) (w io.Writer, err error) {
	if cnf != nil && cnf.GetDir() != "" {
		var dir string
		if dir, err = cnf.GetValidDir(); err != nil {
			return
		}
		name := "error.log"
		if cnf.fileName != "" {
			name = cnf.fileName + "-" + name
		}
		w = writer.NewFileWriter(filepath.Join(dir, name), cnf.GetMaxSize(), cnf.GetMaxBackup(), cnf.GetMaxAge(), true)
	} else {
		err = errors.New("log config invalid")
	}
	return
}

func NewDefAccessWriter(cnf *Config, debug func() bool) (io.Writer, error) {
	var wts []io.Writer
	wts = append(wts, writer.NewDynamicStdWriter(debug, os.Stdout))
	if cnf != nil && cnf.GetDir() != "" {
		if w, err := NewAccessWriter(cnf); err == nil {
			wts = append(wts, w)
		} else {
			return nil, err
		}
	}

	return writer.NewMultiWriter(wts...), nil
}

func NewDefErrorWriter(cnf *Config, debug func() bool) (io.Writer, error) {
	var wts []io.Writer
	wts = append(wts, writer.NewDynamicStdWriter(debug, os.Stderr))
	if cnf != nil && cnf.GetDir() != "" {
		if w, err := NewErrorWriter(cnf); err == nil {
			wts = append(wts, w)
		} else {
			return nil, err
		}

	}

	return writer.NewMultiWriter(wts...), nil
}

func loggerError(msg string) error {
	return utils.TitledError("logger error", msg, nil)
}

func NewFileLogger(name string, cnf *Config, develop bool) (l *zap.Logger, err error) {
	var dir string
	defer func() {
		if l == nil {
			l, _ = NewCliLogger(name, zap.NewAtomicLevelAt(zapcore.DebugLevel), true)
		}
	}()

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
	qr := url2.Values{}
	qr.Add("path", filepath.Join(dir, cnf.GetFilename()+".log"))
	qr.Add("max_size", strconv.Itoa(cnf.GetMaxSize()))
	qr.Add("max_age", strconv.Itoa(cnf.GetMaxAge()))
	qr.Add("max_backup", strconv.Itoa(cnf.GetMaxBackup()))
	qr.Add("compress", "1")
	url := utils.ToStr("lumberjack://file?", qr.Encode())

	eqr := url2.Values{}
	eqr.Add("path", filepath.Join(dir, cnf.GetFilename()+"-error.log"))
	eqr.Add("max_size", strconv.Itoa(cnf.GetMaxSize()))
	eqr.Add("max_age", strconv.Itoa(cnf.GetMaxAge()))
	eqr.Add("max_backup", strconv.Itoa(cnf.GetMaxBackup()))
	eqr.Add("compress", "1")
	urlErr := utils.ToStr("lumberjack://file?", eqr.Encode())

	if cnf.GetFormat() == "json" {
		if l, err = logging.NewJsonLogger(name, cnf.GetLevel(), []string{url}, []string{urlErr}, develop); err != nil {
			err = loggerError("logger init failed, err=" + err.Error())
			return
		}
	} else {
		if l, err = logging.NewConsoleLogger(name, cnf.GetLevel(), []string{url}, []string{urlErr}, develop); err != nil {
			err = loggerError("logger init failed, err=" + err.Error())
			return
		}
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
