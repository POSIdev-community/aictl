package scan

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/spf13/cobra"
)

type CmdScanStartProject struct {
	*cobra.Command
}

type UseCaseScanStartProject interface {
	Execute(ctx context.Context, scanLabel string) error
}

func NewScanStartProjectCmd(cfg *config.Config, uc UseCaseScanStartProject) CmdScanStartProject {
	cmd := &cobra.Command{
		Use:   "project <project-id>",
		Short: "Start project scan",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			args = _utils.ReadArgsFromStdin(args)
			var projectIdFlag string
			if len(args) > 0 {
				projectIdFlag = args[0]
			}

			if err := cfg.UpdateProjectId(projectIdFlag); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanLabel); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'scan start project' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdScanStartProject{cmd}
}
