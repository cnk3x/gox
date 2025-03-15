package cmdx

import (
	"log/slog"

	"github.com/cnk3x/gox/arrs"
	"github.com/cnk3x/gox/strs"
)

type Env []string

func (e Env) Set(k, v string) Env { return arrs.ReplaceOrAppend(e, k+"="+v, strs.PrefixMatch(k)) }

func (e Env) Del(keys ...string) Env { return arrs.DeleteFunc(e, strs.PrefixMatch(keys...)) }

func (e Env) Sets(env ...string) Env {
	for _, item := range env {
		if k, v, ok := strs.Cut(item, "="); ok {
			e = e.Set(k, v)
		} else {
			slog.Debug("env sets ignored", "env", item)
		}
	}
	return e
}

func (e Env) Compact() Env {
	return arrs.CleanFunc(e, func(a, b string) bool { return strs.Equal(strs.Prefix2(a, b, "=")) })
}
