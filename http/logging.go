package http

import (
	"github.com/gsmcwhirter/go-util/v12/logging"
	"github.com/gsmcwhirter/go-util/v12/logging/level"
)

type HTTPLogger struct { //nolint:revive // ok with stutter
	Logger logging.Logger
}

func (h *HTTPLogger) Error(msg string, keysAndValues ...interface{}) {
	level.Error(h.Logger).Message(msg, keysAndValues...)
}

func (h *HTTPLogger) Info(msg string, keysAndValues ...interface{}) {
	level.Info(h.Logger).Message(msg, keysAndValues...)
}

func (h *HTTPLogger) Debug(msg string, keysAndValues ...interface{}) {
	level.Debug(h.Logger).Message(msg, keysAndValues...)
}

func (h *HTTPLogger) Warn(msg string, keysAndValues ...interface{}) {
	level.Info(h.Logger).Message(msg, keysAndValues...)
}
