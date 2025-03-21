package flags

import "context"

// RunFunc cobra run
type RunFunc interface {
	func() | func([]string) | func(*Command, []string) | func(*Command) |
		func() error | func([]string) error | func(*Command, []string) error | func(*Command) error |
		func(context.Context) | func(context.Context) error | func(context.Context, []string) error
}

func buildCobraRun[F RunFunc](fn F) func(*Command, []string) error {
	return func(cmd *Command, args []string) error {
		switch f := any(fn).(type) {
		case func():
			f()
		case func(args []string):
			f(args)
		case func(*Command, []string):
			f(cmd, args)
		case func(*Command):
			f(cmd)
		case func() error:
			return f()
		case func(args []string) error:
			return f(args)
		case func(*Command, []string) error:
			return f(cmd, args)
		case func(*Command) error:
			return f(cmd)
		case func(context.Context):
			f(cmd.Context())
		case func(context.Context) error:
			return f(cmd.Context())
		case func(context.Context, []string) error:
			return f(cmd.Context(), args)
		}
		return nil
	}
}

// Run run the command
func Run[T RunFunc](fn T) Option {
	return func(c *Command) { c.RunE = buildCobraRun(fn) }
}

// PreRun pre run the command
func PreRun[T RunFunc](fn T) Option {
	return func(c *Command) { c.PreRunE = buildCobraRun(fn) }
}

// PostRun post run the command
func PostRun[T RunFunc](fn T) Option {
	return func(c *Command) { c.PostRunE = buildCobraRun(fn) }
}

// PersistentPreRun persistent pre run the command
func PersistentPreRun[T RunFunc](fn T) Option {
	return func(c *Command) { c.PersistentPreRunE = buildCobraRun(fn) }
}

// PersistentPostRun persistent post run the command
func PersistentPostRun[T RunFunc](fn T) Option {
	return func(c *Command) { c.PersistentPostRunE = buildCobraRun(fn) }
}
