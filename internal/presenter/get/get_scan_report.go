package get

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type PersistentPreRunEGetScanReportCmd _utils.RunE

type CmdGetScanReport struct {
	*cobra.Command
}

type UseCaseGetScanReport interface {
	Execute(ctx context.Context, scanId uuid.UUID, customReportName string, outPath string, includeComments, includeDFD, includeGlossary bool, l10n string) error
}

var (
	outPath             string
	forceRewriteOutPath bool

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

		if outPath != "" {
			if fshelper.PathExists(outPath) && !forceRewriteOutPath {
				return errs.NewValidationError("'output' path exists")
			}
		}

		return nil
	})
}

func NewGetScanReportCmd(
	uc UseCaseGetScanReport,
	persistentPreRunE PersistentPreRunEGetScanReportCmd,
	cmdGetScanReportAutocheck CmdGetScanReportAutocheck,
	cmdGetScanReportGitlab CmdGetScanReportGitlab,
	cmdGetScanReportJson CmdGetScanReportJson,
	cmdGetScanReportMarkdown CmdGetScanReportMarkdown,
	cmdGetScanReportNist CmdGetScanReportNist,
	cmdGetScanReportOud4 CmdGetScanReportOud4,
	cmdGetScanReportOwasp CmdGetScanReportOwasp,
	cmdGetScanReportOwaspm CmdGetScanReportOwaspm,
	cmdGetScanReportPcidss CmdGetScanReportPcidss,
	cmdGetScanReportPlain CmdGetScanReportPlain,
	cmdGetScanReportSans CmdGetScanReportSans,
	cmdGetScanReportSarif CmdGetScanReportSarif,
	cmdGetScanReportXml CmdGetScanReportXml) CmdGetScanReport {

	cmd := &cobra.Command{
		Use:               "report <report-name> <scan-id>",
		Short:             "Get scan report",
		Args:              cobra.ExactArgs(2),
		PersistentPreRunE: persistentPreRunE,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, customReportName, outPath, includeComments, includeDFD, includeGlossary, l10n); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan report' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.AddCommand(cmdGetScanReportAutocheck.Command)
	cmd.AddCommand(cmdGetScanReportGitlab.Command)
	cmd.AddCommand(cmdGetScanReportJson.Command)
	cmd.AddCommand(cmdGetScanReportMarkdown.Command)
	cmd.AddCommand(cmdGetScanReportNist.Command)
	cmd.AddCommand(cmdGetScanReportOud4.Command)
	cmd.AddCommand(cmdGetScanReportOwasp.Command)
	cmd.AddCommand(cmdGetScanReportOwaspm.Command)
	cmd.AddCommand(cmdGetScanReportPcidss.Command)
	cmd.AddCommand(cmdGetScanReportPlain.Command)
	cmd.AddCommand(cmdGetScanReportSans.Command)
	cmd.AddCommand(cmdGetScanReportSarif.Command)
	cmd.AddCommand(cmdGetScanReportXml.Command)

	cmd.PersistentFlags().StringVarP(&outPath, "output", "o", "", "Destination path for the report file")
	cmd.PersistentFlags().BoolVarP(&forceRewriteOutPath, "force", "f", false, "Force rewrite output file")

	cmd.PersistentFlags().BoolVar(&includeComments, "include-comments", false, "Include comments in the report file")
	cmd.PersistentFlags().BoolVar(&includeDFD, "include-dfd", false, "Include dfd in the report file")
	cmd.PersistentFlags().BoolVar(&includeGlossary, "include-glossary", false, "Include glossary report")
	cmd.PersistentFlags().StringVar(&l10n, "localization", "en", "Localization language: 'en', 'ru'")

	return CmdGetScanReport{cmd}
}
