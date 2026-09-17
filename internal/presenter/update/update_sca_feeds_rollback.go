package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

type UseCaseRollbackScaFeeds interface {
	Execute(ctx context.Context, skipConfirm bool) error
}

func NewUpdateScaFeedsRollbackCmd(uc UseCaseRollbackScaFeeds) *cobra.Command {
	var skipConfirm bool

	cmd := &cobra.Command{
		Use:   "rollback",
		Short: "Roll back current SCA feeds package",
		Long:  `Roll back the current SCA feeds package to the previous version selected by the server (AIE ≥ 6.3). Prints the new version on stdout.`,
		Example: `  aictl update sca-feeds rollback
  aictl update sca-feeds rollback -y`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, skipConfirm); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'update sca-feeds rollback' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&skipConfirm, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}
