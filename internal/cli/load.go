package cli

import (
	"os"

	"github.com/endobit/stack"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newLoadCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "load stack.yaml",
		Short: "load a stack",
		RunE: func(cmd *cobra.Command, args []string) error {
			var stack stack.Document

			b, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}

			if err := yaml.Unmarshal(b, &stack); err != nil {
				return err
			}

			return loadDocument(&stack)
		},
	}

	return &cmd
}

func loadDocument(doc *stack.Document) error {

	return nil
}
