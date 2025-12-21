package update

import (
	"fmt"

	"github.com/spf13/cobra"
)

type CmdUpdateSourcesGit struct {
	*cobra.Command
}

type UseCaseUpdateSourcesGit interface {
	Execute() error
}

func NewUpdateSourcesGitCmd(uc UseCaseUpdateSourcesGit) CmdUpdateSourcesGit {
	cmd := &cobra.Command{
		Use:   "git",
		Short: "Update sources git",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := uc.Execute(); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("get projects: %w", err)
			}

			return nil
		},
	}

	return CmdUpdateSourcesGit{cmd}
}
