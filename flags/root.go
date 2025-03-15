package flags

import (
	"os"
	"path/filepath"

	"github.com/cnk3x/gox/strs"
)

var rootCmd = &Command{}

var NameDefault = strs.TrimExe(filepath.Base(os.Args[0]))

func init() {
	rootCmd = &Command{Use: NameDefault}
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func RootSet(options ...Option) {
	for _, option := range options {
		option(rootCmd)
	}
}

func AddCommand(use string, options ...Option) *Command {
	c := &Command{Use: use}
	for _, option := range options {
		option(c)
	}
	rootCmd.AddCommand(c)
	return c
}
