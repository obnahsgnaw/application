package logger

import (
	"testing"
)

func TestNewFileLogger(t *testing.T) {
	cnf := &Config{
		Dir:        "/Users/wangshanbo/Documents/Data/projects/application/out",
		MaxSize:    1,
		MaxBackup:  5,
		MaxAge:     5,
		Level:      "info",
		TraceLevel: "error",
	}
	l, err := New("", cnf, true)
	if err != nil {
		t.Errorf("got err=%s", err.Error())
		return
	}
	defer l.Sync()
	l.Debug("this is a debug message 1")
	l.Info("this is a info message 1")
	_ = cnf.SetLevel("debug")
	l.Debug("this is a debug message 2")
	l.Info("this is a info message 2")
}

func TestCopyCnfWithLevel(t *testing.T) {
	cnf := &Config{
		Dir:        "/Users/wangshanbo/Documents/Data/projects/application/out",
		MaxSize:    1,
		MaxBackup:  5,
		MaxAge:     5,
		Level:      "info",
		TraceLevel: "error",
	}
	cnf.InitLevel()
	cnf.InitTraceLevel()
	cnf.AddSubDir("application", "demo")

	cnf1 := CopyCnfWithLevel(cnf)
	cnf1.AddSubDir("rpc", "A")

	l, err := New("", cnf1, true)
	if err != nil {
		t.Errorf("got err=%s", err.Error())
		return
	}
	defer l.Sync()

	l.Debug("this is a debug message 1")
	l.Info("this is a info message 1")
	_ = cnf.SetLevel("debug")
	l.Debug("this is a debug message 2")
	l.Info("this is a info message 2")
}
