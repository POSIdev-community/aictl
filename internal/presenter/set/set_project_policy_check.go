package set

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

type CmdSetProjectPolicyCheck struct {
	*cobra.Command
}

type UseCaseSetProjectPolicyCheck interface {
	Execute(ctx context.Context, enabled bool) error
}

func NewSetProjectPolicyCheckCmd(uc UseCaseSetProjectPolicyCheck) CmdSetProjectPolicyCheck {
	cmd := &cobra.Command{
		Use:   "policy-check <true|false>",
		Short: "Enable or disable project security policy check",
		Long: `Set checkSecurityPoliciesAccordance for the project (UI "use security policies" checkbox).
Existing policy rules are preserved. Project id comes from context or -p.`,
		Example: `  aictl set project policy-check true
  aictl set project policy-check false -p <project-id>`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			enabled, err := parseBoolArg(args[0])
			if err != nil {
				return validation.NewMessageError(err.Error())
			}

			ctx := cmd.Context()
			if err := uc.Execute(ctx, enabled); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'set project policy-check' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdSetProjectPolicyCheck{cmd}
}

func parseBoolArg(s string) (bool, error) {
	v, err := strconv.ParseBool(strings.TrimSpace(s))
	if err != nil {
		return false, fmt.Errorf("argument must be true or false")
	}

	return v, nil
}
