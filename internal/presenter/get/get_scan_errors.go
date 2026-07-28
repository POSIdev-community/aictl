package get

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type CmdGetScanErrors struct {
	*cobra.Command
}

type UseCaseGetScanErrors interface {
	Execute(ctx context.Context, scanId uuid.UUID) error
}

func NewGetScanErrorsCmd(uc UseCaseGetScanErrors) CmdGetScanErrors {
	cmd := &cobra.Command{
		Use:     "errors <scan-id>",
		Short:   "Get scan errors",
		Long:    `Print scan errors. Scan id comes from argument or stdin.`,
		Example: `  aictl get scan errors <scan-id>`,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan errors' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetScanErrors{cmd}
}
