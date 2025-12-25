package get

import (
	"github.com/spf13/cobra"
)

type CmdGetScanLogs struct {
	*cobra.Command
}

func NewGetScanLogsCmd() CmdGetScanLogs {
	cmd := &cobra.Command{
		Use:   "logs <scan-id>",
		Short: "Get scan logs",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			panic("not implemented")
		},
	}

	return CmdGetScanLogs{cmd}
}
