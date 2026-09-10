package update

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
)

type CmdUpdateSbom struct {
	*cobra.Command
}

type UseCaseUpdateSbom interface {
	Execute(ctx context.Context, sbomPath string) error
}

func NewUpdateSbomCmd(cfg *config.Config, uc UseCaseUpdateSbom) CmdUpdateSbom {
	var path string

	cmd := &cobra.Command{
		Use:   "sbom <path>",
		Short: "Update SBOM",
		Long:  `Upload or replace the SBOM file for an SBOM project. Project id comes from context or -p.`,
		Example: `  aictl update sbom ./sbom.json -p <project-id>
  aictl update sbom ./bom.xml`,
		Args: cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cfg.UpdateProjectId(projectIdFlag); err != nil {
				return err
			}

			path = strings.TrimSpace(args[0])
			if path == "" {
				return validation.NewError("empty sbom path")
			}

			if !fshelper.PathExists(path) {
				return validation.NewError("path does not exist")
			}

			if !fshelper.IsFile(path) {
				return validation.NewError("path is not a file")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, path); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'update sbom' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectIdFlag, "project-id", "p", "", "Project id (overrides context)")

	return CmdUpdateSbom{cmd}
}
