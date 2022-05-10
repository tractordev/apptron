package cli

import (
	"flag"
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/template"
	"unicode"
)

// HelpFuncs are used by the help templating system.
var HelpFuncs = template.FuncMap{
	"trim": strings.TrimSpace,
	"trimRight": func(s string) string {
		return strings.TrimRightFunc(s, unicode.IsSpace)
	},
	"padRight": func(s string, padding int) string {
		template := fmt.Sprintf("%%-%ds", padding)
		return fmt.Sprintf(template, s)
	},
}

// HelpTemplate is a template used to generate help.
var HelpTemplate = `Usage:{{if .Runnable}}
{{.UseLine}}{{end}}{{if .HasSubCommands}}
{{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
{{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasSubCommands}}

Available Commands:{{range .Commands}}{{if (or .Available (eq .Name "help"))}}
{{padRight .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasFlags}}

Flags:
{{.FlagUsages | trimRight }}{{end}}{{if .HasSubCommands}}

Use "{{.CommandPath}} [command] -help" for more information about a command.{{end}}
`

// CommandHelp wraps a Command to generate help.
type CommandHelp struct {
	*Command
}

// WriteHelp generates help for the command written to an io.Writer.
func (c *CommandHelp) WriteHelp(w io.Writer) error {
	t := template.Must(template.New("help").Funcs(HelpFuncs).Parse(HelpTemplate))
	return t.Execute(w, c)
}

// Runnable determines if the command is itself runnable.
func (c *CommandHelp) Runnable() bool {
	return c.Run != nil
}

// Available determines if a command is available as a non-help command (this includes all non hidden commands).
func (c *CommandHelp) Available() bool {
	if c.Hidden {
		return false
	}
	if c.Runnable() || c.HasSubCommands() {
		return true
	}
	return false
}

// HasSubCommands determines if a command has available sub commands that need to be
// shown in the usage/help default template under 'available commands'.
func (c *CommandHelp) HasSubCommands() bool {
	for _, sub := range c.commands {
		if (&CommandHelp{sub}).Available() {
			return true
		}
	}
	return false
}

// NameAndAliases returns a list of the command name and all aliases.
func (c *CommandHelp) NameAndAliases() string {
	return strings.Join(append([]string{c.Name()}, c.Aliases...), ", ")
}

// HasExample determines if the command has example.
func (c *CommandHelp) HasExample() bool {
	return len(c.Example) > 0
}

// Commands returns any subcommands as CommandHelp values.
func (c *CommandHelp) Commands() (cmds []*CommandHelp) {
	for _, cmd := range c.commands {
		cmds = append(cmds, &CommandHelp{cmd})
	}
	return
}

// NamePadding returns padding for the name.
func (c *CommandHelp) NamePadding() int {
	// TODO: consider making this dynamic, based on length of all sibling commands
	return 16
}

// HasFlags checks if the command contains flags.
func (c *CommandHelp) HasFlags() bool {
	n := 0
	c.Flags().VisitAll(func(f *flag.Flag) {
		n++
	})
	return n > 0
}

// FlagUsages creates a string for flag usage help.
func (c *CommandHelp) FlagUsages() string {
	var sb strings.Builder
	c.Flags().VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(&sb, "  -%s", f.Name) // Two spaces before -; see next two comments.
		name, usage := flag.UnquoteUsage(f)
		if len(name) > 0 {
			sb.WriteString(" ")
			sb.WriteString(name)
		}
		// Boolean flags of one ASCII letter are so common we
		// treat them specially, putting their usage on the same line.
		if sb.Len() <= 4 { // space, space, '-', 'x'.
			sb.WriteString("\t")
		} else {
			// Four spaces before the tab triggers good alignment
			// for both 4- and 8-space tab stops.
			sb.WriteString("\n    \t")
		}
		sb.WriteString(strings.ReplaceAll(usage, "\n", "\n    \t"))
		f.Usage = ""
		if !isZeroValue(f, f.DefValue) {
			typ, _ := flag.UnquoteUsage(f)
			if typ == "string" {
				// put quotes on the value
				fmt.Fprintf(&sb, " (default %q)", f.DefValue)
			} else {
				fmt.Fprintf(&sb, " (default %v)", f.DefValue)
			}
		}
		sb.WriteString("\n")
	})
	return sb.String()
}

// isZeroValue determines whether the string represents the zero
// value for a flag.
func isZeroValue(f *flag.Flag, value string) bool {
	// Build a zero value of the flag's Value type, and see if the
	// result of calling its String method equals the value passed in.
	// This works unless the Value type is itself an interface type.
	typ := reflect.TypeOf(f.Value)
	var z reflect.Value
	if typ.Kind() == reflect.Ptr {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}
	return value == z.Interface().(flag.Value).String()
}
