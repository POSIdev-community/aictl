package create

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/spf13/cobra"
)

type CmdCreateProject struct {
	*cobra.Command
}

type UseCaseCreateProject interface {
	Execute(ctx context.Context, projectName string, safe bool) error
}

func NewCreateProjectCmd(uc UseCaseCreateProject) CmdCreateProject {

	var (
		projectName string
	)

	cmd := &cobra.Command{
		Use:   "project",
		Short: "Create project",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			args = _utils.ReadArgsFromStdin(args)
			projectName = args[0]

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, projectName, safeFlag); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'create project' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdCreateProject{cmd}
}
