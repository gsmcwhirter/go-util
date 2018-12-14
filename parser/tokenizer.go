package parser

import (
	"strings"

	"github.com/pkg/errors"
)

// ErrTokenizeError represents an error tokenizing where there was an opening quote without a paired closing
var ErrTokenizeError = errors.New("error tokenizing")

// Tokenize will split a string into delim-separated tokens, accounting for quoted
// sections and escapes
func Tokenize(msg string, delim rune, escape rune, quot rune) ([]string, error) {
	msgR := []rune(msg)
	max := strings.Count(msg, string(delim)) + 1
	tokens := make([]string, 0, max)

	buffer := make([]rune, 0, len(msgR))
	var r rune
	var inQuote bool

	for i := 0; i < len(msgR); i++ {
		r = msgR[i]

		if inQuote {
			if r == quot {
				inQuote = false
				continue
			}

			if r == escape && i < len(msgR)-1 && msgR[i+1] == quot {
				buffer = append(buffer, msgR[i+1])
				i++ // skip over next, having added it already
				continue
			}

			buffer = append(buffer, r)
			continue
		}

		if r == delim {
			tokens = append(tokens, string(buffer))
			buffer = buffer[:0]

			continue
		}

		if r == quot {
			inQuote = true
			continue
		}

		if r == escape && i < len(msgR)-1 && (msgR[i+1] == delim || msgR[i+1] == quot) {
			buffer = append(buffer, msgR[i+1])
			i++ // skip over next, having added it already
			continue
		}

		buffer = append(buffer, r)
	}

	if len(buffer) > 0 {
		tokens = append(tokens, string(buffer))
	}

	if inQuote {
		return tokens, errors.Wrap(ErrTokenizeError, "unmatched quotes")
	}

	return tokens, nil
}
