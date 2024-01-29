package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func createLogger() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "ts"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: true,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"log_file.log",
		},
		ErrorOutputPaths: []string{
			"log_file.log",
		},
		InitialFields: map[string]interface{}{},
	}

	l, _ := config.Build()
	return l
}

var L = createLogger()

// Debug logs an debug msg with fields
func Debug(msg string, fields ...zapcore.Field) {
	L.Debug(msg, fields...)
}

// Info logs an info msg with fields
func Info(msg string, fields ...zapcore.Field) {
	L.Info(msg, fields...)
}

// Error logs an error msg with fields
func Error(msg string, fields ...zapcore.Field) {
	L.Error(msg, fields...)
}

// Fatal logs a fatal error msg with fields
func Fatal(msg string, fields ...zapcore.Field) {
	L.Fatal(msg, fields...)
}
