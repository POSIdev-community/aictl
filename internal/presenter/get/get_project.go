package get

import (
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	_utils "github.com/POSIdev-community/aictl/internal/presenter/.utils"
)

type PersistentPreRunEGetProjectCmd _utils.RunE

type CmdGetProject struct {
	*cobra.Command
}

func NewPersistentPreRunEGetProjectCmd(cfg *config.Config, prev PersistentPreRunEGetCmd) PersistentPreRunEGetProjectCmd {
	return _utils.ChainRunE(prev, func(cmd *cobra.Command, args []string) error {
		return cfg.UpdateProjectId(projectIdFlag)
	})
}

func NewGetProjectCmd(persistentPreRunE PersistentPreRunEGetProjectCmd, cmdGetProjectAiproj CmdGetProjectAiproj,
	cmdGetProjectSettings CmdGetProjectSettings, cmdGetProjectPolicies CmdGetProjectPolicies,
	cmdGetProjectPolicyCheck CmdGetProjectPolicyCheck,
	cmdGetProjectExclusions CmdGetProjectExclusions) CmdGetProject {
	cmd := &cobra.Command{
		Use:               "project",
		Short:             "Get project",
		Long:              `Retrieve project-level data. Project id comes from context or -p.`,
		PersistentPreRunE: persistentPreRunE,
	}

	cmd.AddCommand(cmdGetProjectAiproj.Command)
	cmd.AddCommand(cmdGetProjectSettings.Command)
	cmd.AddCommand(cmdGetProjectPolicies.Command)
	cmd.AddCommand(cmdGetProjectPolicyCheck.Command)
	cmd.AddCommand(cmdGetProjectExclusions.Command)

	cmd.PersistentFlags().StringVarP(&projectIdFlag, "project-id", "p", "", "Project id (overrides context)")

	return CmdGetProject{cmd}
}
