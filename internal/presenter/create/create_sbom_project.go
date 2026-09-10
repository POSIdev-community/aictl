package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
)

type CmdCreateSbomProject struct {
	*cobra.Command
}

type UseCaseCreateSbomProject interface {
	Execute(ctx context.Context, projectName, sbomPath string, safe bool) error
}

func NewCreateSbomProjectCmd(uc UseCaseCreateSbomProject) CmdCreateSbomProject {
	var (
		projectName string
		filePath    string
	)

	cmd := &cobra.Command{
		Use:   "sbom-project <project-name>",
		Short: "Create an SBOM project",
		Long:  `Create a new SBOM project and upload an SBOM file. The name may be passed as an argument or via stdin. With --safe, an existing SBOM project is updated instead of an error.`,
		Example: `  aictl create sbom-project my-sbom --file ./sbom.json
  aictl create sbom-project my-sbom --file ./sbom.json --safe
  echo my-sbom | aictl create sbom-project --file ./sbom.json`,
		Args: cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			args = _utils.ReadArgsFromStdin(args)
			projectName = args[0]

			if filePath == "" {
				return validation.NewError("empty sbom file path")
			}

			if !fshelper.PathExists(filePath) {
				return validation.NewError("path does not exist")
			}

			if !fshelper.IsFile(filePath) {
				return validation.NewError("path is not a file")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, projectName, filePath, safeFlag); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'create sbom-project' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&filePath, "file", "", "Path to SBOM file to upload")
	_ = cmd.MarkFlagRequired("file")

	return CmdCreateSbomProject{cmd}
}
