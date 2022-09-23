package logger

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Error(message string)
	Warning(message string)
	Info(message string)
	Debug(message string)
	LogHTTPRequest(r *http.Request, statusCode int, duration time.Duration)
}

type logger struct {
	logFile string
	logger  *zap.Logger
}

func NewLogger(file string, level string) (Logger, error) {
	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(zapConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(zapConfig)

	var logLevel zapcore.Level
	if err := logLevel.Set(level); err != nil {
		return nil, err
	}

	logFile, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}

	writer := zapcore.AddSync(logFile)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, logLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), logLevel),
	)

	loggerZap := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &logger{
		logFile: file,
		logger:  loggerZap,
	}, nil
}

func (l *logger) Error(message string) {
	l.logger.Error(message)
}

func (l *logger) Warning(message string) {
	l.logger.Warn(message)
}

func (l *logger) Info(message string) {
	l.logger.Info(message)
}

func (l *logger) Debug(message string) {
	l.logger.Debug(message)
}

func (l *logger) LogHTTPRequest(r *http.Request, statusCode int, duration time.Duration) {
	l.logger.Debug(fmt.Sprintf("%s, %v, %v", r.RequestURI, statusCode, duration))
}
