package get

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

type CmdGetProjectSettings struct {
	*cobra.Command
}

type UseCaseGetProjectSettings interface {
	Execute(ctx context.Context, jsonOutput bool) error
}

func NewGetProjectSettingsCmd(uc UseCaseGetProjectSettings) CmdGetProjectSettings {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Get project settings",
		Long:  `Print project settings. Project id comes from context or parent -p. Use --json for JSON output.`,
		Example: `  aictl get project settings -p <project-id>
  aictl get project settings -p <project-id> --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, jsonOutput); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get project settings' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	return CmdGetProjectSettings{cmd}
}
