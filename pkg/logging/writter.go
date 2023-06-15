package logging

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

// NewFileWriter Get a writer
func NewFileWriter(file string, maxSize, maxBackUp, maxAge int, compress bool) io.Writer {
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
