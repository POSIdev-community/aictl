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
		login    string
		password string
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
		Args: cobra.ExactArgs(1),
		// Override parent's PersistentPreRunE - agent-token uses login/password, not API token
		PersistentPreRunE: _utils.ChainRunE(
			_utils.InitializeLogger,
			func(cmd *cobra.Command, args []string) error {
				// Only update connection config without token validation
				return _utils.UpdateConnectionConfig(cfg)
			},
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			agentName := args[0]

			if login == "" {
				return fmt.Errorf("login is required (use --login)")
			}
			if password == "" {
				return fmt.Errorf("password is required (use --password)")
			}

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
