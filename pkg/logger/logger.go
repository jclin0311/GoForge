package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger(level string) error {
	var config zap.Config

	if level == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	var err error
	Log, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}

func Sync() error {
	return Log.Sync()
}
