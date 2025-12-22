package context

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

type CmdConfigClear struct {
	*cobra.Command
}

type UseCaseConfigClear interface {
	Execute(ctx context.Context, skipConfirm bool) error
}

func NewConfigClearCommand(uc UseCaseConfigClear) CmdConfigClear {

	var skipConfirm bool

	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear current aictl configuration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			err := uc.Execute(ctx, skipConfirm)
			if err != nil {
				return fmt.Errorf("'ctx clear' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&skipConfirm, "yes", "y", false, "Skip confirmation prompt")

	return CmdConfigClear{cmd}
}
