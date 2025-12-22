package get

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type CmdGetScanState struct {
	*cobra.Command
}

type UseCaseGetScanState interface {
	Execute(ctx context.Context, scanId uuid.UUID) error
}

func NewGetScanStateCmd(uc UseCaseGetScanState) CmdGetScanState {
	cmd := &cobra.Command{
		Short: "Get scan stage",
		Use:   "stage",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan state' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetScanState{cmd}
}
