package parser

import (
	"strings"

	"github.com/gsmcwhirter/go-util/v6/errors"
)

// ErrTokenizeError represents an error tokenizing where there was an opening quote without a paired closing
var ErrTokenizeError = errors.New("error tokenizing")

func isQuote(char rune, quots []rune) bool {
	for _, quot := range quots {
		if char == quot {
			return true
		}
	}

	return false
}

// Tokenize will split a string into delim-separated tokens, accounting for quoted
// sections and escapes
func Tokenize(msg string, delim, escape rune, quots []rune) ([]string, error) {
	msgR := []rune(msg)
	max := strings.Count(msg, string(delim)) + 1
	tokens := make([]string, 0, max)

	buffer := make([]rune, 0, len(msgR))
	var r rune
	var inQuote bool

	for i := 0; i < len(msgR); i++ {
		r = msgR[i]

		if inQuote {
			if isQuote(r, quots) {
				inQuote = false
				continue
			}

			if r == escape && i < len(msgR)-1 && isQuote(msgR[i+1], quots) {
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

		if isQuote(r, quots) {
			inQuote = true
			continue
		}

		if r == escape && i < len(msgR)-1 && (msgR[i+1] == delim || isQuote(msgR[i+1], quots)) {
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
