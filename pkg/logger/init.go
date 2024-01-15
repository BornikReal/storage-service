package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() {
	cfg := zap.NewProductionEncoderConfig()

	f, err := os.OpenFile("storage_service.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defaultLogLevel := chooseLogLevel()
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(cfg), zapcore.AddSync(f), defaultLogLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)
}

func chooseLogLevel() zapcore.Level {
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel // Default log level
	}
}
