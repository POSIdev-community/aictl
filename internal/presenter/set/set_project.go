package set

import (
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
)

func NewSetProjectCmd(cfg *config.Config, depsContainer *application.DependenciesContainer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Project",
		Long:  "Set project parameters",
	}

	cmd.AddCommand(NewSetProjectSettingsCmd(cfg, depsContainer))

	return cmd
}
