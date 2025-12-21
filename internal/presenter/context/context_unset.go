package context

import (
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/POSIdev-community/aictl/pkg/logger"
	"github.com/spf13/cobra"
)

type CmdConfigUnset struct {
	*cobra.Command
}

type UseCaseConfigUnset interface {
	Execute(uriUnset, tokenUnset, tlsUnset, projectIdUnset, branchIdUnset bool) error
}

func NewConfigUnsetCommand(uc UseCaseConfigUnset) CmdConfigUnset {

	var (
		uriUnset       bool
		tokenUnset     bool
		tlsUnset       bool
		projectIdUnset bool
		branchIdUnset  bool
	)

	cmd := &cobra.Command{
		Use:   "unset",
		Short: "Unset context params",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !uriUnset && !tokenUnset && !tlsUnset && !projectIdUnset && !branchIdUnset {
				return errs.NewValidationError("Any configs not provided")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.FromContext(cmd.Context())
			log.StdErrf("aictl ctx")

			err := uc.Execute(uriUnset, tokenUnset, tlsUnset, projectIdUnset, branchIdUnset)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&uriUnset, "uri", "u", false, "Unset uri")
	cmd.Flags().BoolVarP(&tokenUnset, "token", "t", false, "Unset token")
	cmd.Flags().BoolVar(&tlsUnset, "tls-skip", false, "Unset tls-skip")
	cmd.Flags().BoolVarP(&projectIdUnset, "project-id", "p", false, "Unset project id")
	cmd.Flags().BoolVarP(&branchIdUnset, "branch-id", "b", false, "Unset branch id")

	return CmdConfigUnset{cmd}
}
