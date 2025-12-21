package get

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/POSIdev-community/aictl/pkg/errs"
)

type CmdGetScan struct {
	*cobra.Command
}

type UseCaseGetScan interface {
	Execute(ctx context.Context, scanId uuid.UUID) error
}

var (
	projectIdFlag string
	scanId        uuid.UUID
)

func NewGetScanCmd(cfg *config.Config, uc UseCaseGetScan, cmdGetScanAiproj CmdGetScanAiproj,
	cmdGetScanLogs CmdGetScanLogs, cmdGetScanReport CmdGetScanReport, cmdGetScanSbom CmdGetScanSbom,
	cmdGetScanState CmdGetScanState) CmdGetScan {
	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Get scan",
		PersistentPreRunE: _utils.ConcatFuncs(_utils.InitializeLogger, func(cmd *cobra.Command, args []string) error {
			var err error
			if err = cfg.UpdateProjectId(projectIdFlag); err != nil {
				return err
			}

			args = _utils.ReadArgsFromStdin(args)
			if len(args) < 1 {
				return errs.NewValidationError("missing scan id")
			}

			scanIdFlag := args[0]

			scanId, err = uuid.Parse(scanIdFlag)
			if err != nil {
				return errs.NewValidationFieldError(scanIdFlag, "invalid uuid")
			}

			return nil
		}),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("get projects: %w", err)
			}

			return nil
		},
	}

	cmd.AddCommand(cmdGetScanAiproj.Command)
	cmd.AddCommand(cmdGetScanLogs.Command)
	cmd.AddCommand(cmdGetScanReport.Command)
	cmd.AddCommand(cmdGetScanSbom.Command)
	cmd.AddCommand(cmdGetScanState.Command)

	cmd.PersistentFlags().StringVarP(&projectIdFlag, "project-id", "p", "", "project id")

	return CmdGetScan{cmd}
}
