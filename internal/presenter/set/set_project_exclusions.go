package set

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	_utils "github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
)

type CmdSetProjectExclusions struct {
	*cobra.Command
}

type UseCaseSetProjectExclusions interface {
	Execute(ctx context.Context, exclusions string) error
}

func NewSetProjectExclusionsCmd(uc UseCaseSetProjectExclusions) CmdSetProjectExclusions {
	var (
		filePath   string
		exclusions string
	)

	cmd := &cobra.Command{
		Use:   "exclusions [text|-]",
		Short: "Set project exclusions",
		Long:  `Replace project file/folder exclusions with gitignore-like text from an argument, file, or stdin. Existing exclusions are fully replaced.`,
		Example: `  aictl set project exclusions -f .aictlignore -p <project-id>
  aictl set project exclusions -
  aictl set project exclusions '*.tmp' -p <project-id>`,
		Args: cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			exclusions, err = readExclusionsInput(filePath, args)
			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, exclusions); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'set project exclusions' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to exclusions file, or - for stdin")

	return CmdSetProjectExclusions{cmd}
}

func readExclusionsInput(filePath string, args []string) (string, error) {
	if filePath != "" {
		if filePath == "-" {
			content, err := io.ReadAll(os.Stdin)
			if err != nil {
				return "", fmt.Errorf("read exclusions from stdin: %w", err)
			}

			return string(content), nil
		}

		if !fshelper.PathExists(filePath) {
			return "", fmt.Errorf("file %s does not exist", filePath)
		}

		if !fshelper.IsFile(filePath) {
			return "", fmt.Errorf("path %s does not a file", filePath)
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("read exclusions file: %w", err)
		}

		return string(content), nil
	}

	args = _utils.ReadArgsFromStdin(args)
	if len(args) == 0 {
		return "", fmt.Errorf("exclusions text required")
	}

	return args[0], nil
}
