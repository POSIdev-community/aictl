package get

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
)

func NewGetScanReportGitlabCmd(cfg *config.Config, depsContainer *application.DependenciesContainer) *cobra.Command {
	cmd := &cobra.Command{
		Short: "Get scan report gitlab",
		Use:   "gitlab",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			useCase, err := depsContainer.GetScanReportGitlabUseCase(ctx, cfg)
			if err != nil {
				return fmt.Errorf("presenter get scan repot gitlab useCase error: %w", err)
			}

			if err := useCase.Execute(ctx, cfg, scanId, destPath, includeComments, includeDFD, includeGlossary); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("presenter get scan repot sarif: %w", err)
			}

			return nil
		},
	}

	return cmd
}
