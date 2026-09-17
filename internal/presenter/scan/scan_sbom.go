package scan

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
)

type CmdScanSbom struct {
	*cobra.Command
}

type UseCaseScanSbom interface {
	Execute(ctx context.Context, scanLabel string) error
}

func NewScanSbomCmd(cfg *config.Config, uc UseCaseScanSbom) CmdScanSbom {
	var (
		projectIdFlag string
		scanLabel     string
	)

	cmd := &cobra.Command{
		Use:   "sbom <project-id>",
		Short: "Start SBOM scan",
		Long:  `Start a scan on an SBOM project. Project id comes from the argument, -p, or context.`,
		Example: `  aictl scan sbom <project-id>
  aictl scan sbom -p <project-id> --scan-label release`,
		Args: cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			args = _utils.ReadArgsFromStdin(args)
			id := projectIdFlag
			if len(args) > 0 {
				id = args[0]
			}

			if err := cfg.UpdateProjectId(id); err != nil {
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

			if err := uc.Execute(ctx, scanLabel); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'scan sbom' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectIdFlag, "project-id", "p", "", "Project id (overrides context)")
	cmd.Flags().StringVar(&scanLabel, "scan-label", "", "Label for the scan (max 40 chars)")

	return CmdScanSbom{cmd}
}
