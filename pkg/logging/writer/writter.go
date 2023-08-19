package writer

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

// NewFileWriter Get a writer
func NewFileWriter(file string, maxSize, maxBackUp, maxAge int, compress bool) io.Writer {
	return NewLumberjack(file, maxSize, maxBackUp, maxAge, compress)
}
func NewLumberjack(file string, maxSize, maxBackUp, maxAge int, compress bool) *lumberjack.Logger {
	if maxSize <= 0 {
		maxAge = 10
	}
	if maxBackUp <= 0 {
		maxBackUp = 5
	}
	if maxAge <= 0 {
		maxAge = 30
	}
	return &lumberjack.Logger{
		Filename:   file,
		MaxSize:    maxSize,
		MaxBackups: maxBackUp,
		MaxAge:     maxAge,
		Compress:   compress,
	}
}

// NewStdWriter std writer
func NewStdWriter() io.Writer {
	return os.Stdout
}

func NewNullWriter() io.Writer {
	return &NullWriter{}
}

type NullWriter struct {
}

func (NullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
