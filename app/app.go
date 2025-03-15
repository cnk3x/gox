package app

import (
	"context"

	"github.com/cnk3x/gox/flags"
)

type App interface {
	Run(ctx context.Context)
}

func Run(app App, options ...flags.Option) {
	flags.RootSet(flags.Run(app.Run), flags.Flags(flags.FlagStruct(app)))
	flags.RootSet(options...)
	flags.Execute()
}
