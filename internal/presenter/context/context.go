package context

import (
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/spf13/cobra"
)

func NewContextCmd(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ctx",
		Short: "aictl context",
	}

	cmd.AddCommand(NewConfigClearCommand(depsContainer))
	cmd.AddCommand(NewConfigSetCommand(cfg, depsContainer))
	cmd.AddCommand(NewConfigShowCommand(cfg, depsContainer))
	cmd.AddCommand(NewConfigUnsetCommand(cfg, depsContainer))

	return cmd
}
