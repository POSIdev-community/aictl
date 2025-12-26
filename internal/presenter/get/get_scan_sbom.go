package get

import (
	"github.com/spf13/cobra"
)

type CmdGetScanSbom struct {
	*cobra.Command
}

func NewGetScanSbomCmd() CmdGetScanSbom {
	cmd := &cobra.Command{
		Use:   "sbom <scan-id>",
		Short: "Get scan sbom",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			panic("not implemented")
		},
	}

	return CmdGetScanSbom{cmd}
}
