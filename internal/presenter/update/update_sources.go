package update

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
	"github.com/POSIdev-community/aictl/pkg/gitignore"
)

type CmdUpdateSources struct {
	*cobra.Command
}

type UseCaseUpdateSources interface {
	Execute(ctx context.Context, sourcePath string, exclusions gitignore.Exclusions) error
}

func NewUpdateSourcesCmd(cfg *config.Config, uc UseCaseUpdateSources) CmdUpdateSources {

	var (
		path             string
		excludeFlags     []string
		excludeFromFlags []string
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
				return validation.NewError("empty sources path")
			}

			if !fshelper.PathExists(path) {
				return validation.NewError("path does not exist")
			}

			for _, excludeFrom := range excludeFromFlags {
				if !fshelper.PathExists(excludeFrom) {
					return validation.NewError(fmt.Sprintf("exclude-from path '%s' does not exist", excludeFrom))
				}
				if !fshelper.IsFile(excludeFrom) {
					return validation.NewError(fmt.Sprintf("exclude-from path '%s' is not a file", excludeFrom))
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			exclusions := gitignore.Exclusions{
				Patterns:  excludeFlags,
				FromFiles: excludeFromFlags,
			}

			if err := uc.Execute(ctx, path, exclusions); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'update sources' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectIdFlag, "project-id", "p", "", "project id")
	cmd.Flags().StringVarP(&branchIdFlag, "branch-id", "b", "", "branch id")
	cmd.Flags().StringArrayVarP(&excludeFlags, "exclude", "e", nil, "exclude file or directory (gitignore pattern)")
	cmd.Flags().StringArrayVar(&excludeFromFlags, "exclude-from", nil, "path to file with exclude patterns in gitignore format")

	return CmdUpdateSources{cmd}
}
