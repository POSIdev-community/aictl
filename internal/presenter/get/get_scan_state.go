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
	Execute(ctx context.Context, scanId uuid.UUID, failOnScanFailed bool) error
}

func NewGetScanStateCmd(uc UseCaseGetScanState) CmdGetScanState {
	var failOnScanFailed bool

	cmd := &cobra.Command{
		Use:   "stage <scan-id>",
		Short: "Get scan stage",
		Long:  `Print the current scan stage. Scan id comes from argument or stdin. Use --fail-on-scan-failed to exit with code 1 when stage is Failed or Aborted.`,
		Example: `  aictl get scan stage <scan-id>
  aictl get scan stage <scan-id> --fail-on-scan-failed`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, failOnScanFailed); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan stage' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&failOnScanFailed, "fail-on-scan-failed", false, "Exit with code 1 when scan stage is Failed or Aborted")

	return CmdGetScanState{cmd}
}
