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
) CmdUpdateProject {
	cmd := &cobra.Command{
		Use:               "project",
		Short:             "Update project",
		PersistentPreRunE: persistentPreRunE,
	}

	cmd.AddCommand(cmdUpdateProjectSettings.Command)

	cmd.PersistentFlags().StringVarP(&projectIdFlag, "project-id", "p", "", "project id")

	return CmdUpdateProject{cmd}
}
