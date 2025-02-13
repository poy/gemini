package config

import (
	"context"
	"fmt"

	"github.com/poy/gemini/pkg/config"
	"github.com/poy/go-dependency-injection/pkg/injection"
	"github.com/spf13/cobra"
)

func init() {
	injection.Register(func(ctx context.Context) injection.Group[SubCommand] {
		return injection.AddToGroup(ctx, buildGet(ctx))
	})
}

func buildGet(ctx context.Context) SubCommand {
	cfg := injection.Resolve[config.Config](ctx)

	cmd := &cobra.Command{
		Use:   "get <KEY_NAME>",
		Short: "Get configuration value",
		Long:  "Get configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			fmt.Fprintln(cmd.OutOrStdout(), cfg.Get(args[0]))
			return nil
		},
	}
	return cmd
}
