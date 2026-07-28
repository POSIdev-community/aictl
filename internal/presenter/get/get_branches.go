package get

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/regexfilter"
)

type CmdGetBranches struct {
	*cobra.Command
}

type UseCaseGetBranches interface {
	Execute(ctx context.Context, filter regexfilter.RegexFilter, quite bool) error
}

func NewGetBranchesCmd(cfg *config.Config, uc UseCaseGetBranches) CmdGetBranches {
	var (
		projectIdFlag string
		quite         bool
		regexFilter   regexfilter.RegexFilter
	)

	cmd := &cobra.Command{
		Use:   "branches [<regex>]",
		Short: "Get AI branches",
		Long:  `List branches for the current project matching an optional regex filter. Project id comes from context or -p.`,
		Example: `  aictl get branches -p <project-id>
  aictl get branches main -p <project-id>
  aictl get branches -p <project-id> -q`,
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

			if err := uc.Execute(ctx, regexFilter, quite); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get branches' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectIdFlag, "project-id", "p", "", "Project id (overrides context)")
	cmd.Flags().BoolVarP(&quite, "quite", "q", false, "Print only ids")

	return CmdGetBranches{cmd}
}
