package cli

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"strings"
)

// Command is a command or subcommand that can be run with Execute.
type Command struct {
	// Use is the one-line usage message.
	// Recommended syntax is as follow:
	//   [ ] identifies an optional argument. Arguments that are not enclosed in brackets are required.
	//   ... indicates that you can specify multiple values for the previous argument.
	//   |   indicates mutually exclusive information. You can use the argument to the left of the separator or the
	//       argument to the right of the separator. You cannot use both arguments in a single use of the command.
	//   { } delimits a set of mutually exclusive arguments when one of the arguments is required. If the arguments are
	//       optional, they are enclosed in brackets ([ ]).
	// Example: add [-F file | -D dir]... [-f format] <profile>
	Usage string

	// Short is the short description shown in the 'help' output.
	Short string

	// Long is the long message shown in the 'help <this-command>' output.
	Long string

	// Hidden defines, if this command is hidden and should NOT show up in the list of available commands.
	Hidden bool

	// Aliases is an array of aliases that can be used instead of the first word in Use.
	Aliases []string

	// Example is examples of how to use the command.
	Example string

	// Annotations are key/value pairs that can be used by applications to identify or
	// group commands.
	Annotations map[string]interface{}

	// Version defines the version for this command. If this value is non-empty and the command does not
	// define a "version" flag, a "version" boolean flag will be added to the command and, if specified,
	// will print content of the "Version" variable. A shorthand "v" flag will also be added if the
	// command does not define one.
	Version string

	// Expected arguments
	Args PositionalArgs

	// Run is the function that performs the command
	Run func(ctx context.Context, args []string)

	commands []*Command
	parent   *Command
	flags    *flag.FlagSet
}

// Flags returns the complete FlagSet that applies to this command.
func (c *Command) Flags() *flag.FlagSet {
	if c.flags == nil {
		c.flags = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		var null bytes.Buffer
		c.flags.SetOutput(&null)
	}
	return c.flags
}

// AddCommand adds one or more commands to this parent command.
func (c *Command) AddCommand(sub *Command) {
	if sub == c {
		panic("command can't be a child of itself")
	}
	sub.parent = c
	c.commands = append(c.commands, sub)
}

// CommandPath returns the full path to this command.
func (c *Command) CommandPath() string {
	if c.parent != nil {
		return c.parent.CommandPath() + " " + c.Name()
	}
	return c.Name()
}

// UseLine puts out the full usage for a given command (including parents).
func (c *Command) UseLine() string {
	use := c.Usage
	if use == c.Name() && len(c.commands) > 0 {
		use = fmt.Sprintf("%s [command]", c.Name())
	}
	if c.parent != nil {
		return c.parent.CommandPath() + " " + use
	} else {
		return use
	}
}

// Name returns the command's name: the first word in the usage line.
func (c *Command) Name() string {
	name := c.Usage
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

// Find the target command given the args and command tree.
// Meant to be run on the highest node. Only searches down.
// Also returns the arguments consumed to reach the command.
func (c *Command) Find(args []string) (cmd *Command, n int) {
	cmd = c
	if len(args) == 0 {
		return
	}
	var arg string
	for n, arg = range args {
		if cc := cmd.findSub(arg); cc != nil {
			cmd = cc
		} else {
			return
		}
	}
	n += 1
	return
}

func (c *Command) findSub(name string) *Command {
	for _, cmd := range c.commands {
		if cmd.Name() == name || hasAlias(cmd, name) {
			return cmd
		}
	}
	return nil
}

func hasAlias(cmd *Command, name string) bool {
	for _, a := range cmd.Aliases {
		if a == name {
			return true
		}
	}
	return false
}

// PositionalArgs is a function type used by the Command Args field
// for detecting whether the arguments match a given expectation.
type PositionalArgs func(cmd *Command, args []string) error

// MinArgs returns an error if there is not at least N args.
func MinArgs(n int) PositionalArgs {
	return func(cmd *Command, args []string) error {
		if len(args) < n {
			return fmt.Errorf("requires at least %d arg(s), only received %d", n, len(args))
		}
		return nil
	}
}

// MaxArgs returns an error if there are more than N args.
func MaxArgs(n int) PositionalArgs {
	return func(cmd *Command, args []string) error {
		if len(args) > n {
			return fmt.Errorf("accepts at most %d arg(s), received %d", n, len(args))
		}
		return nil
	}
}

// ExactArgs returns an error if there are not exactly n args.
func ExactArgs(n int) PositionalArgs {
	return func(cmd *Command, args []string) error {
		if len(args) != n {
			return fmt.Errorf("accepts %d arg(s), received %d", n, len(args))
		}
		return nil
	}
}

// RangeArgs returns an error if the number of args is not within the expected range.
func RangeArgs(min int, max int) PositionalArgs {
	return func(cmd *Command, args []string) error {
		if len(args) < min || len(args) > max {
			return fmt.Errorf("accepts between %d and %d arg(s), received %d", min, max, len(args))
		}
		return nil
	}
}
