package scan

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
)

type CmdScanAwait struct {
	*cobra.Command
}

type UseCaseScanAwait interface {
	Execute(ctx context.Context, scanId uuid.UUID, failOnScanFailed bool) error
}

func NewScanAwaitCmd(cfg *config.Config, uc UseCaseScanAwait) CmdScanAwait {
	var (
		projectIdFlag    string
		scanIdFlag       string
		scanId           uuid.UUID
		failOnScanFailed bool
	)

	cmd := &cobra.Command{
		Use:   "await <scan-id>",
		Short: "Await scan completion",
		Long:  `Wait until a scan reaches a terminal stage. Project id comes from context or -p.`,
		Example: `  aictl scan await <scan-id>
  aictl scan await <scan-id> -p <project-id> --fail-on-scan-failed`,
		Args: cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if err = cfg.UpdateProjectId(projectIdFlag); err != nil {
				return err
			}

			args = _utils.ReadArgsFromStdin(args)
			if len(args) < 1 {
				return validation.NewError("missing scan id")
			}

			scanIdFlag = args[0]
			scanId, err = uuid.Parse(scanIdFlag)
			if err != nil {
				return validation.NewFieldError(scanIdFlag, "invalid uuid")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, failOnScanFailed); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'scan await' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectIdFlag, "project-id", "p", "", "Project id (overrides context)")
	cmd.Flags().BoolVar(&failOnScanFailed, "fail-on-scan-failed", false, "Exit 1 if scan stage is Failed or Aborted")

	return CmdScanAwait{cmd}
}
