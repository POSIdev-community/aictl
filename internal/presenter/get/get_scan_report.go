package get

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
)

var (
	destPath string
)

func NewGetScanReportCmd(cfg *config.Config, depsContainer *application.DependenciesContainer) *cobra.Command {
	cmd := &cobra.Command{
		Short: "Get scan report",
		Use:   "report",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if destPath == "" {
				return fmt.Errorf("must specify -o")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			panic("not implemented")
		},
	}

	cmd.AddCommand(NewGetScanReportPlainCmd(cfg, depsContainer))
	cmd.AddCommand(NewGetScanReportSarifCmd(cfg, depsContainer))

	cmd.PersistentFlags().StringVarP(&destPath, "output", "o", "", "Destination path for the report file")

	return cmd
}
