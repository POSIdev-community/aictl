package scan

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/scantype"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
)

type CmdScanProject struct {
	*cobra.Command
}

func NewScanProjectCmd(cfg *config.Config, uc UseCaseScanStartProject) CmdScanProject {
	var (
		scanLabel string
		fullScan  bool
	)

	cmd := &cobra.Command{
		Use:   "project <project-id>",
		Short: "Start project scan",
		Long:  `Start a scan on an entire project. Project id comes from the argument or context.`,
		Example: `  aictl scan project <project-id>
  aictl scan project --scan-label release --full-scan`,
		Args: cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			args = _utils.ReadArgsFromStdin(args)
			var projectIdFlag string
			if len(args) > 0 {
				projectIdFlag = args[0]
			}

			if err := cfg.UpdateProjectId(projectIdFlag); err != nil {
				return err
			}

			if scanLabel != "" {
				if len(scanLabel) > 40 {
					return fmt.Errorf("label length must be less than 40")
				}

				if strings.ContainsAny(scanLabel, "#%?/;,\"\r\n\\") {
					return fmt.Errorf("label contains invalid characters")
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			scanType := scantype.Incremental
			if fullScan {
				scanType = scantype.Full
			}

			if err := uc.Execute(ctx, scanLabel, scanType); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'scan project' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&scanLabel, "scan-label", "", "Label for the scan (max 40 chars)")
	cmd.Flags().BoolVar(&fullScan, "full-scan", false, "Run a full scan instead of incremental")

	return CmdScanProject{cmd}
}
