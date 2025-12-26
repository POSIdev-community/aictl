package delete

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type CmdDeleteProjects struct {
	*cobra.Command
}

type UseCaseDeleteProjects interface {
	Execute(ctx context.Context, projectIds []uuid.UUID) error
}

func NewDeleteProjectsCommand(uc UseCaseDeleteProjects) CmdDeleteProjects {

	var projectIds []uuid.UUID

	cmd := &cobra.Command{
		Use:   "projects <project-id>...",
		Short: "Delete AI projects",
		Args:  cobra.MinimumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			args = _utils.ReadArgsFromStdin(args)
			projectIds, err = _utils.ParseUUIDs(args)
			if err != nil {
				return fmt.Errorf("project ids parse: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, projectIds); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'delete projects' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdDeleteProjects{cmd}
}
