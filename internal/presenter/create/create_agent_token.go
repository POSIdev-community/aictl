package create

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	_utils "github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/spf13/cobra"
)

type CmdCreateAgentToken struct {
	*cobra.Command
}

type UseCaseCreateAgentToken interface {
	Execute(ctx context.Context, login, password, agentName string) error
}

func NewCreateAgentTokenCmd(cfg *config.Config, uc UseCaseCreateAgentToken) CmdCreateAgentToken {
	var (
		login     string
		password  string
		agentName string
	)

	cmd := &cobra.Command{
		Use:   "agent-token <agent-name>",
		Short: "Create access token for scan agent",
		Long: `Create access token for AIE scan agent.

Requires admin credentials (login/password) to authenticate and create
a token with ScanAgent scope. The generated token can be used to
configure scan agents.

Example:
  aictl create agent-token my-agent --login admin --password secret -u https://aie-server:443`,
		Args:              cobra.ExactArgs(1),
		PersistentPreRunE: _utils.ChainRunE(_utils.InitializeLogger, _utils.UpdateConfig(cfg)),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if login == "" {
				return fmt.Errorf("login is required (use --login)")
			}
			if password == "" {
				return fmt.Errorf("password is required (use --password)")
			}

			agentName = args[0]

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, login, password, agentName); err != nil {
				cmd.SilenceUsage = true
				return fmt.Errorf("'create agent-token' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&login, "login", "", "admin user login (required)")
	cmd.Flags().StringVar(&password, "password", "", "admin user password (required)")

	return CmdCreateAgentToken{cmd}
}
