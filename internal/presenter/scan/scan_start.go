package scan

import (
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/spf13/cobra"
)

func NewScanStartCommand(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start scan",
	}

	cmd.AddCommand(NewScanStartBranchCommand(cfg, depsContainer))
	cmd.AddCommand(NewScanStartProjectCommand(cfg, depsContainer))

	return cmd
}
