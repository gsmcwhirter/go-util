package parser

import (
	"errors"
	"strings"
)

// Parser TODOC
type Parser interface {
	ParseCommand(line string) (cmd string, rest string, err error)
	KnownCommand(cmd string) bool
	LearnCommand(cmd string)
	LeadChar() string
	IsCaseSensitive() bool
}

// ErrNotACommand TODOC
var ErrNotACommand = errors.New("not a command")

// ErrUnknownCommand TODOC
var ErrUnknownCommand = errors.New("unknown command")

type parser struct {
	CmdIndicator  string
	knownCommands map[string]bool
	caseSensitive bool
}

// Options TODOC
type Options struct {
	CmdIndicator  string
	KnownCommands []string
	CaseSensitive bool
}

// NewParser TODOC
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

// IsCaseSensitive TODOC
func (p *parser) IsCaseSensitive() bool {
	return p.caseSensitive
}

// KnownCommand TODOC
func (p *parser) KnownCommand(cmd string) bool {
	if p.caseSensitive {
		return p.knownCommands[cmd]
	}
	return p.knownCommands[strings.ToLower(cmd)]
}

func (p *parser) LearnCommand(cmd string) {
	if p.caseSensitive {
		p.knownCommands[cmd] = true
	}
	p.knownCommands[strings.ToLower(cmd)] = true
}

// LeadChar returns the character that identifies commands
func (p *parser) LeadChar() string {
	return p.CmdIndicator
}

func (p *parser) ParseCommand(line string) (cmd string, rest string, err error) {
	if len(line) == 0 {
		if p.CmdIndicator == "" && p.KnownCommand("") {
			return "", line, nil
		}

		err = ErrNotACommand
		return
	}

	if !strings.HasPrefix(line, p.CmdIndicator) {
		err = ErrNotACommand
		return
	}

	for i := range line {
		if i == 0 {
			continue
		}

		if line[i] == ' ' {
			cmd = line[1:i]
			rest = line[i:]

			if p.KnownCommand(cmd) {
				return
			}
		}
	}

	cmd = line[1:]
	rest = line[0:0]
	if !p.KnownCommand(cmd) {
		err = ErrUnknownCommand
	}
	return
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

// MaybeCount TODOC
func MaybeCount(line string) (l string, c string) {
	l = line

	for i := len(line) - 1; i >= 0; i-- {
		_, isDigit := digits[line[i]]
		if !isDigit {
			if line[i] == 'x' {
				l = line[:i]
				c = line[i+1:]
			} else {
				l = line[:i+1]
				c = line[i+1:]
			}
			return
		}
	}

	return
}
