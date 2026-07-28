package get

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
)

type PersistentPreRunEGetScanCmd _utils.RunE

type CmdGetScan struct {
	*cobra.Command
}

func NewPersistentPreRunEGetScanCmd(cfg *config.Config, prev PersistentPreRunEGetCmd) PersistentPreRunEGetScanCmd {
	return _utils.ChainRunE(prev, func(cmd *cobra.Command, args []string) error {
		var err error
		if err = cfg.UpdateProjectId(projectIdFlag); err != nil {
			return err
		}

		args = _utils.ReadArgsFromStdin(args)
		if len(args) < 1 {
			return validation.NewError("missing scan id")
		}

		if len(args) > 2 {
			return validation.NewError("too many arguments")
		}

		scanIdFlag := args[len(args)-1]
		customReportName = args[0]

		if customReportName == "" {
			return validation.NewFieldError("custom-report-name", "cannot be empty")
		}

		scanId, err = uuid.Parse(scanIdFlag)
		if err != nil {
			return validation.NewFieldError(scanIdFlag, "invalid uuid")
		}

		return nil
	})
}

type UseCaseGetScan interface {
	Execute(ctx context.Context, scanId uuid.UUID) error
}

var (
	projectIdFlag    string
	scanId           uuid.UUID
	customReportName string
)

func NewGetScanCmd(persistentPreRunE PersistentPreRunEGetScanCmd, uc UseCaseGetScan, cmdGetScanAiproj CmdGetScanAiproj,
	cmdGetScanLogs CmdGetScanLogs, cmdGetScanErrors CmdGetScanErrors, cmdGetScanReport CmdGetScanReport,
	cmdGetScanSbom CmdGetScanSbom, cmdGetScanState CmdGetScanState, cmdGetScanStatistic CmdGetScanStatistic) CmdGetScan {

	cmd := &cobra.Command{
		Use:               "scan <scan-id>",
		Short:             "Get scan",
		Long:              `Retrieve scan-level data. Scan id may be passed as an argument or via stdin. Project id comes from context or -p.`,
		PersistentPreRunE: persistentPreRunE,
		Args:              cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.AddCommand(cmdGetScanAiproj.Command)
	cmd.AddCommand(cmdGetScanLogs.Command)
	cmd.AddCommand(cmdGetScanErrors.Command)
	cmd.AddCommand(cmdGetScanReport.Command)
	cmd.AddCommand(cmdGetScanSbom.Command)
	cmd.AddCommand(cmdGetScanState.Command)
	cmd.AddCommand(cmdGetScanStatistic.Command)

	cmd.PersistentFlags().StringVarP(&projectIdFlag, "project-id", "p", "", "Project id (overrides context)")

	return CmdGetScan{cmd}
}
