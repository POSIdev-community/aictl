package update

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
)

type CmdUpdateScaFeeds struct {
	*cobra.Command
}

type UseCaseUpdateScaFeeds interface {
	Execute(ctx context.Context, path, version string) error
}

func NewUpdateScaFeedsCmd(uc UseCaseUpdateScaFeeds) CmdUpdateScaFeeds {
	var (
		path    string
		version string
	)

	cmd := &cobra.Command{
		Use:   "sca-feeds <path>",
		Short: "Upload SCA feeds package",
		Long:  `Upload a SCA feeds zip archive to the server (AIE ≥ 6.3). Requires --version.`,
		Example: `  aictl update sca-feeds ./AI.SCA.Feeds.47.zip --version 47
  aictl update sca-feeds ./feeds.zip --version 1.2.3`,
		Args: cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			path = strings.TrimSpace(args[0])
			if path == "" {
				return validation.NewError("empty sca feeds path")
			}

			version = strings.TrimSpace(version)
			if version == "" {
				return validation.NewError("empty version")
			}

			if !fshelper.PathExists(path) {
				return validation.NewError("path does not exist")
			}

			if !fshelper.IsFile(path) {
				return validation.NewError("path is not a file")
			}

			if !fshelper.IsArchive(path) {
				return validation.NewError("path is not a zip archive")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, path, version); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'update sca-feeds' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&version, "version", "", "Package version (required)")
	_ = cmd.MarkFlagRequired("version")

	return CmdUpdateScaFeeds{cmd}
}
