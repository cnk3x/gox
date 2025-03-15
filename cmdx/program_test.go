package cmdx

import (
	"log/slog"
	"testing"
	"time"
)

func TestProgram(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	s := Run(t.Context(), WithOptions(Options{Execute: "sleep", Args: []string{"100s"}}))
	time.AfterFunc(time.Second*2, s.Restart)
	time.AfterFunc(time.Second*4, s.Restart)
	time.AfterFunc(time.Second*6, s.Restart)
	time.AfterFunc(time.Second*8, s.Restart)
	time.AfterFunc(time.Second*10, s.Stop)

loop:
	for {
		select {
		case code := <-s.Changed:
			slog.Debug("status", "code", code)
			if code == StatusStopped {
				break loop
			}
		}
	}

	s.Wait()
}
