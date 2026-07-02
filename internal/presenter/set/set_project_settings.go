package set

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	_utils "github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
	"github.com/spf13/cobra"
)

type CmdSetProjectSettings struct {
	*cobra.Command
}

type UseCaseSetProjectSettings interface {
	Execute(ctx context.Context, rawAiproj []byte) error
}

func NewSetProjectSettingsCmd(uc UseCaseSetProjectSettings) CmdSetProjectSettings {
	var (
		filePath  string
		rawAiproj []byte
	)

	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Set project settings",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if filePath == "" {
				args = _utils.ReadArgsFromStdin(args)
				if len(args) == 0 {
					return fmt.Errorf("aiproj data required")
				}

				rawAiproj = []byte(args[0])
			} else {
				if !fshelper.PathExists(filePath) {
					return fmt.Errorf("file %s does not exist", filePath)
				}

				if !fshelper.IsFile(filePath) {
					return fmt.Errorf("path %s does not a file", filePath)
				}

				content, err := os.ReadFile(filePath)
				if err != nil {
					return fmt.Errorf("read aiproj file: %w", err)
				}

				rawAiproj = content
			}

			if !json.Valid(rawAiproj) {
				return fmt.Errorf("invalid aiproj data: not valid json")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, rawAiproj); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'set project settings' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "path to aiproj.json")

	return CmdSetProjectSettings{cmd}
}
