package get

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/regexfilter"
)

type CmdGetScansSbomProject struct {
	*cobra.Command
}

type UseCaseGetScansSbomProject interface {
	Execute(ctx context.Context, filter regexfilter.RegexFilter, quiet bool, latest bool) error
}

func NewGetScansSbomProjectCmd(cfg *config.Config, uc UseCaseGetScansSbomProject) CmdGetScansSbomProject {
	var (
		projectIdFlag string
		quiet         bool
		latest        bool
		regexFilter   regexfilter.RegexFilter
	)

	cmd := &cobra.Command{
		Use:   "sbom-project [<regex>]",
		Short: "Get scans for an SBOM project",
		Long:  `List scans for an SBOM project matching an optional regex filter. Resolves the virtual branch from the project id (context or -p).`,
		Example: `  aictl get scans sbom-project -p <project-id>
  aictl get scans sbom-project 2024 -p <project-id>
  aictl get scans sbom-project --latest -p <project-id>`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cfg.UpdateProjectId(projectIdFlag); err != nil {
				return err
			}

			filter := strings.Join(args, " ")

			var err error
			regexFilter, err = regexfilter.NewRegexFilter(filter)
			if err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("new filter: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, regexFilter, quiet, latest); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scans sbom-project' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectIdFlag, "project-id", "p", "", "Project id (overrides context)")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Print only ids")
	cmd.Flags().BoolVar(&latest, "latest", false, "Return only the latest scan result")

	return CmdGetScansSbomProject{cmd}
}
