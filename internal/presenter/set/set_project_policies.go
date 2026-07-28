package set

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	_utils "github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
)

type CmdSetProjectPolicies struct {
	*cobra.Command
}

type UseCaseSetProjectPolicies interface {
	Execute(ctx context.Context, rawJSON []byte) error
}

func NewSetProjectPoliciesCmd(uc UseCaseSetProjectPolicies) CmdSetProjectPolicies {
	var (
		filePath string
		rawJSON  []byte
	)

	cmd := &cobra.Command{
		Use:   "policies [json|-]",
		Short: "Set project security policies",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			rawJSON, err = readPoliciesInput(filePath, args)
			if err != nil {
				return err
			}

			if !json.Valid(rawJSON) {
				return fmt.Errorf("invalid policies data: not valid json")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, rawJSON); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'set project policies' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "path to policies JSON file, or - for stdin")

	return CmdSetProjectPolicies{cmd}
}

func readPoliciesInput(filePath string, args []string) ([]byte, error) {
	if filePath != "" {
		if filePath == "-" {
			return io.ReadAll(os.Stdin)
		}

		if !fshelper.PathExists(filePath) {
			return nil, fmt.Errorf("file %s does not exist", filePath)
		}

		if !fshelper.IsFile(filePath) {
			return nil, fmt.Errorf("path %s does not a file", filePath)
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("read policies file: %w", err)
		}

		return content, nil
	}

	args = _utils.ReadArgsFromStdin(args)
	if len(args) == 0 {
		return nil, fmt.Errorf("policies JSON required")
	}

	return []byte(args[0]), nil
}
