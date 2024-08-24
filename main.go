package main

import (
	"github.com/spf13/cobra"
)

func init() {
	DeepSpace.AddCommand(
		startCommand(),
		listCommand(),
		inspectCommand(),
		cleanupCommand(),
		exportCommand(),
	)
}

var (
	DeepSpace = &cobra.Command{
		Use:           "deepspace",
		Version:       "v0.0.1",
		Short:         "DeepSpace is a command-line tool for debugging the DeepSeek AI HTTP API",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
)

func main() {
	if err := DeepSpace.Execute(); err != nil {
		logFatal(err)
	}
}
