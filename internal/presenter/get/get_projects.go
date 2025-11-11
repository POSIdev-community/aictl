package get

import (
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/regexfilter"
	"github.com/spf13/cobra"
)

func NewGetProjectsCmd(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {

	var (
		filter string
		quite  bool
	)

	var regexFilter regexfilter.RegexFilter

	cmd := &cobra.Command{
		Use:   "projects",
		Short: "Get AI projects",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			regexFilter, err = regexfilter.NewRegexFilter(filter)
			if err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("new filter error: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			useCase, err := depsContainer.GetProjectsUseCase(ctx, cfg)
			if err != nil {
				return fmt.Errorf("get projects useCase error: %w", err)
			}

			if err := useCase.Execute(ctx, regexFilter, quite); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("get projects: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&filter, "name", "n", "", "Filter projects by name. Support regular expression")
	cmd.Flags().BoolVarP(&quite, "quite", "q", false, "Get only ids")

	return cmd
}
