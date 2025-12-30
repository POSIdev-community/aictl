package scan

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type CmdScanAwait struct {
	*cobra.Command
}

type UseCaseScanAwait interface {
	Execute(ctx context.Context, scanId uuid.UUID) error
}

func NewScanAwaitCmd(cfg *config.Config, uc UseCaseScanAwait) CmdScanAwait {
	var (
		projectIdFlag string
		scanIdFlag    string
		scanId        uuid.UUID
	)

	cmd := &cobra.Command{
		Use:   "await <scan-id>",
		Short: "Await scan",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if err = cfg.UpdateProjectId(projectIdFlag); err != nil {
				return err
			}

			args = _utils.ReadArgsFromStdin(args)
			scanIdFlag = args[0]

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

				return fmt.Errorf("'scan await' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectIdFlag, "project-id", "p", "", "project id")

	return CmdScanAwait{cmd}
}
