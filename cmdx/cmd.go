package cmdx

import (
	"cmp"
	"context"
	"errors"
	"io"
	"log/slog"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/cnk3x/gox/chans"
	"github.com/cnk3x/gox/fss"
	"github.com/cnk3x/gox/strs"

	"github.com/valyala/fasttemplate"
)

var (
	ErrDone    = errors.New("done")
	ErrStop    = errors.New("user stop")
	ErrRestart = errors.New("user restart")
	ErrStatus  = errors.New("error status")
)

type Options struct {
	Execute      string        `json:"execute,omitempty" yaml:"execute,omitempty"`
	Args         []string      `json:"args,omitempty" yaml:"args,omitempty"`
	Dir          string        `json:"dir,omitempty" yaml:"dir,omitempty"`
	Env          []string      `json:"env,omitempty" yaml:"env,omitempty"`
	RestartDelay time.Duration `json:"restart_delay,omitempty" yaml:"restart_delay,omitempty"`
	Logger       *Logger       `json:"logger,omitempty" yaml:"logger,omitempty"`

	PreStart []func(c *exec.Cmd) error `json:"-" yaml:"-"`
}

type Logger struct {
	*RotateOptions `json:",inline" yaml:",inline"`
	Stderr         *RotateOptions `json:"stderr,omitempty" yaml:"stderr,omitempty"`
	Stdout         *RotateOptions `json:"stdout,omitempty" yaml:"stdout,omitempty"`
}

// Status is a status code.
type Status int

const (
	StatusUnknown    Status = iota // unknown
	StatusStarting                 // starting
	StatusRunning                  // running
	StatusRestarting               // restarting
	StatusStopping                 // stopping
	StatusStopped                  // stopped
)

// String returns the string representation of the status code.
func (s Status) String() string {
	switch s {
	case StatusStarting:
		return "starting"
	case StatusRunning:
		return "running"
	case StatusRestarting:
		return "restarting"
	case StatusStopping:
		return "stopping"
	case StatusStopped:
		return "stopped"
	default:
		return "unknown"
	}
}

type Result struct {
	Command string

	Restart func()
	Stop    func()
	Wait    func()

	Status Status // status
	Pid    int    // pid
	Exit   int    // exit code
	Err    error  // error

	StartTs int64
	StopTs  int64

	Changed <-chan Status
}

type Option func(*Options)

func Run(ctx context.Context, options ...Option) *Result {
	var opts Options
	for _, apply := range options {
		apply(&opts)
	}
	return cmdRun(ctx, opts)
}

func WithOptions(options Options) Option {
	return func(opts *Options) { *opts = options }
}

func cmdRun(ctx context.Context, options Options) (s *Result) {
	var (
		allDone, closeAllDone = chans.StructChan()
		statusc, closeStatusc = chans.MakeChan[Status](5)
		pDone                 <-chan struct{}
		pCancel               context.CancelFunc
	)

	s = &Result{Changed: statusc}

	statusUpdate := func(status Status) {
		s.Status = status
		select {
		case statusc <- status:
		default:
			_ = <-statusc
			statusc <- status
		}
		if status == StatusStopped {
			closeAllDone()
			closeStatusc()
		}
	}

	run := func() {
		s.Status = StatusUnknown
		s.Err = nil
		s.Exit = 0
		s.StopTs = 0
		s.StartTs = time.Now().UnixNano()

		var (
			dir      = filepath.Clean(options.Dir)
			replArgs = map[string]string{"dir": dir}
			execute  = filepath.Clean(strRepl(options.Execute, replArgs))
			args     = strReplAll(options.Args, replArgs)
			env      = Env(os.Environ()).Sets(strReplAll(options.Env, replArgs)...)
		)

		done, closeDone := chans.StructChan()
		pDone = done

		ctx, cancel := context.WithCancel(ctx)
		pCancel = cancel
		chans.AfterChan(done, cancel)

		c := setProcessGroup(exec.CommandContext(ctx, execute, args...))
		c.Dir, c.Env = dir, env
		c.Cancel = func() error { return terminateProcess(c.Process.Pid) }

		if options.Logger != nil {
			loggerFactory := createLoggerFactory()
			c.Stdout = loggerFactory.Create(options.Logger.Stdout, options.Logger.RotateOptions)
			c.Stderr = loggerFactory.Create(options.Logger.Stderr, options.Logger.RotateOptions)
			chans.AfterChan(done, fss.NoErr(loggerFactory))
		}

		s.Command = c.String()

		if s.Status != StatusRestarting {
			statusUpdate(StatusStarting)
		}

		started, closeStarted := chans.StructChan()
		go func() {
			defer closeDone()

			err := func() (err error) {
				defer closeStarted()
				for _, ps := range options.PreStart {
					if err = ps(c); err != nil {
						return
					}
				}

				if err = c.Start(); err != nil {
					return
				}

				s.Pid = c.Process.Pid
				return
			}()

			if s.Err = err; s.Err != nil {
				statusUpdate(StatusStopped)
				return
			}

			statusUpdate(StatusRunning)

			err = c.Wait()
			if err != nil {
				if s.Status == StatusRestarting {
					return
				}

				var ee *exec.ExitError
				if errors.As(err, &ee) {
					s.Exit = ee.ExitCode()
				}
				s.Err = err
			}

			if s.Status != StatusRestarting {
				statusUpdate(StatusStopped)
			}
		}()
		<-started
	}

	s.Restart = func() {
		if s.Status != StatusRunning {
			return
		}
		statusUpdate(StatusRestarting)
		pCancel()
		<-pDone
		run()
	}

	s.Stop = func() {
		if s.Status != StatusRunning {
			return
		}
		statusUpdate(StatusStopping)
		pCancel()
	}

	s.Wait = func() { <-allDone }

	run()

	return
}

/** logger **/

type LoggerFactory struct {
	create func(options ...*RotateOptions) (w io.Writer)
	close  func() error
}

func (f LoggerFactory) Create(options ...*RotateOptions) io.Writer { return f.create(options...) }
func (f LoggerFactory) Close() error                               { return f.close() }

func createLoggerFactory() (factory LoggerFactory) {
	writers := make(map[string]io.WriteCloser, 2)

	factory.create = func(options ...*RotateOptions) (w io.Writer) {
		var opts RotateOptions
		for _, it := range options {
			if it != nil {
				opts.Path = cmp.Or(strs.TrimSpace(opts.Path), strs.TrimSpace(it.Path))
				opts.MaxBackups = cmp.Or(opts.MaxBackups, it.MaxBackups)
				opts.MaxSize = cmp.Or(opts.MaxSize, it.MaxSize)
				opts.Std = cmp.Or(opts.Std, it.Std)
			}
		}

		if p := strs.Lower(opts.Path); p != "" {
			if opts.Path, _ = filepath.Abs(opts.Path); opts.Path != "" {
				if w = writers[opts.Path]; w == nil {
					w = Rotate(opts)
					writers[opts.Path] = w.(io.WriteCloser)
					slog.Debug("[cmdx] created rotate writer", "path", opts.Path, "size", opts.MaxSize, "backups", opts.MaxBackups)
				}
			}
		}

		if opts.Std != "" {
			var std io.Writer
			switch strs.Lower(strs.TrimSpace(opts.Std)) {
			case "err":
				std = os.Stderr
			case "out":
				std = os.Stdout
			}

			if std != nil {
				if w != nil {
					w = io.MultiWriter(w, std)
				} else {
					w = std
				}
				slog.Debug("[cmdx] created std writer", "std", opts.Std)
			}
		}
		return
	}

	factory.close = func() (err error) {
		var errs []error
		for writer := range maps.Values(writers) {
			errs = append(errs, writer.Close())
		}
		clear(writers)
		return errors.Join(errs...)
	}

	return
}

/** var replaces **/

func strRepl(src string, args map[string]string, keepUnknownTags ...bool) string {
	keepUnknown := cmp.Or(keepUnknownTags...)
	output, err := fasttemplate.ExecuteFuncStringWithErr(src, "{", "}", func(w io.Writer, tag string) (int, error) {
		sTag := strs.TrimSpace(tag)
		if args != nil {
			if v, found := args[sTag]; found {
				return w.Write([]byte(v))
			}
		}

		if v, found := os.LookupEnv(sTag); found {
			return w.Write([]byte(v))
		}

		if keepUnknown {
			return w.Write([]byte(`{` + tag + `}`))
		}

		return 0, nil
	})
	if err != nil {
		output = src
	}
	return output
}

func strReplAll(src []string, args map[string]string, keepUnknownTags ...bool) (dst []string) {
	dst = make([]string, len(src))
	for i, it := range src {
		dst[i] = strRepl(it, args, keepUnknownTags...)
	}
	return dst
}

/** msic **/

func Elapsed(t ...time.Time) time.Duration {
	var start, end time.Time

	if len(t) > 0 {
		start = t[0]
	}

	if len(t) > 1 {
		end = t[1]
	}

	if start.IsZero() {
		return 0
	}

	if end.IsZero() {
		return time.Since(start)
	}

	return end.Sub(start)
}
