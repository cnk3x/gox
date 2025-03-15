package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/cnk3x/gox/flags"
	"github.com/cnk3x/gox/strs"
)

type App interface {
	Run(context.Context) error
}

var envPrefix = strs.Replace(strs.Upper(strs.TrimSuffix(flags.NameDefault, "-"+runtime.GOARCH, "-"+runtime.GOOS)), "-", "_")

func SetEnvPrefix(s string) { envPrefix = s }

func Run(app App, options ...flags.Option) {
	flags.RootSet(
		flags.FlagStruct(app, flags.EnvPrefix(envPrefix)),
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
