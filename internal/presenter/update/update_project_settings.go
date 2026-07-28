package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	domainsettings "github.com/POSIdev-community/aictl/internal/core/domain/settings"
)

type CmdUpdateProjectSettings struct {
	*cobra.Command
}

type UseCaseUpdateProjectSettings interface {
	Execute(ctx context.Context, patch domainsettings.ProjectSettingsPatch) error
}

func NewUpdateProjectSettingsCmd(uc UseCaseUpdateProjectSettings) CmdUpdateProjectSettings {
	var (
		priorityFlag              string
		agentsFlag                string
		preferredAgentsOnlyFlag   bool
		noPreferredAgentsOnlyFlag bool
	)

	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Update project settings",
		Long:  `Patch project scan settings (priority, preferred agents). At least one of --priority, --agents, --preferred-agents-only, or --no-preferred-agents-only is required. Project id comes from context or -p.`,
		Example: `  aictl update project settings --priority High -p <project-id>
  aictl update project settings --agents agent1,agent2 --preferred-agents-only`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("priority") &&
				!cmd.Flags().Changed("agents") &&
				!cmd.Flags().Changed("preferred-agents-only") &&
				!cmd.Flags().Changed("no-preferred-agents-only") {
				return fmt.Errorf("at least one of --priority, --agents, --preferred-agents-only or --no-preferred-agents-only is required")
			}

			if cmd.Flags().Changed("preferred-agents-only") && cmd.Flags().Changed("no-preferred-agents-only") {
				return fmt.Errorf("--preferred-agents-only and --no-preferred-agents-only cannot be used together")
			}

			if cmd.Flags().Changed("priority") {
				if _, err := domainsettings.ParsePriority(priorityFlag); err != nil {
					return err
				}
			}

			if cmd.Flags().Changed("agents") {
				if _, err := domainsettings.ParsePreferredAgentsCSV(agentsFlag); err != nil {
					return err
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			patch := domainsettings.ProjectSettingsPatch{}

			if cmd.Flags().Changed("priority") {
				p, err := domainsettings.ParsePriority(priorityFlag)
				if err != nil {
					cmd.SilenceUsage = true

					return err
				}
				patch.Priority = &p
			}

			if cmd.Flags().Changed("agents") {
				agents, err := domainsettings.ParsePreferredAgentsCSV(agentsFlag)
				if err != nil {
					cmd.SilenceUsage = true

					return err
				}
				patch.PreferredAgents = &agents
			}

			if cmd.Flags().Changed("no-preferred-agents-only") {
				preferredOnly := false
				patch.PreferredAgentsOnly = &preferredOnly
			} else if cmd.Flags().Changed("preferred-agents-only") {
				patch.PreferredAgentsOnly = &preferredAgentsOnlyFlag
			}

			if err := uc.Execute(ctx, patch); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'update project settings' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&priorityFlag, "priority", "", "Scan priority: None, Low, Medium, High, or Critical")
	cmd.Flags().StringVar(&agentsFlag, "agents", "", "Comma-separated preferred scan agent ids")
	cmd.Flags().BoolVar(&preferredAgentsOnlyFlag, "preferred-agents-only", false, "Use only the selected scan agents")
	cmd.Flags().BoolVar(&noPreferredAgentsOnlyFlag, "no-preferred-agents-only", false, "Allow all scan agents, not only selected ones")

	return CmdUpdateProjectSettings{cmd}
}
