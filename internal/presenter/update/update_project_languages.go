package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

type CmdUpdateProjectLanguages struct {
	*cobra.Command
}

type UseCaseUpdateProjectLanguages interface {
	Execute(ctx context.Context) error
}

func NewUpdateProjectLanguagesCmd(uc UseCaseUpdateProjectLanguages) CmdUpdateProjectLanguages {
	cmd := &cobra.Command{
		Use:     "languages",
		Short:   "Update project languages",
		Long:    `Detect and update project languages from uploaded sources. Project id comes from context or -p.`,
		Example: `  aictl update project languages -p <project-id>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'update project languages' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdUpdateProjectLanguages{cmd}
}
