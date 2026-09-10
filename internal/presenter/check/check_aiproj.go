package check

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	_utils "github.com/POSIdev-community/aictl/internal/presenter/.utils"
)

type CmdCheckAiproj struct {
	*cobra.Command
}

type UseCaseCheckAiproj interface {
	Execute(ctx context.Context, raw []byte, schemaVersion string, jsonOut bool) error
}

func NewCheckAiprojCmd(uc UseCaseCheckAiproj) CmdCheckAiproj {
	var (
		filePath      string
		schemaVersion string
		jsonOut       bool
		rawAiproj     []byte
	)

	cmd := &cobra.Command{
		Use:   "aiproj",
		Short: "Check aiproj schema",
		Long:  `Validate an aiproj/JSON document against JSON Schema offline. Input via -f, positional JSON, or stdin (same as set project settings).`,
		Example: `  aictl check aiproj -f aiproj.json
  aictl check aiproj -f aiproj.json --schema-version 1.11
  aictl check aiproj -f aiproj.json --json
  aictl check aiproj '{"Version":"1.11","ProjectName":"demo","ProgrammingLanguages":["Go"],"ScanModules":["StaticCodeAnalysis"]}'`,
		Args: cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			content, err := _utils.ReadAiprojInput(filePath, args)
			if err != nil {
				return err
			}

			rawAiproj = content

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, rawAiproj, schemaVersion, jsonOut); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'check aiproj' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to aiproj.json file")
	cmd.Flags().StringVar(&schemaVersion, "schema-version", "", "Schema version to validate against (1.8, 1.9, 1.10, 1.11); default: auto-detect")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "Print result as pretty JSON on stdout")

	return CmdCheckAiproj{cmd}
}
