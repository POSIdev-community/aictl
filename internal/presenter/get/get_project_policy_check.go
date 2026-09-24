package get

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

type CmdGetProjectPolicyCheck struct {
	*cobra.Command
}

type UseCaseGetProjectPolicyCheck interface {
	Execute(ctx context.Context) error
}

func NewGetProjectPolicyCheckCmd(uc UseCaseGetProjectPolicyCheck) CmdGetProjectPolicyCheck {
	cmd := &cobra.Command{
		Use:     "policy-check",
		Short:   "Get project security policy check flag",
		Long:    `Print checkSecurityPoliciesAccordance as true or false. Project id comes from context or parent -p.`,
		Example: `  aictl get project policy-check -p <project-id>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get project policy-check' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetProjectPolicyCheck{cmd}
}
