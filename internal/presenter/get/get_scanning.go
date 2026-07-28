package get

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

type CmdGetScanning struct {
	*cobra.Command
}

type UseCaseGetScanning interface {
	Execute(ctx context.Context, projectFilter *uuid.UUID) error
}

func NewGetScanningCmd(uc UseCaseGetScanning) CmdGetScanning {
	var projectIdFlag string

	cmd := &cobra.Command{
		Use:   "scanning",
		Short: "Get active scans",
		Long:  `List scans currently in progress. Optionally filter by project id; context project id is not used.`,
		Example: `  aictl get scanning
  aictl get scanning -p <project-id>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			var filter *uuid.UUID
			if projectIdFlag != "" {
				id, err := uuid.Parse(projectIdFlag)
				if err != nil {
					return validation.NewFieldError(projectIdFlag, "invalid uuid")
				}
				filter = &id
			}

			if err := uc.Execute(ctx, filter); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scanning' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectIdFlag, "project-id", "p", "", "filter by project id (ctx -p is ignored)")

	return CmdGetScanning{cmd}
}
