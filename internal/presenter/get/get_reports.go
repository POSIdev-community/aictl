package get

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type CmdGetReports struct {
	*cobra.Command
}

type UseCaseGetReports interface {
	Execute(ctx context.Context, projectIds []uuid.UUID, sarif bool, plain bool, destPath string, includeComments, includeDFD, includeGlossary bool) error
}

func NewGetReportsCmd(uc UseCaseGetReports) CmdGetReports {

	var (
		sarif    bool
		plain    bool
		destPath string
	)

	var projectIds []uuid.UUID

	cmd := &cobra.Command{
		Use:   "reports",
		Short: "Get AI reports",
		Args:  cobra.MinimumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !sarif && !plain || sarif && plain {
				return fmt.Errorf("must specify only --sarif or --plain")
			}

			if destPath == "" {
				return fmt.Errorf("must specify --dest-path")
			}

			var err error

			args = _utils.ReadArgsFromStdin(args)
			projectIds, err = _utils.ParseUUIDs(args)
			if err != nil {
				return fmt.Errorf("project ids parse: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, projectIds, sarif, plain, destPath, includeComments, includeDFD, includeGlossary); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get reports' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&sarif, "sarif", false, "Get sarif report")
	cmd.Flags().BoolVar(&plain, "plain", false, "Get plaint report")
	cmd.Flags().StringVarP(&destPath, "dest-path", "d", ".", "Destination path")
	cmd.Flags().BoolVarP(&includeComments, "include-comments", "", false, "Include comments in the report file")
	cmd.Flags().BoolVarP(&includeDFD, "include-dfd", "", false, "Include dfd in the report file")
	cmd.Flags().BoolVarP(&includeGlossary, "include-glossary", "", false, "Include glossary report")

	return CmdGetReports{cmd}
}
