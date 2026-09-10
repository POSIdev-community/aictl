package get

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
)

type PersistentPreRunEGetScanReportWithFiltersCmd _utils.RunE

type CmdGetScanReportWithFilters struct {
	*cobra.Command
}

type UseCaseGetScanReportWithFilters interface {
	Execute(ctx context.Context, scanId uuid.UUID, customReportName string, outPath string, includeComments, includeDFD, includeGlossary bool, l10n string, filters report.Filters) error
}

func NewPersistentPreRunEGetScanReportWithFiltersCmd(prev PersistentPreRunEGetScanReportCmd) PersistentPreRunEGetScanReportWithFiltersCmd {
	return _utils.ChainRunE(prev, func(cmd *cobra.Command, args []string) error {
		filters, err := buildReportFiltersFromFlags(cmd)
		if err != nil {
			return err
		}
		reportFilters = filters
		return nil
	})
}

func NewGetScanReportWithFiltersCmd(
	uc UseCaseGetScanReportWithFilters,
	defaultUC useCaseDefaultReportWithFilters,
	persistentPreRunE PersistentPreRunEGetScanReportWithFiltersCmd,
) CmdGetScanReportWithFilters {
	cmd := &cobra.Command{
		Use:   "with-filters <report-name> <scan-id>",
		Short: "Get scan report with vulnerability filters",
		Long: `Download a scan report with vulnerability filters applied (useFilters=true).
Requires at least one filter flag. Same format subcommands and output flags as get scan report.
Boolean filter flags send true only when present; unset flags are omitted from the API payload.
Array filters (--type, --language, --scan-module) that are not passed are sent as empty arrays.`,
		Example: `  aictl get scan report with-filters sarif <scan-id> -o ./out.sarif --level-high --level-medium
  aictl get scan report with-filters MyTemplate <scan-id> -o ./out.html --status-confirmed --language Java`,
		Args:              cobra.ExactArgs(2),
		PersistentPreRunE: persistentPreRunE,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, customReportName, outPath, includeComments, includeDFD, includeGlossary, l10n, reportFilters); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan report with-filters' usecase call: %w", err)
			}

			return nil
		},
	}

	// -o / -f / --include-* / --localization are inherited from parent `report` PersistentFlags.
	addReportFilterFlags(cmd)

	addDefaultReportFormatSubcommands(cmd, defaultUC, "get scan report with-filters")

	return CmdGetScanReportWithFilters{cmd}
}
