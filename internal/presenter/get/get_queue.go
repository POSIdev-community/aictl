package get

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

type CmdGetQueue struct {
	*cobra.Command
}

type UseCaseGetQueue interface {
	Execute(ctx context.Context, projectFilter *uuid.UUID) error
}

func NewGetQueueCmd(uc UseCaseGetQueue) CmdGetQueue {
	var projectIdFlag string

	cmd := &cobra.Command{
		Use:   "queue",
		Short: "Get scan queue",
		Long:  `List scans waiting in the queue. Optionally filter by project id; context project id is not used.`,
		Example: `  aictl get queue
  aictl get queue -p <project-id>`,
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

				return fmt.Errorf("'get queue' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectIdFlag, "project-id", "p", "", "filter by project id (ctx -p is ignored)")

	return CmdGetQueue{cmd}
}
