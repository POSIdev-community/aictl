package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
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
		Use:   "project <project-name>",
		Short: "Create a project",
		Long:  `Create a new AI project by name. The name may be passed as an argument or via stdin. With --safe, an existing project id is returned instead of an error.`,
		Example: `  aictl create project my-app
  aictl create project my-app --safe
  echo my-app | aictl create project`,
		Args: cobra.MaximumNArgs(1),
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
