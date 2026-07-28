package get

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/regexfilter"
)

type CmdGetReportTemplates struct {
	*cobra.Command
}

type UseCaseGetReportTemplates interface {
	Execute(ctx context.Context, filter regexfilter.RegexFilter, quite bool, localization string) error
}

func NewGetReportTemplatesCmd(uc UseCaseGetReportTemplates) CmdGetReportTemplates {
	var (
		quite        bool
		localization string
		regexFilter  regexfilter.RegexFilter
	)

	cmd := &cobra.Command{
		Use:   "report-templates [<regex>]",
		Short: "Get report templates",
		Long:  `List available report templates matching an optional regex filter.`,
		Example: `  aictl get report-templates
  aictl get report-templates owasp
  aictl get report-templates -q --localization ru`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if localization == "" || (localization != "en" && localization != "ru") {
				return fmt.Errorf("the localization language '%s' is unknown, but 'en' or 'ru' may be used", localization)
			}

			filter := strings.Join(args, " ")
			var err error
			regexFilter, err = regexfilter.NewRegexFilter(filter)
			if err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("new filter: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, regexFilter, quite, localization); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get report-templates' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&quite, "quite", "q", false, "Print only ids")
	cmd.Flags().StringVar(&localization, "localization", "en", "Report localization language: 'en' or 'ru'")

	return CmdGetReportTemplates{cmd}
}
