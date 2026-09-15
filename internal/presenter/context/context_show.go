package context

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

type CmdConfigShow struct {
	*cobra.Command
}

type UseCaseConfigShow interface {
	Execute(ctx context.Context, json bool, yaml bool) error
}

func NewConfigShowCommand(uc UseCaseConfigShow) CmdConfigShow {

	var (
		json bool
		yaml bool
	)

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show current aictl context",
		Long:  `Print the local aictl context. Default output is human-readable with a masked token; use --json or --yaml for machine-readable formats with the raw token (mutually exclusive). In JSON/YAML, unset fields are null instead of "<unset>".`,
		Example: `  aictl ctx show
  aictl ctx show --json
  aictl ctx show --yaml`,
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if json && yaml {
				return fmt.Errorf("cannot use both json and yaml flags")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			err := uc.Execute(ctx, json, yaml)
			if err != nil {
				return fmt.Errorf("'ctx show' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&json, "json", false, "Print context as JSON")
	cmd.Flags().BoolVar(&yaml, "yaml", false, "Print context as YAML")

	return CmdConfigShow{cmd}
}
