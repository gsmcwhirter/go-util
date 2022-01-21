package logging

import (
	"context"
	"fmt"
	"io"
	stdLog "log" //nolint:depguard // this is the package that wraps the stdlib
	"net/http"
	"os"

	"github.com/go-kit/kit/log"       //nolint:depguard // uses this internally to do the logging
	"github.com/go-kit/kit/log/level" //nolint:depguard // uses this internally to do the logging

	"github.com/gsmcwhirter/go-util/v9/errors"
	"github.com/gsmcwhirter/go-util/v9/request"
)

// DefaultTimestampUTC is a passthrough to the go-kit object of the same name
var DefaultTimestampUTC = log.DefaultTimestampUTC

// DefaultCaller is an alternative to the go-kit object of the same name to account for wrapping
var DefaultCaller = log.Caller(4)

type BaseLogger interface {
	Log(keyvals ...interface{}) error
}

// Logger is the extended logging interface for the corvee applications
type Logger interface {
	Log(keyvals ...interface{}) error
	Message(string, ...interface{})
	Err(string, error, ...interface{})
	Printf(string, ...interface{})
}

type logger struct {
	base BaseLogger
}

func (l *logger) Log(args ...interface{}) error {
	return l.base.Log(args...)
}

func (l *logger) Printf(f string, args ...interface{}) {
	m := fmt.Sprintf(f, args...)
	if err := l.base.Log("message", m); err != nil {
		panic(errors.WithDetails(err, "message", m))
	}
}

func (l *logger) Message(msg string, args ...interface{}) {
	args = append([]interface{}{"message", msg}, args...)
	if err := l.base.Log(args...); err != nil {
		panic(errors.WithDetails(err, args...))
	}
}

func (l *logger) Err(msg string, err error, args ...interface{}) {
	if e, ok := err.(errors.Error); ok {
		args = append([]interface{}{"message", msg, "error", e.Msg()}, args...)
		args = append(args, e.Data()...)
		if logErr := l.base.Log(args...); logErr != nil {
			panic(errors.WithDetails(logErr, args...))
		}

		return
	}

	args = append([]interface{}{"message", msg, "error", err}, args...)
	if logErr := l.base.Log(args...); logErr != nil {
		panic(errors.WithDetails(logErr, args...))
	}
}

// NewFrom wraps a BaseLogger (e.g., go-kit) in our custom extension
func NewFrom(l BaseLogger) Logger {
	if l2, ok := l.(*logger); ok {
		return &logger{base: l2.base}
	}

	return &logger{base: l}
}

// NewFromKitLogger wraps a BaseLogger (e.g., go-kit) in our custom extension
var NewFromKitLogger = NewFrom

func BaseFrom(l Logger) BaseLogger {
	var base BaseLogger = l
	if lbase, ok := l.(*logger); ok {
		base = lbase.base
	}

	return base
}

// NewJSONLogger creates a new logger that writes json to stdout
func NewJSONLogger() Logger {
	return NewFrom(log.NewJSONLogger(log.NewSyncWriter(os.Stdout)))
}

// NewLogfmtLogger creates a new logger that writes logfmt to stdout
func NewLogfmtLogger() Logger {
	return NewFrom(log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout)))
}

// WithLevel wraps a logger to filter out logs lower than the designated level
func WithLevel(l Logger, levelStr string) Logger {
	base := BaseFrom(l)

	switch levelStr {
	case "debug":
		return NewFrom(level.NewFilter(base, level.AllowDebug()))
	case "info":
		return NewFrom(level.NewFilter(base, level.AllowInfo()))
	case "warn":
		return NewFrom(level.NewFilter(base, level.AllowWarn()))
	case "error":
		return NewFrom(level.NewFilter(base, level.AllowError()))
	default:
		return NewFrom(level.NewFilter(base, level.AllowAll()))
	}
}

// With wraps a logger so that every emitted line contains the provided key/val pairs
func With(l Logger, keyvals ...interface{}) Logger {
	return NewFrom(log.With(BaseFrom(l), keyvals...))
}

// WithContext wraps a logger to include the request_id from a context in log messages
func WithContext(ctx context.Context, logger Logger, keyvals ...interface{}) Logger {
	if rid, ok := request.GetRequestID(ctx); ok {
		keyvals = append(keyvals, "request_id", rid)
	} else {
		keyvals = append(keyvals, "request_id", "unknown")
	}

	return With(logger, keyvals...)
}

// WithRequest wraps a logger with fields from a http.Request
func WithRequest(req *http.Request, l Logger, keyvals ...interface{}) Logger {
	keyvals = append(keyvals,
		"request_host", req.Host,
		"request_method", req.Method,
		"request_uri", req.RequestURI,
	)

	return WithContext(req.Context(), l, keyvals...)
}

// PatchStdLib sets up the stdlib global logger to run through the provided one instead
func PatchStdLib(l Logger) {
	var w io.Writer = writer{l}

	stdLog.SetOutput(w)
	stdLog.SetFlags(0)
}
