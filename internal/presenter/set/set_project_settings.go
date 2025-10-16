package set

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/aiproj"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	_utils "github.com/POSIdev-community/aictl/internal/presenter/.utils"
)

func NewSetProjectSettingsCmd(cfg *config.Config, depsContainer *application.DependenciesContainer) *cobra.Command {
	cmd := &cobra.Command{
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			useCase, err := depsContainer.SetProjectSettingsUseCase(ctx, cfg)
			if err != nil {
				return fmt.Errorf("presenter set project settings useCase error: %w", err)
			}

			args = _utils.ReadArgsFromStdin(args)
			if len(args) == 0 {
				return fmt.Errorf("aiproj data required")
			}

			aiprojData, err := aiproj.FromString(args[0])
			if err != nil {
				return fmt.Errorf("invalid aiproj data: %w", err)
			}

			if err := useCase.Execute(ctx, cfg.ProjectId(), &aiprojData); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("presenter set project settings: %w", err)
			}

			return nil
		},
		Short: "Set project settings",
		Use:   "settings",
	}

	return cmd
}
