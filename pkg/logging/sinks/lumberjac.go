package sinks

import (
	"errors"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/url"
	"strconv"
)

type SizedFileLog struct {
	lumberjack.Logger
}

func (l *SizedFileLog) Sync() error {
	return nil
}

// NewLumberjackSink new a lumberjack sink writer
func NewLumberjackSink(filename string, maxSize, maxAge, maxBackup int, compress bool) *SizedFileLog {
	if maxSize <= 0 {
		maxAge = 10
	}
	if maxBackup <= 0 {
		maxBackup = 5
	}
	if maxAge <= 0 {
		maxAge = 30
	}
	return &SizedFileLog{lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackup,
		LocalTime:  false,
		Compress:   compress,
	}}
}

// RegisterLumberjackSink schema://filename?max_size=?&max_age=?&max_backup=?&compress=1
func RegisterLumberjackSink(schema string) error {
	if schema == "" {
		schema = "lumberjack"
	}
	return zap.RegisterSink(schema, func(url *url.URL) (zap.Sink, error) {
		var maxSize = 10
		var maxAge = 30
		var maxBackup = 5
		fileName := url.Host + url.Path
		if fileName == "" {
			return nil, errors.New("lumberjack filename is empty")
		}
		query := url.Query()
		if sizeStr := query.Get("max_size"); sizeStr != "" {
			size, err := strconv.Atoi(sizeStr)
			if err != nil {
				return nil, errors.New("max size err=" + err.Error())
			}
			if size > 0 {
				maxSize = size
			}
		}

		if ageStr := query.Get("max_age"); ageStr != "" {
			age, err := strconv.Atoi(ageStr)
			if err != nil {
				return nil, errors.New("max age err=" + err.Error())
			}
			if age > 0 {
				maxAge = age
			}
		}

		if backupStr := query.Get("max_backup"); backupStr != "" {
			backup, err := strconv.Atoi(backupStr)
			if err != nil {
				return nil, errors.New("max backup err=" + err.Error())
			}
			if backup > 0 {
				maxBackup = backup
			}
		}

		compress := query.Get("compress") == "1"

		return NewLumberjackSink(fileName, maxSize, maxAge, maxBackup, compress), nil
	})
}
