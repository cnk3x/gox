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
	hideHelp(rootCmd)
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

func hideHelp(c *Command) {
	c.InitDefaultHelpFlag()
	c.InitDefaultHelpCmd()
	c.Flags().MarkHidden("help")
	c.PersistentFlags().MarkHidden("help")
	for _, sub := range c.Commands() {
		if sub.Name() == "help" {
			sub.Hidden = true
		} else {
			hideHelp(sub)
		}
	}
}
