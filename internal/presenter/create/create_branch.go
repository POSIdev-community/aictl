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
	Execute(ctx context.Context, cfg *config.Config, branchName, scanTarget string, safe bool, exclusions gitignore.Exclusions) error
}

func NewCreateBranchCmd(cfg *config.Config, uc UseCaseCreateBranch) CmdCreateBranch {
	var (
		projectIdFlag    string
		branchName       string
		scanTarget       string
		excludeFlags     []string
		excludeFromFlags []string
	)

	cmd := &cobra.Command{
		Use:   "branch <branch-name>",
		Short: "Create branch",
		Args:  cobra.MaximumNArgs(1),
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

			args = _utils.ReadArgsFromStdin(args)
			branchName = args[0]

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			exclusions := gitignore.Exclusions{
				Patterns:  excludeFlags,
				FromFiles: excludeFromFlags,
			}

			if err := uc.Execute(ctx, cfg, branchName, scanTarget, safeFlag, exclusions); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'create branch' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectIdFlag, "project-id", "p", "", "project id")
	cmd.Flags().StringVarP(&scanTarget, "scan-target", "s", "", "scan target path")
	cmd.Flags().StringArrayVarP(&excludeFlags, "exclude", "e", nil, "exclude file or directory (gitignore pattern)")
	cmd.Flags().StringArrayVar(&excludeFromFlags, "exclude-from", nil, "path to file with exclude patterns in gitignore format")

	return CmdCreateBranch{cmd}
}
