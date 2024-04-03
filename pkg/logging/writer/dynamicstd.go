package writer

import (
	"go.uber.org/zap/zapcore"
)

type DynamicStdWriter struct {
	enable func() bool
	w      zapcore.WriteSyncer
}

func NewDynamicStdWriter(enable func() bool, w zapcore.WriteSyncer) *DynamicStdWriter {
	return &DynamicStdWriter{
		enable: enable,
		w:      w,
	}
}

func (s *DynamicStdWriter) Write(p []byte) (n int, err error) {
	if s.enable() {
		return s.w.Write(p)
	}
	return len(p), nil
}

func (s *DynamicStdWriter) Sync() error {
	return s.w.Sync()
}
