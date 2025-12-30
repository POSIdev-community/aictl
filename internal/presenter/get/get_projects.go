package get

import (
	"context"
	"fmt"
	"strings"

	"github.com/POSIdev-community/aictl/internal/core/domain/regexfilter"
	"github.com/spf13/cobra"
)

type CmdGetProjects struct {
	*cobra.Command
}

type UseCaseGetProjects interface {
	Execute(ctx context.Context, filter regexfilter.RegexFilter, quite bool) error
}

func NewGetProjectsCmd(uc UseCaseGetProjects) CmdGetProjects {

	var (
		filter string
		quite  bool
	)

	var regexFilter regexfilter.RegexFilter

	cmd := &cobra.Command{
		Use:   "projects <regex>",
		Short: "Get AI projects",
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

			if err := uc.Execute(ctx, regexFilter, quite); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get projects' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&quite, "quite", "q", false, "Get only ids")

	return CmdGetProjects{cmd}
}
