package main

import (
	"context"
	"log"
	"os"

	"github.com/poy/gemini/pkg/cli/commands"
	"github.com/poy/gemini/pkg/config"
	"github.com/poy/go-dependency-injection/pkg/injection"

	// Register the STDOUT logger.
	_ "github.com/poy/go-router/pkg/observability/cli"

	// Register subcommands
	_ "github.com/poy/gemini/pkg/cli/commands/config"
)

func main() {
	log.SetFlags(0)
	ctx := injection.WithInjection(context.Background())
	cfg := injection.Resolve[config.Config](ctx)
	defer cfg.Store()
	if err := commands.BuildRoot(ctx).ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
