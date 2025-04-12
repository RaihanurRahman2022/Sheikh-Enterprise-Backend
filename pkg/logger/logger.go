package logger

import (
	"os"

	"Sheikh-Enterprise-Backend/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Initialize sets up the logger with the given configuration
func Initialize(cfg *config.LoggerConfig) error {
	// Create encoder configuration
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = ""

	// Create encoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// Create core with file and console output
	var core zapcore.Core
	if cfg.File != "" {
		// Open log file
		logFile, err := os.OpenFile(cfg.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		// Create file writer
		fileWriter := zapcore.AddSync(logFile)

		// Create console writer
		consoleWriter := zapcore.AddSync(os.Stdout)

		// Set log level
		level := getLogLevel(cfg.Level)

		// Create core with multiple writers
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, fileWriter, level),
			zapcore.NewCore(encoder, consoleWriter, level),
		)
	} else {
		// If no file specified, log only to console
		core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), getLogLevel(cfg.Level))
	}

	// Create logger
	log = zap.New(core)
	defer log.Sync()

	return nil
}

// getLogLevel converts string level to zapcore.Level
func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

// With creates a child logger with additional fields
func With(fields ...zap.Field) *zap.Logger {
	return log.With(fields...)
}
