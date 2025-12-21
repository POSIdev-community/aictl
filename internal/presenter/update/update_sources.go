package update

import (
	"context"
	"fmt"
	"strings"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
	"github.com/spf13/cobra"
)

type CmdUpdateSources struct {
	*cobra.Command
}

type UseCaseUpdateSources interface {
	Execute(ctx context.Context, sourcePath string) error
}

func NewUpdateSourcesCmd(cfg *config.Config, uc UseCaseUpdateSources, cmdUpdateSourcesGit CmdUpdateSourcesGit) CmdUpdateSources {

	var (
		path string
	)

	cmd := &cobra.Command{
		Use:   "sources",
		Short: "Update sources",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {

			var err error

			if err = cfg.UpdateProjectId(projectIdFlag); err != nil {
				return err
			}

			if err = cfg.UpdateBranchId(branchIdFlag); err != nil {
				return err
			}

			path = strings.TrimSpace(args[0])
			if path == "" {
				return errs.NewValidationError("empty sources path")
			}

			if !fshelper.PathExists(path) {
				return errs.NewValidationError("path does not exist")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, path); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("get projects: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectIdFlag, "project-id", "p", "", "project id")
	cmd.Flags().StringVarP(&branchIdFlag, "branch-id", "b", "", "branch id")

	cmd.AddCommand(cmdUpdateSourcesGit.Command)

	return CmdUpdateSources{cmd}
}
