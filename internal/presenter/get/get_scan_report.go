package get

import (
	"fmt"

	"github.com/spf13/cobra"
)

type CmdGetScanReport struct {
	*cobra.Command
}

var (
	destPath        string
	includeComments bool
	includeDFD      bool
	includeGlossary bool
)

func NewGetScanReportCmd(
	cmdGetScanReportGitlab CmdGetScanReportGitlab,
	cmdGetScanReportPlain CmdGetScanReportPlain,
	cmdGetScanReportSarif CmdGetScanReportSarif) CmdGetScanReport {

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

	cmd.AddCommand(cmdGetScanReportGitlab.Command)
	cmd.AddCommand(cmdGetScanReportPlain.Command)
	cmd.AddCommand(cmdGetScanReportSarif.Command)

	cmd.PersistentFlags().StringVarP(&destPath, "output", "o", "", "Destination path for the report file")
	cmd.PersistentFlags().BoolVarP(&includeComments, "include-comments", "", false, "Include comments in the report file")
	cmd.PersistentFlags().BoolVarP(&includeDFD, "include-dfd", "", false, "Include dfd in the report file")
	cmd.PersistentFlags().BoolVarP(&includeGlossary, "include-glossary", "", false, "Include glossary report")

	return CmdGetScanReport{cmd}
}
