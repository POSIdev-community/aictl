package get

import (
	"fmt"
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func NewGetScansCmd(cfg *config.Config, depsContainer *application.DependenciesContainer) *cobra.Command {
	var branchId uuid.UUID

	cmd := &cobra.Command{
		Short: "Get scans",
		Use:   "scans",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			args = _utils.ReadArgsFromStdin(args)
			branchIdFlag := args[0]

			var err error
			branchId, err = uuid.Parse(branchIdFlag)
			if err != nil {
				return errs.NewValidationFieldError(branchIdFlag, "invalid uuid")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			useCase, err := depsContainer.GetScansUseCase(ctx, cfg)
			if err != nil {
				return fmt.Errorf("presenter get reports useCase error: %w", err)
			}

			if err := useCase.Execute(ctx, branchId); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("presenter get reports: %w", err)
			}

			return nil
		},
	}

	return cmd
}
