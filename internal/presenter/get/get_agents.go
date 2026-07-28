package get

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

type CmdGetAgents struct {
	*cobra.Command
}

type UseCaseGetAgents interface {
	Execute(ctx context.Context, quite bool) error
}

func NewGetAgentsCmd(uc UseCaseGetAgents) CmdGetAgents {
	var quite bool

	cmd := &cobra.Command{
		Use:   "agents",
		Short: "Get agents",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, quite); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get agents' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&quite, "quite", "q", false, "Get only ids")

	return CmdGetAgents{cmd}
}
