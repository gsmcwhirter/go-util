package logging

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"       //nolint:depguard
	"github.com/go-kit/kit/log/level" //nolint:depguard

	"github.com/gsmcwhirter/go-util/v4/errors"
	"github.com/gsmcwhirter/go-util/v4/request"
)

// DefaultTimestampUTC is a passthrough to the go-kit object of the same name
var DefaultTimestampUTC = log.DefaultTimestampUTC

// DefaultCaller is an alternative to the go-kit object of the same name to account for wrapping
var DefaultCaller = log.Caller(6)

// Logger is the extended logging interface for the corvee applications
type Logger interface {
	Log(keyvals ...interface{}) error
	Message(string, ...interface{})
	Err(string, error, ...interface{})
	Printf(string, ...interface{})
}

type logger struct {
	log.Logger
}

func (l *logger) Log(args ...interface{}) error {
	return l.Logger.Log(args...)
}

func (l *logger) Printf(f string, args ...interface{}) {
	m := fmt.Sprintf(f, args...)
	if err := l.Logger.Log("message", m); err != nil {
		panic(errors.WithDetails(err, "message", m))
	}
}

func (l *logger) Message(msg string, args ...interface{}) {
	args = append([]interface{}{"message", msg}, args...)
	if err := l.Logger.Log(args...); err != nil {
		panic(errors.WithDetails(err, args...))
	}
}

func (l *logger) Err(msg string, err error, args ...interface{}) {
	if e, ok := err.(errors.Error); ok {
		args = append([]interface{}{"message", msg, "error", e.Msg()}, args...)
		args = append(args, e.Data()...)
		if logErr := l.Logger.Log(args...); logErr != nil {
			panic(errors.WithDetails(logErr, args...))
		}
	}

	args = append([]interface{}{"message", msg, "error", err}, args...)
	if logErr := l.Logger.Log(args...); logErr != nil {
		panic(errors.WithDetails(logErr, args...))
	}
}

// NewFromKitLogger wraps a go-kit logger in our custom extension
func NewFromKitLogger(l log.Logger) Logger {
	return &logger{l}
}

// NewJSONLogger creates a new logger that writes json to stdout
func NewJSONLogger() Logger {
	return NewFromKitLogger(log.NewJSONLogger(log.NewSyncWriter(os.Stdout)))
}

// NewLogfmtLogger creates a new logger that writes logfmt to stdout
func NewLogfmtLogger() Logger {
	return NewFromKitLogger(log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout)))
}

// WithLevel wraps a logger to filter out logs lower than the designated level
func WithLevel(l log.Logger, levelStr string) Logger {
	switch levelStr {
	case "debug":
		return NewFromKitLogger(level.NewFilter(l, level.AllowDebug()))
	case "info":
		return NewFromKitLogger(level.NewFilter(l, level.AllowInfo()))
	case "warn":
		return NewFromKitLogger(level.NewFilter(l, level.AllowWarn()))
	case "error":
		return NewFromKitLogger(level.NewFilter(l, level.AllowError()))
	default:
		return NewFromKitLogger(level.NewFilter(l, level.AllowAll()))
	}
}

// With wraps a logger so that every emitted line contains the provided key/val pairs
func With(l log.Logger, keyvals ...interface{}) Logger {
	return &logger{log.With(l, keyvals...)}
}

// WithContext wraps a logger to include the request_id from a context in log messages
func WithContext(ctx context.Context, logger log.Logger) Logger {
	rid, ok := request.GetRequestID(ctx)
	if !ok {
		return With(logger, "request_id", "unknown")
	}

	return With(logger, "request_id", rid)
}

// WithRequest wraps a logger with fields from a http.Request
func WithRequest(req *http.Request, l log.Logger) Logger {
	logger := WithContext(req.Context(), l)
	logger = With(logger,
		"request_host", req.Host,
		"request_method", req.Method,
		"request_uri", req.RequestURI,
	)
	return logger
}
