package context

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/pkg/logger"
)

type CmdConfigUnset struct {
	*cobra.Command
}

type UseCaseConfigUnset interface {
	Execute(uriUnset, tokenUnset, tlsUnset, cacertUnset, projectIdUnset, branchIdUnset bool) error
}

func NewConfigUnsetCommand(uc UseCaseConfigUnset) CmdConfigUnset {

	var (
		uriUnset       bool
		tokenUnset     bool
		tlsUnset       bool
		cacertUnset    bool
		projectIdUnset bool
		branchIdUnset  bool
	)

	cmd := &cobra.Command{
		Use:   "unset",
		Short: "Unset context parameters",
		Long:  `Clear selected fields from the local aictl context. At least one flag is required.`,
		Example: `  aictl ctx unset -p -b
  aictl ctx unset -u -t
  aictl ctx unset --cacert`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !uriUnset && !tokenUnset && !tlsUnset && !cacertUnset && !projectIdUnset && !branchIdUnset {
				return validation.NewError("Any configs not provided")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.FromContext(cmd.Context())
			log.StdErrf("aictl ctx")

			err := uc.Execute(uriUnset, tokenUnset, tlsUnset, cacertUnset, projectIdUnset, branchIdUnset)
			if err != nil {
				return fmt.Errorf("'ctx unset' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&uriUnset, "uri", "u", false, "Unset URI")
	cmd.Flags().BoolVarP(&tokenUnset, "token", "t", false, "Unset access token")
	cmd.Flags().BoolVar(&tlsUnset, "tls-skip", false, "Unset TLS skip setting")
	cmd.Flags().BoolVar(&cacertUnset, "cacert", false, "Unset CA certificate path")
	cmd.Flags().BoolVarP(&projectIdUnset, "project-id", "p", false, "Unset project id")
	cmd.Flags().BoolVarP(&branchIdUnset, "branch-id", "b", false, "Unset branch id")

	return CmdConfigUnset{cmd}
}
