package get

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/regexfilter"
)

type CmdGetProjects struct {
	*cobra.Command
}

type UseCaseGetProjects interface {
	Execute(ctx context.Context, filter regexfilter.RegexFilter, quiet bool) error
}

func NewGetProjectsCmd(uc UseCaseGetProjects) CmdGetProjects {

	var (
		filter string
		quiet  bool
	)

	var regexFilter regexfilter.RegexFilter

	cmd := &cobra.Command{
		Use:   "projects <regex>",
		Short: "Get AI projects",
		Long:  `List projects matching an optional regex filter. Filter may be passed as an argument or via stdin.`,
		Example: `  aictl get projects
  aictl get projects my-.*
  aictl get projects -q`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			filter = strings.Join(args, " ")

			regexFilter, err = regexfilter.NewRegexFilter(filter)
			if err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("new filter: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, regexFilter, quiet); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get projects' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Print only ids")

	return CmdGetProjects{cmd}
}
