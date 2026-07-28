package get

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

type CmdGetProjectExclusions struct {
	*cobra.Command
}

type UseCaseGetProjectExclusions interface {
	Execute(ctx context.Context) error
}

func NewGetProjectExclusionsCmd(uc UseCaseGetProjectExclusions) CmdGetProjectExclusions {
	cmd := &cobra.Command{
		Use:     "exclusions",
		Short:   "Get project file/folder exclusions",
		Long:    `Print project file and folder exclusions (gitignore-style). Project id comes from context or parent -p.`,
		Example: `  aictl get project exclusions -p <project-id>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get project exclusions' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetProjectExclusions{cmd}
}
