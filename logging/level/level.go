package level

import (
	"github.com/go-kit/kit/log/level" //nolint:depguard

	"github.com/gsmcwhirter/go-util/v4/logging"
)

func Debug(logger logging.Logger) logging.Logger {
	return logging.NewFromKitLogger(level.Debug(logger))
}

func Info(logger logging.Logger) logging.Logger {
	return logging.NewFromKitLogger(level.Info(logger))
}

func Error(logger logging.Logger) logging.Logger {
	return logging.NewFromKitLogger(level.Error(logger))
}
