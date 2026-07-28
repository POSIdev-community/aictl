package get

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

type CmdGetProjectPolicies struct {
	*cobra.Command
}

type UseCaseGetProjectPolicies interface {
	Execute(ctx context.Context) error
}

func NewGetProjectPoliciesCmd(uc UseCaseGetProjectPolicies) CmdGetProjectPolicies {
	cmd := &cobra.Command{
		Use:     "policies",
		Short:   "Get project security policies",
		Long:    `Print project security policies. Project id comes from context or parent -p.`,
		Example: `  aictl get project policies -p <project-id>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get project policies' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetProjectPolicies{cmd}
}
