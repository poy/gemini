package config

import (
	"context"

	"github.com/poy/gemini/pkg/cli/commands"
	"github.com/poy/go-dependency-injection/pkg/injection"
	"github.com/spf13/cobra"
)

type SubCommand *cobra.Command

func init() {
	injection.Register(func(ctx context.Context) injection.Group[commands.SubCommand] {
		return injection.AddToGroup(ctx, buildConfig(ctx))
	})
}

func buildConfig(ctx context.Context) commands.SubCommand {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Config for CLI",
		Long:  "Config for CLI",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	for _, subCmd := range injection.Resolve[injection.Group[SubCommand]](ctx).Vals() {
		cmd.AddCommand(subCmd)
	}

	return cmd
}
