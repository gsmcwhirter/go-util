package level

import (
	"github.com/go-kit/kit/log/level" //nolint:depguard // used to implement levels

	"github.com/gsmcwhirter/go-util/v10/logging"
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
