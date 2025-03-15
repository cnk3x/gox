package fss

import (
	"fmt"
	"io"
	"log/slog"
)

func NoErr(closer io.Closer) func() {
	return func() {
		if err := closer.Close(); err != nil {
			slog.Debug("close error", "type", fmt.Sprintf("%T", closer), "err", err)
		}
	}
}
