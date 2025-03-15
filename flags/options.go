package flags

import (
	"cmp"
	"net"
	"time"

	"github.com/cnk3x/gox/strs"
)

type (
	Option     func(*Command) // Option option
	FlagOption func(*FlagSet) // FlagOption flag option
)

func (fo FlagOption) Apply(fs *FlagSet) { fo(fs) }

// Options options
func Options(options ...Option) Option {
	return func(c *Command) {
		for _, option := range options {
			option(c)
		}
	}
}

// Use is the one-line usage message.
// Recommended syntax is as follows:
//
//	[ ] identifies an optional argument. Arguments that are not enclosed in brackets are required.
//	... indicates that you can specify multiple values for the previous argument.
//	|   indicates mutually exclusive information. You can use the argument to the left of the separator or the
//	    argument to the right of the separator. You cannot use both arguments in a single use of the command.
//	{ } delimits a set of mutually exclusive arguments when one of the arguments is required. If the arguments are
//	    optional, they are enclosed in brackets ([ ]).
//
// Example: add [-F file | -D dir]... [-f format] profile
//
// Aliases is an array of aliases that can be used instead of the first word in Use.
func Use(use string, aliases ...string) Option {
	return func(c *Command) { c.Use, c.Aliases = use, aliases }
}

// Aliases is an array of aliases that can be used instead of the first word in Use.
func Aliases(aliases ...string) Option {
	return func(c *Command) { c.Aliases = aliases }
}

// Description
//
//	Short is the short description shown in the 'help' output.
//	Long is the long message shown in the 'help <this-command>' output.
func Description(short string, long ...string) Option {
	return func(c *Command) {
		c.Short = short
		if len(long) > 0 {
			c.Long = cmp.Or(long...)
		}
	}
}

// SuggestFor is an array of command names for which this command will be suggested -
// similar to aliases but only suggests.
func SuggestFor(suggestFor ...string) Option {
	return func(c *Command) { c.SuggestFor = suggestFor }
}

// The group id under which this subcommand is grouped in the 'help' output of its parent.
func GroupID(groupID string) Option {
	return func(c *Command) { c.GroupID = groupID }
}

func Hidden(hide bool) Option { return func(c *Command) { c.Hidden = hide } }

// Deprecated defines, if this command is deprecated and should print this string when used.
func Deprecated(deprecated string) Option {
	return func(c *Command) { c.Deprecated = deprecated }
}

// Example is examples of how to use the command.
func Example(example string) Option {
	return func(c *Command) { c.Example = example }
}

// Flags returns the complete FlagSet that applies
// to this command (local and persistent declared here and by all parents).
func Flags(sets ...FlagOption) Option {
	return func(c *Command) {
		f := c.Flags()
		for _, set := range sets {
			set(f)
		}
	}
}

func SortFlags(enabled bool) Option {
	return Flags(func(fs *FlagSet) { fs.SortFlags = enabled })
}

// PersistentFlags returns the persistent FlagSet specifically set in the current command.
func PersistentFlags(sets ...FlagOption) Option {
	return func(c *Command) {
		f := c.PersistentFlags()
		for _, set := range sets {
			set(f)
		}
	}
}

func SortPersistentFlags(enabled bool) Option {
	return PersistentFlags(func(fs *FlagSet) { fs.SortFlags = enabled })
}

func MarkHidden(names ...string) Option {
	return func(c *Command) {
		flag := c.Flags()
		for _, name := range names {
			flag.MarkHidden(name)
		}
	}
}

// Var is a helper function to add a variable to the flag set.
func Var[T VarT](name, short string, val T, usage, env string) FlagOption {
	return Val(&val, name, short, usage, env)
}

// Val is a helper function to add a variable to the flag set.
func Val[T VarT](val *T, name, short string, usage, env string) FlagOption {
	return func(fs *FlagSet) { anyFlag(fs, val, name, short, usage, env) }
}

// VarT is the type of the variable to be added to the flag set.
type VarT interface {
	strs.Ints |
		~bool | ~string | time.Duration | net.IP |
		~[]bool | ~[]string |
		~[]int | ~[]int32 | ~[]int64 |
		~[]uint |
		~[]float32 | ~[]float64 |
		~[]time.Duration |
		~[]net.IP
}
