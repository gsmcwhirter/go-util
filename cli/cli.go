// Package cli is a light wrapper around the cobra argument parsing/command
// library. The goal is to make it a bit easier to add examples and to limit
// some options in the name of ease-of-use.
package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// CommandOptions is a way to specify the structure of a cli command
//
// - ShortHelp is the short description of the command (should be one line)
// - LongHelp is a longer description of the command
// - Example is a string containing examples (can also be constructed and expanded with Command.AddExamples)
// - PosArgsUsage is a string to represent the positional arguments expected
// - Deprecated is a string that will be displayed to indicate the command is deprecated
// - Args is a specifier for the number of positional arguments (see NoArgs, ExactArgs, MinimumNArgs, MaximumNArgs, RangeArgs)
// - Aliases is a list of aliases for the command
// - Hidden enables hiding the command from help messages
// - SilenceUsage will prevent a usage message from being displayed if an error occurs
type CommandOptions struct {
	ShortHelp    string
	LongHelp     string
	Example      string
	PosArgsUsage string
	Deprecated   string
	Args         cobra.PositionalArgs

	Aliases       []string
	Hidden        bool
	SilenceUsage  bool
	SilenceErrors bool
}

// NoArgs is a possible CommandOptions.Args value that indicates that no positional
// arguments should occur
//
// Pass as CommandOptions.Args = NoArgs
var NoArgs = cobra.NoArgs

// ExactArgs is a possible CommandOptions.Args value that indicates an exact number N
// of positional arguments is required
//
// Pass as CommandOptions.Args = ExactArgs(3) for exactly 3 args
var ExactArgs = cobra.ExactArgs

// MinimumNArgs is a possible CommandOptions.Args value that indicates at least N
// positional arguments are required
//
// Pass as CommandOptions.Args = MinimumNArgs(3) for at least 3 args
var MinimumNArgs = cobra.MinimumNArgs

// MaximumNArgs is a possible CommandOpitions.Args value that indicates at most N
// positional arguments are allowed
//
// Pass as CommandOptions.Args = MaximumNArgs(3) for at most 3 args
var MaximumNArgs = cobra.MaximumNArgs

// RangeArgs is a possible CommandOptions.Args vlue that indicates between M and N
// positional arguments are required
//
// Pass as CommandOptions.Args = RangeArgs(3, 5) for at least 3 and at most 5 args
var RangeArgs = cobra.RangeArgs

// Command is a thin wrapper around a cobra.Command pointer to provide
// some convenience functions
type Command struct {
	*cobra.Command
}

// NewCLI constructs a new Command struct that is intended to be the root command for a cli
// application
//
// - appName is used to name the application
// - buildVersion, buildDate, and buildSHA are used to construct a version string
// - opts determines how the command behaves
func NewCLI(appName, buildVersion, buildSHA, buildDate string, opts CommandOptions) *Command {
	c := NewCommand(appName, opts)
	c.Version = fmt.Sprintf("%s (%s, %s)", buildVersion, buildDate, buildSHA)
	return c
}

// NewCommand constructs a new Command struct that can be used as a subcommand
//
// - cmdName is the name of the command
// - opts determines how the command behaves
func NewCommand(cmdName string, opts CommandOptions) *Command {
	var use string
	if opts.PosArgsUsage != "" {
		use = fmt.Sprintf("%s %s", cmdName, opts.PosArgsUsage)
	} else {
		use = cmdName
	}

	return &Command{
		&cobra.Command{
			Use:     use,
			Short:   opts.ShortHelp,
			Long:    opts.LongHelp,
			Example: opts.Example,

			Deprecated:    opts.Deprecated,
			SilenceUsage:  opts.SilenceUsage,
			SilenceErrors: opts.SilenceErrors,

			Args: opts.Args,

			Aliases: opts.Aliases,
			Hidden:  opts.Hidden,

			RunE: func(cmd *cobra.Command, _ []string) error {
				return cmd.Help()
			},
		},
	}
}

// AddExamples extends the existing Command.Example value with more examples
//
// - descCmds is a list whose alternating elements are example descriptions and
// example invocations (so the first, third, fifth entries, etc., are the descriptions,
// and the second, fourth, sixth, etc are the commands being described). If there are an
// uneven number of descriptions and commands, the extra descriptions are ignored.
func (c *Command) AddExamples(descCmds ...string) {
	b := strings.Builder{}
	_, _ = b.WriteString(c.Example)
	for i := 0; i < len(descCmds)/2*2; i += 2 {
		_, _ = b.WriteString(fmt.Sprintf(`
  %s:
	$ %s
`, descCmds[i], descCmds[i+1]))
	}
	c.Example = b.String()
}

// SetRunFunc sets the function that is called when the command executes.
func (c *Command) SetRunFunc(run func(cmd *Command, args []string) error) {
	c.RunE = func(ccmd *cobra.Command, args []string) error {
		return run(&Command{ccmd}, args)
	}
}

// AddSubCommands attaches a list of Command structs as subcommands to the current command
func (c *Command) AddSubCommands(cmds ...*Command) {
	ccmds := make([]*cobra.Command, 0, len(cmds))
	for _, cmd := range cmds {
		ccmds = append(ccmds, cmd.Command)
	}
	c.AddCommand(ccmds...)
}
