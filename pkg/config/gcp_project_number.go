package config

import (
	"context"
	"os/exec"
	"strings"

	"github.com/poy/go-dependency-injection/pkg/injection"
	"github.com/poy/go-router/pkg/observability"
)

type GCPProjectNumber string

func init() {
	injection.Register(func(ctx context.Context) GCPProjectNumber {
		log := injection.Resolve[observability.Logger](ctx)
		cfg := injection.Resolve[Config](ctx)
		val := cfg.Get("GCPProjectNumber")
		if val == "" {
			cfg := injection.Resolve[Config](ctx)
			log.Infof("GCPProjectNumber is empty, reading from gcloud to get value")
			cmd := exec.CommandContext(ctx, "gcloud", "config", "get-value", "project")
			output, err := cmd.Output()
			if err != nil {
				log.Fatalf("failed to fetch GCP project ID from gcloud: %v", err)
			}
			gcpProjectID := strings.TrimSpace(string(output))
			if gcpProjectID == "" {
				log.Fatalf("empty project ID. Please set it with the following command\ngcloud config set-value project <PROJECT>")
			}

			log.Infof("Read project ID: %s", gcpProjectID)

			// Now we need to lookup the project number.
			cmd = exec.CommandContext(ctx, "gcloud", "projects", "list", "--filter="+gcpProjectID, "--format=value(PROJECT_NUMBER)")
			output, err = cmd.CombinedOutput()
			if err != nil {
				println(string(output))
				log.Fatalf("failed to fetch GCP project number from gcloud: %v", err)
			}
			val = strings.TrimSpace(string(output))
			log.Infof("Read project number: h%s", val)
			cfg.Set("GCPProjectNumber", val)
		}

		return GCPProjectNumber(val)
	})
	injection.Register(func(ctx context.Context) injection.Group[KeyName] {
		return injection.AddToGroup(ctx, KeyName("GCPProjectNumber"))
	})
}
