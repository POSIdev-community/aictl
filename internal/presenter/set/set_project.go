package set

import (
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
)

type CmdSetProject struct {
	*cobra.Command
}

var projectIdFlag string

func NewSetProjectCmd(cfg *config.Config, setProjectSettingsCmd CmdSetProjectSettings) CmdSetProject {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Project",
		Long:  "Set project parameters",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cfg.UpdateProjectId(projectIdFlag); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.AddCommand(setProjectSettingsCmd.Command)

	cmd.PersistentFlags().StringVarP(&projectIdFlag, "project-id", "p", "", "project id")

	return CmdSetProject{cmd}
}
