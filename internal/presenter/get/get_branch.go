package get

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	_utils "github.com/POSIdev-community/aictl/internal/presenter/.utils"
)

type PersistentPreRunEGetBranchCmd _utils.RunE

type CmdGetBranch struct {
	*cobra.Command
}

func NewPersistentPreRunEGetBranchCmd(prev PersistentPreRunEGetCmd) PersistentPreRunEGetBranchCmd {
	return _utils.ChainRunE(prev, func(cmd *cobra.Command, args []string) error {
		args = _utils.ReadArgsFromStdin(args)
		if len(args) < 1 {
			return validation.NewError("missing branch id")
		}

		if len(args) > 1 {
			return validation.NewError("too many arguments")
		}

		var err error
		branchId, err = uuid.Parse(args[0])
		if err != nil {
			return validation.NewFieldError(args[0], "invalid uuid")
		}

		return nil
	})
}

type UseCaseGetBranch interface {
	Execute(ctx context.Context, branchId uuid.UUID) error
}

var branchId uuid.UUID

func NewGetBranchCmd(persistentPreRunE PersistentPreRunEGetBranchCmd, uc UseCaseGetBranch) CmdGetBranch {
	cmd := &cobra.Command{
		Use:               "branch <branch-id>",
		Short:             "Get branch",
		Long:              `Retrieve branch details by id. Branch id may be passed as an argument or via stdin.`,
		Example:           `  aictl get branch <branch-id>`,
		PersistentPreRunE: persistentPreRunE,
		Args:              cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, branchId); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get branch' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetBranch{cmd}
}
