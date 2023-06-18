// Package main is the entrypoint for the stack command line tool.
package main

import (
	"os"

	"github.com/spf13/cobra"
)

var version string

func main() {
	root := newRootCmd()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	cmd := cobra.Command{
		Version: version,
		Use:     "start",
		Short:   "stack is a command line tool for interacting with stackd",
	}

	return &cmd
}
