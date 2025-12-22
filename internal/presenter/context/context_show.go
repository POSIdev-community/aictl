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
		Args:  cobra.NoArgs,
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

	cmd.Flags().BoolVar(&json, "json", false, "Json format context")
	cmd.Flags().BoolVar(&yaml, "yaml", false, "Yaml format context")

	return CmdConfigShow{cmd}
}
