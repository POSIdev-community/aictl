package set

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	_utils "github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
)

type CmdSetProjectPolicies struct {
	*cobra.Command
}

type UseCaseSetProjectPolicies interface {
	Execute(ctx context.Context, policiesJSON []byte, check *bool) error
}

func NewSetProjectPoliciesCmd(uc UseCaseSetProjectPolicies) CmdSetProjectPolicies {
	var (
		filePath     string
		policiesJSON []byte
		check        *bool
	)

	cmd := &cobra.Command{
		Use:   "policies [json|-]",
		Short: "Set project security policies",
		Long: `Replace project security policies. Project id comes from context or -p.

Accepts the API SecurityPoliciesModel object, or a policies JSON array
(as in aisa --policy-settings-file); arrays are wrapped automatically.
When checkSecurityPoliciesAccordance is omitted, the current server value is preserved.
Comments in the policies file are preserved inside securityPolicies.`,
		Example: `  aictl set project policies -f policies.json -p <project-id>
  aictl set project policies -`,
		Args: cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			content, err := readPoliciesInput(filePath, args)
			if err != nil {
				return err
			}

			normalized, err := _utils.NormalizePoliciesJSON(content)
			if err != nil {
				return validation.NewMessageError("invalid policies data: " + err.Error())
			}

			policiesJSON = normalized.PoliciesJSON
			check = normalized.Check

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, policiesJSON, check); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'set project policies' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to policies JSON file, or - for stdin")

	return CmdSetProjectPolicies{cmd}
}

func readPoliciesInput(filePath string, args []string) ([]byte, error) {
	if filePath != "" {
		if filePath == "-" {
			return io.ReadAll(os.Stdin)
		}

		if !fshelper.PathExists(filePath) {
			return nil, validation.NewMessageError(fmt.Sprintf("file %s does not exist", filePath))
		}

		if !fshelper.IsFile(filePath) {
			return nil, validation.NewMessageError(fmt.Sprintf("path %s does not a file", filePath))
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, validation.NewMessageError(fmt.Sprintf("read policies file: %v", err))
		}

		return content, nil
	}

	args = _utils.ReadArgsFromStdin(args)
	if len(args) == 0 {
		return nil, validation.NewMessageError("policies JSON required")
	}

	return []byte(args[0]), nil
}
