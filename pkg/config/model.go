package config

import (
	"context"

	"github.com/poy/go-dependency-injection/pkg/injection"
	"github.com/poy/go-router/pkg/observability"
)

type Model string

func init() {
	injection.Register(func(ctx context.Context) Model {
		log := injection.Resolve[observability.Logger](ctx)
		cfg := injection.Resolve[Config](ctx)
		val := cfg.Get("Model")
		if val == "" {
			cfg := injection.Resolve[Config](ctx)
			log.Infof("Model is empty, using gemini-2.0-flash")
			val = "gemini-2.0-flash"
			cfg.Set("Model", val)
		}

		return Model(val)
	})
	injection.Register(func(ctx context.Context) injection.Group[KeyName] {
		return injection.AddToGroup(ctx, KeyName("Model"))
	})
}
