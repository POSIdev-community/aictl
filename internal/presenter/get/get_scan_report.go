package get

import (
	"fmt"

	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/spf13/cobra"
)

type PersistentPreRunEGetScanReportCmd _utils.RunE

type CmdGetScanReport struct {
	*cobra.Command
}

var (
	destPath        string
	includeComments bool
	includeDFD      bool
	includeGlossary bool
	l10n            string
)

func NewPersistentPreRunEGetScanReportCmd(prev PersistentPreRunEGetScanCmd) PersistentPreRunEGetScanReportCmd {
	return _utils.ChainRunE(prev, func(cmd *cobra.Command, args []string) error {
		if l10n == "" || (l10n != "en" && l10n != "ru") {
			return fmt.Errorf("the localization language '%s' is unknown, but 'en' or 'ru' may be used", l10n)
		}

		return nil
	})
}

func NewGetScanReportCmd(
	persistentPreRunE PersistentPreRunEGetScanReportCmd,
	cmdGetScanReportGitlab CmdGetScanReportGitlab,
	cmdGetScanReportPlain CmdGetScanReportPlain,
	cmdGetScanReportSarif CmdGetScanReportSarif) CmdGetScanReport {

	cmd := &cobra.Command{
		Short: "Get scan report",
		Use:   "report <scan-id>",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if destPath == "" {
				return fmt.Errorf("must specify -o")
			}

			return nil
		},
		PersistentPreRunE: persistentPreRunE,
		RunE: func(cmd *cobra.Command, args []string) error {
			panic("not implemented")
		},
	}

	cmd.AddCommand(cmdGetScanReportGitlab.Command)
	cmd.AddCommand(cmdGetScanReportPlain.Command)
	cmd.AddCommand(cmdGetScanReportSarif.Command)

	cmd.PersistentFlags().StringVarP(&destPath, "output", "o", "", "Destination path for the report file")
	cmd.PersistentFlags().BoolVar(&includeComments, "include-comments", false, "Include comments in the report file")
	cmd.PersistentFlags().BoolVar(&includeDFD, "include-dfd", false, "Include dfd in the report file")
	cmd.PersistentFlags().BoolVar(&includeGlossary, "include-glossary", false, "Include glossary report")
	cmd.PersistentFlags().StringVar(&l10n, "localization", "en", "Localization language: 'en', 'ru'")

	return CmdGetScanReport{cmd}
}
