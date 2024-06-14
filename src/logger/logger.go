package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var l *zap.Logger

func Init(level, fileName string, maxSize, maxAge, maxBackups int, compress bool) {
	// 日志轮转
	writer := &lumberjack.Logger{
		// 日志名称
		Filename: fileName,
		// 日志大小限制，单位MB
		MaxSize: maxSize,
		// 历史日志文件保留天数
		MaxAge: maxAge,
		// 最大保留历史日志数量
		MaxBackups: maxBackups,
		// 本地时区
		LocalTime: true,
		// 历史日志文件压缩标识
		Compress: compress,
	}
	logCfg := zap.NewProductionEncoderConfig()
	logCfg.EncodeTime = zapcore.RFC3339TimeEncoder

	zapCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(logCfg),
		zapcore.AddSync(writer),
		logLevel(level),
	)
	l = zap.New(zapCore, zap.AddCaller())
}

func logLevel(level string) zap.AtomicLevel {
	atomicLevel := zap.NewAtomicLevel()
	switch strings.ToUpper(level) {
	case "DEBUG":
		atomicLevel.SetLevel(zapcore.DebugLevel)
	case "INFO":
		atomicLevel.SetLevel(zapcore.InfoLevel)
	case "WARN":
		atomicLevel.SetLevel(zapcore.WarnLevel)
	case "ERROR":
		atomicLevel.SetLevel(zapcore.ErrorLevel)
	case "DPANIC":
		atomicLevel.SetLevel(zapcore.DPanicLevel)
	case "PANIC":
		atomicLevel.SetLevel(zapcore.PanicLevel)
	case "FATAL":
		atomicLevel.SetLevel(zapcore.FatalLevel)
	default:
		atomicLevel.SetLevel(zapcore.DebugLevel)
	}
	return atomicLevel
}

func Sugar() *zap.SugaredLogger {
	return l.Sugar()
}

func Sync() {
	if l != nil {
		l.Sync()
	}
}
