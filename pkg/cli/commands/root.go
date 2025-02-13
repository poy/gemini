package commands

import (
	"context"

	"github.com/poy/go-dependency-injection/pkg/injection"
	"github.com/spf13/cobra"
)

type SubCommand *cobra.Command

func BuildRoot(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gemini",
		Short: "Gemini CLI",
		Long:  "Gemini CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	for _, subCmd := range injection.Resolve[injection.Group[SubCommand]](ctx).Vals() {
		cmd.AddCommand(subCmd)
	}

	return cmd
}
