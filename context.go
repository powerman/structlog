package structlog

import "context"

type contextKey int

const contextKeyLog contextKey = 0

// NewContext returns a new Context that carries value log.
func NewContext(ctx context.Context, log *Logger) context.Context {
	return context.WithValue(ctx, contextKeyLog, log)
}

// FromContext returns the Logger value stored in ctx or defaultLog or
// New() if defaultLog is nil.
func FromContext(ctx context.Context, defaultLog *Logger) *Logger {
	log, _ := ctx.Value(contextKeyLog).(*Logger)
	switch {
	case log != nil:
		return log
	case defaultLog != nil:
		return defaultLog
	default:
		return New()
	}
}
