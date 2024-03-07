package logger

import (
	"github.com/obnahsgnaw/application/pkg/logging/writer"
	"github.com/obnahsgnaw/application/pkg/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

// 需求： 1. 1个logger，
//       2. 基本日志 和 错误日志分别写入不同的文件（设置了日志目录），同时也直接可输出控制台（开启了debug）
//       3. 文件采用json格式， 控制台采用console格式
//       4. 如果都没有，返回一个空的logger

func NewLogger(cnf *Config, debug bool) (l *zap.Logger, err error) {
	var dir string
	var ww zapcore.WriteSyncer
	var cw = zapcore.Lock(os.Stdout)
	var cores []zapcore.Core

	if cnf == nil {
		err = loggerError("log config required")
		cnf = &Config{}
	}
	if cnf.Dir != "" {
		if dir, err = utils.ValidDir(cnf.Dir); err != nil {
			err = loggerError("log dir invalid, err=" + err.Error())
			dir = ""
		}
	}
	if err = cnf.InitLevel(); err != nil {
		err = loggerError("level is invalid, err=" + err.Error())
		cnf.level.SetLevel(zap.DebugLevel)
	}
	if err = cnf.InitTraceLevel(); err != nil {
		err = loggerError("trace level is invalid, err=" + err.Error())
		cnf.traceLevel.SetLevel(zap.ErrorLevel)
	}

	if dir != "" {
		jsonEncodeConfig := zap.NewProductionEncoderConfig()
		jsonEncodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		jsonEncoder := zapcore.NewJSONEncoder(jsonEncodeConfig)
		ww = zapcore.AddSync(writer.NewFileWriter(filepath.Join(dir, cnf.GetFilename()+".log"), cnf.GetMaxSize(), cnf.GetMaxBackup(), cnf.GetMaxAge(), true))
		cores = append(cores, zapcore.NewCore(jsonEncoder, ww, cnf.GetLevel()))
	}
	if debug {
		consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
		consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)
		cores = append(cores, zapcore.NewCore(consoleEncoder, cw, cnf.GetLevel()))
	}
	if len(cores) == 0 {
		cores = append(cores, zapcore.NewNopCore())
	}

	core := zapcore.NewTee(cores...)

	l = zap.New(core, zap.AddStacktrace(cnf.GetTraceLevel()))
	return
}
