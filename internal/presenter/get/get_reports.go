package get

import (
	"fmt"
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func NewGetReportsCmd(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {

	var (
		sarif    bool
		plain    bool
		destPath string
	)

	var projectIds []uuid.UUID

	cmd := &cobra.Command{
		Use:   "reports",
		Short: "Get AI reports",
		Args:  cobra.MinimumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !sarif && !plain || sarif && plain {
				return fmt.Errorf("must specify only --sarif or --plain")
			}

			if destPath == "" {
				return fmt.Errorf("must specify --dest-path")
			}

			var err error

			args = _utils.ReadArgsFromStdin(args)
			projectIds, err = _utils.ParseUUIDs(args)
			if err != nil {
				return fmt.Errorf("get reports project ids parse error: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			useCase, err := depsContainer.GetReportsUseCase(ctx, cfg)
			if err != nil {
				return fmt.Errorf("presenter get reports useCase error: %w", err)
			}

			if err := useCase.Execute(ctx, projectIds, sarif, plain, destPath); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("presenter get reports: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&sarif, "sarif", false, "Get sarif report")
	cmd.Flags().BoolVar(&plain, "plain", false, "Get plaint report")
	cmd.Flags().StringVarP(&destPath, "dest-path", "d", ".", "Destination path")

	return cmd
}
