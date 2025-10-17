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
	var aiprojData aiproj.AIProj

	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Set project settings",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			args = _utils.ReadArgsFromStdin(args)
			if len(args) == 0 {
				return fmt.Errorf("aiproj data required")
			}

			var err error
			aiprojData, err = aiproj.FromString(args[0])
			if err != nil {
				return fmt.Errorf("invalid aiproj data: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			useCase, err := depsContainer.SetProjectSettingsUseCase(ctx, cfg)
			if err != nil {
				return fmt.Errorf("presenter set project settings useCase error: %w", err)
			}

			if err := useCase.Execute(ctx, cfg.ProjectId(), &aiprojData); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("presenter set project settings: %w", err)
			}

			return nil
		},
	}

	return cmd
}
