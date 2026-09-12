package get

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/scafeeds"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
)

type CmdGetScaFeeds struct {
	*cobra.Command
}

type UseCaseGetScaFeeds interface {
	Execute(ctx context.Context, statuses []scafeeds.Status) error
}

type UseCaseDownloadScaFeeds interface {
	Execute(ctx context.Context, version, outPath string) error
}

func NewGetScaFeedsCmd(listUC UseCaseGetScaFeeds, downloadUC UseCaseDownloadScaFeeds) CmdGetScaFeeds {
	var (
		statuses []string
		outPath  string
	)

	cmd := &cobra.Command{
		Use:   "sca-feeds [version]",
		Short: "List or download SCA feeds packages",
		Long: `List SCA feeds packages or download one by version (AIE ≥ 6.3).

Without arguments — print table. With <version> — download archive.
Optional -o: new file path, or existing directory (writes default file name inside).`,
		Example: `  aictl get sca-feeds
  aictl get sca-feeds --status current --status active
  aictl get sca-feeds 1.2.3
  aictl get sca-feeds 1.2.3 -o ./feeds.zip
  aictl get sca-feeds 1.2.3 -o ./outdir/`,
		Args: cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			hasVersion := len(args) == 1
			if hasVersion && len(statuses) > 0 {
				return validation.NewError("--status cannot be used with version download")
			}
			if !hasVersion && outPath != "" {
				return validation.NewError("-o requires version argument")
			}

			for _, s := range statuses {
				if !scafeeds.Status(s).Valid() {
					return validation.NewError(fmt.Sprintf("unknown status '%s'", s))
				}
			}

			if outPath != "" && fshelper.PathExists(outPath) && fshelper.IsFile(outPath) {
				return validation.NewError("'output' path exists")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if len(args) == 1 {
				version := strings.TrimSpace(args[0])
				if version == "" {
					return validation.NewError("empty version")
				}

				if err := downloadUC.Execute(ctx, version, outPath); err != nil {
					cmd.SilenceUsage = true

					return fmt.Errorf("'get sca-feeds' download usecase call: %w", err)
				}

				return nil
			}

			parsed := make([]scafeeds.Status, 0, len(statuses))
			for _, s := range statuses {
				parsed = append(parsed, scafeeds.Status(s))
			}

			if err := listUC.Execute(ctx, parsed); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get sca-feeds' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringArrayVar(&statuses, "status", nil, "Filter by status (repeatable): active, current, archived, rolled_back")
	cmd.Flags().StringVarP(&outPath, "output", "o", "", "Output file path or existing directory")

	return CmdGetScaFeeds{cmd}
}
