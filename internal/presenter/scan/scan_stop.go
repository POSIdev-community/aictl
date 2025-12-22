package scan

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type CmdScanStop struct {
	*cobra.Command
}

type UseCaseScanStop interface {
	Execute(ctx context.Context, scanResultId uuid.UUID) error
}

func NewScanStopCmd(uc UseCaseScanStop) CmdScanStop {

	var scanId uuid.UUID

	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop scan",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			args = _utils.ReadArgsFromStdin(args)
			scanIdFlag := args[0]

			var err error
			scanId, err = uuid.Parse(scanIdFlag)
			if err != nil {
				return errs.NewValidationFieldError(scanIdFlag, "invalid uuid")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'scan stop' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdScanStop{cmd}
}
