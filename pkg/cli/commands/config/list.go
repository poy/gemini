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
		return injection.AddToGroup(ctx, buildList(ctx))
	})
}

func buildList(ctx context.Context) SubCommand {
	cfg := injection.Resolve[config.Config](ctx)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List configuration keys",
		Long:  "List configuration keys",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			for _, k := range cfg.List() {
				fmt.Fprintln(cmd.OutOrStdout(), k)
			}
			return nil
		},
	}
	return cmd
}
