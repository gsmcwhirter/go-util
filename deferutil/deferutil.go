package deferutil

import (
	"fmt"
	"os"
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
