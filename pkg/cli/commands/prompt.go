package commands

import (
	"context"
	"fmt"

	"cloud.google.com/go/vertexai/genai"
	"github.com/poy/gemini/pkg/config"
	"github.com/poy/go-dependency-injection/pkg/injection"
	"github.com/spf13/cobra"
)

func init() {
	injection.Register(
		func(ctx context.Context) injection.Group[SubCommand] {
			return injection.AddToGroup(ctx, buildPrompt(ctx))
		},
	)
}

func buildPrompt(ctx context.Context) SubCommand {
	var model string
	cmd := &cobra.Command{
		Use:   "prompt <PROMPT>",
		Short: "Prompt Gemini",
		Long:  "Prompt Gemini",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			prompt := args[0]
			projectNumber := string(injection.Resolve[config.GCPProjectNumber](ctx))

			// Initialize the Generative AI client
			c, err := genai.NewClient(ctx, projectNumber, "")
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer c.Close()

			// Ideally, we'd be able to go and find the models and use the lastet,
			// however there doesn't appear to be a way to list these types of models
			// yet.
			if model == "" {
				model = string(injection.Resolve[config.Model](ctx))
			}
			resp, err := c.GenerativeModel(model).GenerateContent(ctx, genai.Text(prompt))
			if err != nil {
				return fmt.Errorf("failed to generate content: %w", err)
			}

			if len(resp.Candidates) == 0 {
				return fmt.Errorf("failed to generate any candidates")
			}

			for _, part := range resp.Candidates[0].Content.Parts {
				fmt.Fprintln(cmd.OutOrStdout(), part)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&model, "model", "", "Set the model name to use (e.g., gemini-2.0-flash)")

	return cmd
}
