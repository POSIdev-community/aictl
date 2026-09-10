package scan

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/scantype"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
)

type CmdScanBranch struct {
	*cobra.Command
}

func NewScanBranchCmd(cfg *config.Config, uc UseCaseScanStartBranch) CmdScanBranch {
	var (
		projectIdFlag string
		scanLabel     string
		fullScan      bool
	)

	cmd := &cobra.Command{
		Use:   "branch <branch-id>",
		Short: "Start branch scan",
		Long:  `Start a scan on a branch. Branch id comes from the argument or context; project id from context or -p.`,
		Example: `  aictl scan branch <branch-id> -p <project-id>
  aictl scan branch --scan-label nightly --full-scan`,
		Args: cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cfg.UpdateProjectId(projectIdFlag); err != nil {
				return err
			}

			args = _utils.ReadArgsFromStdin(args)

			var branchIdFlag string
			if len(args) > 0 {
				branchIdFlag = args[0]
			}

			if err := cfg.UpdateBranchId(branchIdFlag); err != nil {
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

				return fmt.Errorf("'scan branch' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectIdFlag, "project-id", "p", "", "Project id (overrides context)")
	cmd.Flags().StringVar(&scanLabel, "scan-label", "", "Label for the scan (max 40 chars)")
	cmd.Flags().BoolVar(&fullScan, "full-scan", false, "Run a full scan instead of incremental")

	return CmdScanBranch{cmd}
}
