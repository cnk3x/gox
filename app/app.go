package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/cnk3x/gox/flags"
	"github.com/cnk3x/gox/strs"
)

type App interface {
	Run(context.Context) error
}

func Run(app App, options ...flags.Option) {
	flags.RootSet(
		flags.FlagStruct(app, flags.EnvPrefix(strs.Replace(flags.NameDefault, "-", "_"))),
		flags.Run(func(c *flags.Command) {
			ctx, cancel := signal.NotifyContext(c.Context(), os.Interrupt, syscall.SIGTERM)
			defer cancel()
			if err := app.Run(ctx); err != nil {
				slog.Error("app exit", "err", err)
			} else {
				slog.Info("app exit")
			}
			return
		}),
	)
	flags.RootSet(options...)
	flags.Execute()
}
