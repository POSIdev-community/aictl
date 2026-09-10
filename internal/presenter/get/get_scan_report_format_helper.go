package get

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
)

type defaultReportFormatSpec struct {
	use         string
	short       string
	long        string
	examplePath string
	reportType  report.ReportType
}

var defaultReportFormatSpecs = []defaultReportFormatSpec{
	{use: "autocheck", short: "Get scan report in Autocheck format", long: "Download the scan report in Autocheck format for the given scan id. Project id comes from context or parent -p. Output path via -o; use -f to overwrite.", examplePath: "./out.autocheck", reportType: report.AutoCheck},
	{use: "gitlab", short: "Get scan report in GitLab format", long: "Download the scan report in GitLab format for the given scan id. Project id comes from context or parent -p. Output path via -o; use -f to overwrite.", examplePath: "./out.gitlab", reportType: report.Gitlab},
	{use: "json", short: "Get scan report in JSON format", long: "Download the scan report in JSON format for the given scan id. Project id comes from context or parent -p. Output path via -o; use -f to overwrite.", examplePath: "./out.json", reportType: report.Json},
	{use: "json-v2", short: "Get scan report in JSON v2 format", long: "Download the scan report in JSON v2 format for the given scan id. Project id comes from context or parent -p. Output path via -o; use -f to overwrite.", examplePath: "./out.json-v2", reportType: report.JsonV2},
	{use: "markdown", short: "Get scan report in Markdown format", long: "Download the scan report in Markdown format for the given scan id. Project id comes from context or parent -p. Output path via -o; use -f to overwrite.", examplePath: "./out.markdown", reportType: report.Markdown},
	{use: "nist", short: "Get scan report in NIST format", long: "Download the scan report in NIST format for the given scan id. Project id comes from context or parent -p. Output path via -o; use -f to overwrite.", examplePath: "./out.nist", reportType: report.Nist},
	{use: "oud4", short: "Get scan report in OUD4 format", long: "Download the scan report in OUD4 format for the given scan id. Project id comes from context or parent -p. Output path via -o; use -f to overwrite.", examplePath: "./out.oud4", reportType: report.Oud4},
	{use: "owasp", short: "Get scan report in OWASP format", long: "Download the scan report in OWASP format for the given scan id. Project id comes from context or parent -p. Output path via -o; use -f to overwrite.", examplePath: "./out.xml", reportType: report.Owasp},
	{use: "owaspm", short: "Get scan report in OWASP Mobile format", long: "Download the scan report in OWASP Mobile format for the given scan id. Project id comes from context or parent -p. Output path via -o; use -f to overwrite.", examplePath: "./out.owaspm", reportType: report.Owaspm},
	{use: "pcidss", short: "Get scan report in PCI DSS format", long: "Download the scan report in PCI DSS format for the given scan id. Project id comes from context or parent -p. Output path via -o; use -f to overwrite.", examplePath: "./out.pcidss", reportType: report.Pcidss},
	{use: "plain", short: "Get scan report in plain HTML format", long: "Download the scan report in plain HTML format for the given scan id. Project id comes from context or parent -p. Output path via -o; use -f to overwrite.", examplePath: "./out.html", reportType: report.PlainReport},
	{use: "sans", short: "Get scan report in SANS format", long: "Download the scan report in SANS format for the given scan id. Project id comes from context or parent -p. Output path via -o; use -f to overwrite.", examplePath: "./out.sans", reportType: report.Sans},
	{use: "sarif", short: "Get scan report in SARIF format", long: "Download the scan report in SARIF format for the given scan id. Project id comes from context or parent -p. Output path via -o; use -f to overwrite.", examplePath: "./out.sarif", reportType: report.Sarif},
	{use: "xml", short: "Get scan report in XML format", long: "Download the scan report in XML format for the given scan id. Project id comes from context or parent -p. Output path via -o; use -f to overwrite.", examplePath: "./out.xml", reportType: report.Xml},
}

type useCaseDefaultReportWithFilters interface {
	Execute(ctx context.Context, scanId uuid.UUID, reportType report.ReportType, fullDestPath string, includeComments, includeDFD, includeGlossary bool, l10n string, filters report.Filters) error
}

func addDefaultReportFormatSubcommands(parent *cobra.Command, uc useCaseDefaultReportWithFilters, commandPathPrefix string) {
	for _, spec := range defaultReportFormatSpecs {
		spec := spec
		parent.AddCommand(&cobra.Command{
			Use:   spec.use + " <scan-id>",
			Short: spec.short,
			Long:  spec.long,
			Example: fmt.Sprintf(
				"  aictl %s %s <scan-id> -o %s --level-high\n  aictl %s %s <scan-id> -o %s -f --level-high",
				commandPathPrefix, spec.use, spec.examplePath, commandPathPrefix, spec.use, spec.examplePath,
			),
			Args: cobra.MaximumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx := cmd.Context()
				if err := uc.Execute(ctx, scanId, spec.reportType, outPath, includeComments, includeDFD, includeGlossary, l10n, reportFilters); err != nil {
					cmd.SilenceUsage = true
					return fmt.Errorf("'%s %s' usecase call: %w", commandPathPrefix, spec.use, err)
				}
				return nil
			},
		})
	}
}
