package set

import (
	"fmt"
	"os"

	_utils "github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/POSIdev-community/aictl/pkg/fshelper"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/aiproj"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
)

func NewSetProjectSettingsCmd(cfg *config.Config, depsContainer *application.DependenciesContainer) *cobra.Command {
	var (
		filePath   string
		aiprojData aiproj.AIProj
	)

	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Set project settings",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var aiprojString string
			if filePath == "" {
				args = _utils.ReadArgsFromStdin(args)
				if len(args) == 0 {
					return fmt.Errorf("aiproj data required")
				}

				aiprojString = args[0]
			} else {
				if !fshelper.PathExists(filePath) {
					return fmt.Errorf("file %s does not exist", filePath)
				}

				if !fshelper.IsFile(filePath) {
					return fmt.Errorf("path %s does not a file", filePath)
				}

				content, err := os.ReadFile(filePath)
				if err != nil {
					return fmt.Errorf("read aiproj file error: %w", err)
				}

				aiprojString = string(content)
			}

			var err error
			aiprojData, err = aiproj.FromString(aiprojString)
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

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "path to aiproj.json")

	return cmd
}
