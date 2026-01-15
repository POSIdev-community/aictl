package get

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

type CmdGetVersion struct {
	*cobra.Command
}

type UseCaseGetVersion interface {
	Execute(ctx context.Context) error
}

func NewGetVersionCmd(uc UseCaseGetVersion) CmdGetVersion {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Get aie version",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get version' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetVersion{cmd}
}
