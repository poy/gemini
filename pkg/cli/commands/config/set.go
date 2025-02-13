package config

import (
	"context"

	"github.com/poy/gemini/pkg/config"
	"github.com/poy/go-dependency-injection/pkg/injection"
	"github.com/spf13/cobra"
)

func init() {
	injection.Register(func(ctx context.Context) injection.Group[SubCommand] {
		return injection.AddToGroup(ctx, buildSet(ctx))
	})
}

func buildSet(ctx context.Context) SubCommand {
	cfg := injection.Resolve[config.Config](ctx)

	cmd := &cobra.Command{
		Use:   "set <KEY_NAME> <VALUE>",
		Short: "Set configuration value",
		Long:  "Set configuration value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cfg.Set(args[0], args[1])
			return nil
		},
	}
	return cmd
}
