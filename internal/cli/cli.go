// Package cli provides the command line interface for stack.
package cli

import (
	"github.com/spf13/cobra"
)

// NewRootCmd creates a new root command.
func NewRootCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "start",
		Short: "stack is a command line tool for interacting with stackd",
	}

	cmd.AddCommand(newDumpCmd())
	cmd.AddCommand(newLoadCmd())

	return &cmd
}

func newDumpCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "dump",
		Short: "dump a stack",
	}

	return &cmd
}
