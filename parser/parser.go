package parser

import (
	"strings"

	"github.com/pkg/errors"
)

// Parser is an interface describing a repl/text interface command parser
type Parser interface {
	ParseCommand(line string) (cmd string, err error)
	KnownCommand(cmd string) bool
	LearnCommand(cmd string)
	LeadChar() string
	IsCaseSensitive() bool
}

// ErrNotACommand is the error returned when the string to be parsed does not
// represent a command syntactically
var ErrNotACommand = errors.New("not a command")

// ErrUnknownCommand is the error returned when the string to be parsed is a command
// syntactically but that command is not registered
var ErrUnknownCommand = errors.New("unknown command")

type parser struct {
	CmdIndicator  string
	knownCommands map[string]bool
	caseSensitive bool
}

// Options specifies options when constructing a new Parser
//
// - CmdIndicator is the string that all commands must be prefixed with (e.g., "!" or "/")
// - KnownCommands is a list of commands that should be recognized (expanded with LearnCommand)
// - CaseSensitive will make the parser case-sensitive when determining if it knows about a command
type Options struct {
	CmdIndicator  string
	KnownCommands []string
	CaseSensitive bool
}

// NewParser constructs a new Parser
func NewParser(opts Options) Parser {
	p := &parser{
		CmdIndicator:  opts.CmdIndicator,
		knownCommands: map[string]bool{},
		caseSensitive: opts.CaseSensitive,
	}

	for _, cmd := range opts.KnownCommands {
		p.LearnCommand(cmd)
	}
	return p
}

// IsCaseSensitive reports whether the parser is case-sensitive or not
func (p *parser) IsCaseSensitive() bool {
	return p.caseSensitive
}

// KnownCommand reports whether the parser knows about a command or not
func (p *parser) KnownCommand(cmd string) bool {
	if p.caseSensitive {
		return p.knownCommands[cmd]
	}
	return p.knownCommands[strings.ToLower(cmd)]
}

// LearnCommand adds the command to the parser's list of known commands if it
// is not already present
func (p *parser) LearnCommand(cmd string) {
	if p.caseSensitive {
		p.knownCommands[cmd] = true
		return
	}
	p.knownCommands[strings.ToLower(cmd)] = true
}

// LeadChar returns the character that identifies commands
func (p *parser) LeadChar() string {
	return p.CmdIndicator
}

// ParseCommand attempts to parse a user-entered string as a command
func (p *parser) ParseCommand(line string) (string, error) {
	var cmd string
	var err error

	if !strings.HasPrefix(line, p.CmdIndicator) {
		return "", ErrNotACommand
	}

	cmd = strings.TrimPrefix(line, p.CmdIndicator)
	if !p.KnownCommand(cmd) {
		err = ErrUnknownCommand
	}

	if !p.IsCaseSensitive() {
		cmd = strings.ToLower(cmd)
	}

	return cmd, err
}

var digits = map[byte]bool{
	'0': true,
	'1': true,
	'2': true,
	'3': true,
	'4': true,
	'5': true,
	'6': true,
	'7': true,
	'8': true,
	'9': true,
}

// MaybeCount attempts to split some text that might contain a "count" at the end
//
// Recognized "count" format is a string of digits, possibly preceded by an x, +, or -, preceeded by a space
func MaybeCount(line string) (string, string) {
	l := line
	c := ""

	for i := len(line) - 1; i >= 0; i-- {
		_, isDigit := digits[line[i]]
		if !isDigit {
			switch line[i] {
			case ' ':
				l = line[:i+1]
				c = line[i+1:]
			case 'x', '+':
				if i > 0 && line[i-1] != ' ' {
					l = line
					c = ""
				} else {
					l = line[:i]
					c = line[i+1:]
				}
			case '-':
				if i > 0 && line[i-1] != ' ' {
					l = line
					c = ""
				} else {
					l = line[:i]
					c = line[i:]
				}
			default:
				l = line
				c = ""
			}
			return l, c
		}
	}

	return l, c
}
