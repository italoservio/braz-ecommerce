package logger

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
)

type ContextKey string

const (
	CorrelationId ContextKey = "X-Correlation-ID"
)

type LoggerInterface interface {
	Info(msg string)
	Error(msg string)
	WithCtx(ctx context.Context) *Logger
}

type Logger struct {
	CorrelationId string
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Info(msg string) {
	message := fmt.Sprintf("%s - %s", l.CorrelationId, msg)

	slog.Info(message)
}

func (l *Logger) Error(msg string) {
	_, file, line, _ := runtime.Caller(1)
	caller := fmt.Sprintf("%s:%d", file, line)
	message := fmt.Sprintf("%s - %s at %s", l.CorrelationId, msg, caller)

	slog.Error(message)
}

func (l *Logger) WithCtx(ctx context.Context) *Logger {
	correlationId := ctx.Value(string(CorrelationId))

	if correlationId == nil {
		correlationId = "unknown"
	}

	l.CorrelationId = correlationId.(string)

	return l
}
