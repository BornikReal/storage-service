package logger

import (
	"fmt"

	"go.uber.org/zap"
)

const (
	ErrorField   = "error"
	JobNameField = "job_name"
)

func Debug(msg string, fields ...zap.Field) {
	zap.L().Debug(msg, fields...)
}

func Debugf(msg string, args ...interface{}) {
	zap.L().Debug(fmt.Sprintf(msg, args...))
}

func Info(msg string, fields ...zap.Field) {
	zap.L().Info(msg, fields...)
}

func Infof(msg string, args ...interface{}) {
	zap.L().Info(fmt.Sprintf(msg, args...))
}

func Warn(msg string, fields ...zap.Field) {
	zap.L().Warn(msg, fields...)
}

func Warnf(msg string, args ...interface{}) {
	zap.L().Warn(fmt.Sprintf(msg, args...))
}

func Error(msg string, fields ...zap.Field) {
	zap.L().Error(msg, fields...)
}

func Errorf(msg string, args ...interface{}) {
	zap.L().Error(fmt.Sprintf(msg, args...))
}

func DPanic(msg string, fields ...zap.Field) {
	zap.L().DPanic(msg, fields...)
}

func DPanicf(msg string, args ...interface{}) {
	zap.L().DPanic(fmt.Sprintf(msg, args...))
}

func Panic(msg string, fields ...zap.Field) {
	zap.L().Panic(msg, fields...)
}

func Panicf(msg string, args ...interface{}) {
	zap.L().Panic(fmt.Sprintf(msg, args...))
}

func Fatal(msg string, fields ...zap.Field) {
	zap.L().Fatal(msg, fields...)
}

func Fatalf(msg string, args ...interface{}) {
	zap.L().Fatal(fmt.Sprintf(msg, args...))
}
