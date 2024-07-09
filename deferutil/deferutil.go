package deferutil

import (
	"fmt"
	"os"

	"github.com/gsmcwhirter/go-util/v12/logging"
	"github.com/gsmcwhirter/go-util/v12/logging/level"
)

// CheckDefer is a wrapper for use with defer that will check error values returned and
// print to stderr if one is found
func CheckDefer(fs ...func() error) {
	for i := len(fs) - 1; i >= 0; i-- {
		if err := fs[i](); err != nil {
			if _, lastResortErr := fmt.Fprintf(os.Stderr, "Error in defer: %s\n", err); lastResortErr != nil {
				panic(lastResortErr)
			}
		}
	}
}

func CheckDeferLog(logger logging.Logger, fs ...func() error) {
	for i := len(fs) - 1; i >= 0; i-- {
		if err := fs[i](); err != nil {
			level.Error(logger).Err("Error in defer", err)
		}
	}
}
