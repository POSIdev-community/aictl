package set

import (
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	_utils "github.com/POSIdev-community/aictl/internal/presenter/.utils"
)

type PersistentPreRunESetProjectCmd _utils.RunE

type CmdSetProject struct {
	*cobra.Command
}

func NewPersistentPreRunESetProjectCmd(cfg *config.Config, prev PersistentPreRunESetCmd) PersistentPreRunESetProjectCmd {
	return _utils.ChainRunE(prev, func(cmd *cobra.Command, args []string) error {
		if err := cfg.UpdateProjectId(projectIdFlag); err != nil {
			return err
		}

		return nil
	})
}

var projectIdFlag string

func NewSetProjectCmd(persistentPreRunESetProjectCmd PersistentPreRunESetProjectCmd,
	setProjectSettingsCmd CmdSetProjectSettings, setProjectPoliciesCmd CmdSetProjectPolicies,
	setProjectPolicyCheckCmd CmdSetProjectPolicyCheck,
	setProjectExclusionsCmd CmdSetProjectExclusions) CmdSetProject {
	cmd := &cobra.Command{
		Use:               "project",
		Short:             "Set project configuration",
		Long:              `Set project-level configuration on the server. Project id comes from context or -p.`,
		PersistentPreRunE: persistentPreRunESetProjectCmd,
	}

	cmd.AddCommand(setProjectSettingsCmd.Command)
	cmd.AddCommand(setProjectPoliciesCmd.Command)
	cmd.AddCommand(setProjectPolicyCheckCmd.Command)
	cmd.AddCommand(setProjectExclusionsCmd.Command)

	cmd.PersistentFlags().StringVarP(&projectIdFlag, "project-id", "p", "", "Project id (overrides context)")

	return CmdSetProject{cmd}
}
