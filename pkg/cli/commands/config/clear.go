package config

import (
	"context"

	"github.com/poy/gemini/pkg/config"
	"github.com/poy/go-dependency-injection/pkg/injection"
	"github.com/spf13/cobra"
)

func init() {
	injection.Register(func(ctx context.Context) injection.Group[SubCommand] {
		return injection.AddToGroup(ctx, buildClear(ctx))
	})
}

func buildClear(ctx context.Context) SubCommand {
	cfg := injection.Resolve[config.Config](ctx)

	cmd := &cobra.Command{
		Use:   "clear <KEY_NAME>",
		Short: "Clear configuration value",
		Long:  "Clear configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cfg.Set(args[0], "")
			return nil
		},
	}
	return cmd
}
