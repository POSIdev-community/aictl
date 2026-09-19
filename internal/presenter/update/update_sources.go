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
	Execute(ctx context.Context, sourcePath string, exclusions gitignore.Exclusions, tempDir string) error
}

func NewUpdateSourcesCmd(cfg *config.Config, uc UseCaseUpdateSources, ucLanguages UseCaseUpdateProjectLanguages) CmdUpdateSources {

	var (
		path             string
		excludeFlags     []string
		excludeFromFlags []string
		tempDir          string
		updateLanguages  bool
	)

	cmd := &cobra.Command{
		Use:   "sources <path>",
		Short: "Update sources",
		Long:  `Upload or update sources for a project/branch with optional gitignore-style exclusions. Optionally recalculate project languages (--update-languages). Project and branch ids come from context or -p/-b.`,
		Example: `  aictl update sources ./src -p <project-id> -b <branch-id>
  aictl update sources ./src --update-languages
  aictl update sources ./src -e '*.tmp' --exclude-from .aictlignore`,
		Args: cobra.ExactArgs(1),
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

			if tempDir != "" {
				if !fshelper.PathExists(tempDir) {
					return validation.NewError(fmt.Sprintf("temp-dir path '%s' does not exist", tempDir))
				}
				if !fshelper.IsDirectory(tempDir) {
					return validation.NewError(fmt.Sprintf("temp-dir path '%s' is not a directory", tempDir))
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

			if err := uc.Execute(ctx, path, exclusions, tempDir); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'update sources' usecase call: %w", err)
			}

			if updateLanguages {
				if err := ucLanguages.Execute(ctx); err != nil {
					cmd.SilenceUsage = true

					return fmt.Errorf("'update project languages' usecase call: %w", err)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectIdFlag, "project-id", "p", "", "Project id (overrides context)")
	cmd.Flags().StringVarP(&branchIdFlag, "branch-id", "b", "", "Branch id (overrides context)")
	cmd.Flags().StringArrayVarP(&excludeFlags, "exclude", "e", nil, "Exclude path (gitignore pattern); repeatable")
	cmd.Flags().StringArrayVar(&excludeFromFlags, "exclude-from", nil, "File with gitignore-style exclude patterns")
	cmd.Flags().StringVar(&tempDir, "temp-dir", "", "Directory for temporary zip when packing sources")
	cmd.Flags().BoolVar(&updateLanguages, "update-languages", false, "Recalculate project languages after upload")

	return CmdUpdateSources{cmd}
}
