package level

import (
	"github.com/go-kit/log/level" //nolint:depguard,staticcheck // used to implement levels

	"github.com/gsmcwhirter/go-util/v11/logging"
)

func Debug(logger logging.Logger) logging.Logger {
	return logging.NewFrom(level.Debug(logging.BaseFrom(logger)))
}

func Info(logger logging.Logger) logging.Logger {
	return logging.NewFrom(level.Info(logging.BaseFrom(logger)))
}

func Error(logger logging.Logger) logging.Logger {
	return logging.NewFrom(level.Error(logging.BaseFrom(logger)))
}
