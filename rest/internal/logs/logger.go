package logs

import (
	"fmt"
	"strings"

	"github.com/zhlii/wechat-box/rest/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func CreateLogger() {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "ts"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	loglevel, _ := parseLogLevel(config.Data.Common.LogLevel)

	outputPath := "stdout"
	errOutputPath := "stderr"
	if config.Data.Common.IsProd {
		outputPath = "rest.log"
		errOutputPath = "rest.log"
	}

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(loglevel),
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: true,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			outputPath,
		},
		ErrorOutputPaths: []string{
			errOutputPath,
		},
		InitialFields: map[string]interface{}{},
	}

	_logger, _ = config.Build()
}

func parseLogLevel(levelStr string) (zapcore.Level, error) {
	switch strings.ToLower(levelStr) {
	case "debug":
		return zap.DebugLevel, nil
	case "info":
		return zap.InfoLevel, nil
	case "warn", "warning":
		return zap.WarnLevel, nil
	case "error":
		return zap.ErrorLevel, nil
	case "fatal":
		return zap.FatalLevel, nil
	case "panic":
		return zap.PanicLevel, nil
	default:
		return zap.InfoLevel, fmt.Errorf("unknown log level: %s, use default level info", levelStr)
	}
}

var _logger *zap.Logger

// Debug logs an debug msg with fields
func Debug(msg string, fields ...zapcore.Field) {
	_logger.Debug(msg, fields...)
}

// Info logs an info msg with fields
func Info(msg string, fields ...zapcore.Field) {
	_logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zapcore.Field) {
	_logger.Warn(msg, fields...)
}

// Error logs an error msg with fields
func Error(msg string, fields ...zapcore.Field) {
	_logger.Error(msg, fields...)
}

// Fatal logs a fatal error msg with fields
func Fatal(msg string, fields ...zapcore.Field) {
	_logger.Fatal(msg, fields...)
}
