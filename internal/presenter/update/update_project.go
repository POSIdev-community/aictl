package update

import (
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	_utils "github.com/POSIdev-community/aictl/internal/presenter/.utils"
)

type PersistentPreRunEUpdateProjectCmd _utils.RunE

type CmdUpdateProject struct {
	*cobra.Command
}

func NewPersistentPreRunEUpdateProjectCmd(cfg *config.Config, prev PersistentPreRunEUpdateCmd) PersistentPreRunEUpdateProjectCmd {
	return _utils.ChainRunE(prev, func(cmd *cobra.Command, args []string) error {
		return cfg.UpdateProjectId(projectIdFlag)
	})
}

func NewUpdateProjectCmd(
	persistentPreRunE PersistentPreRunEUpdateProjectCmd,
	cmdUpdateProjectSettings CmdUpdateProjectSettings,
	cmdUpdateProjectLanguages CmdUpdateProjectLanguages,
) CmdUpdateProject {
	cmd := &cobra.Command{
		Use:               "project",
		Short:             "Update project",
		Long:              `Update project settings and metadata on the server. Project id comes from context or -p.`,
		PersistentPreRunE: persistentPreRunE,
	}

	cmd.AddCommand(cmdUpdateProjectSettings.Command)
	cmd.AddCommand(cmdUpdateProjectLanguages.Command)

	cmd.PersistentFlags().StringVarP(&projectIdFlag, "project-id", "p", "", "Project id (overrides context)")

	return CmdUpdateProject{cmd}
}
