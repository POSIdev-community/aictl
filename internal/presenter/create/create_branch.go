package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	_utils "github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
	"github.com/POSIdev-community/aictl/pkg/gitignore"
)

type CmdCreateBranch struct {
	*cobra.Command
}

type UseCaseCreateBranch interface {
	Execute(ctx context.Context, cfg *config.Config, branchName, scanTarget string, safe bool, exclusions gitignore.Exclusions, tempDir string) error
}

func NewCreateBranchCmd(cfg *config.Config, uc UseCaseCreateBranch) CmdCreateBranch {
	var (
		projectIdFlag    string
		branchName       string
		scanTarget       string
		excludeFlags     []string
		excludeFromFlags []string
		tempDir          string
	)

	cmd := &cobra.Command{
		Use:   "branch <branch-name>",
		Short: "Create a branch",
		Long:  `Create a branch under a project. The name may be passed as an argument or via stdin. Optionally pack and upload sources from --scan-target with gitignore-style exclusions. Project id comes from context or -p.`,
		Example: `  aictl create branch main -p <project-id>
  aictl create branch main -p <project-id> -s ./src -e '*.tmp' --exclude-from .aictlignore
  aictl create branch main --safe
  echo main | aictl create branch - -p <project-id>`,
		Args: cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cfg.UpdateProjectId(projectIdFlag); err != nil {
				return err
			}

			if scanTarget != "" {
				if !fshelper.PathExists(scanTarget) {
					return validation.NewError(fmt.Sprintf("scan-target path '%s' not exists", scanTarget))
				}
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

			args = _utils.ReadArgsFromStdin(args)
			if len(args) < 1 || args[0] == "" {
				return validation.NewRequiredError("branch-name")
			}

			branchName = args[0]

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			exclusions := gitignore.Exclusions{
				Patterns:  excludeFlags,
				FromFiles: excludeFromFlags,
			}

			if err := uc.Execute(ctx, cfg, branchName, scanTarget, safeFlag, exclusions, tempDir); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'create branch' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectIdFlag, "project-id", "p", "", "Project id (overrides context)")
	cmd.Flags().StringVarP(&scanTarget, "scan-target", "s", "", "Path to sources to pack and upload")
	cmd.Flags().StringArrayVarP(&excludeFlags, "exclude", "e", nil, "Exclude path (gitignore pattern); repeatable")
	cmd.Flags().StringArrayVar(&excludeFromFlags, "exclude-from", nil, "File with gitignore-style exclude patterns")
	cmd.Flags().StringVar(&tempDir, "temp-dir", "", "Directory for temporary zip when packing sources")

	return CmdCreateBranch{cmd}
}
