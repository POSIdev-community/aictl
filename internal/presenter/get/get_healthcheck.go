package get

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

type CmdGetHealthcheck struct {
	*cobra.Command
}

type UseCaseGetHealthcheck interface {
	Execute(ctx context.Context) error
}

func NewGetHealthcheckCmd(uc UseCaseGetHealthcheck) CmdGetHealthcheck {
	cmd := &cobra.Command{
		Use:   "healthcheck",
		Short: "Get aie healthcheck",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get healthcheck' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetHealthcheck{cmd}
}
