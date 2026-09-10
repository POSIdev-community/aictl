package get

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

var (
	filterTypes       []string
	filterLanguages   []string
	filterScanModules []string

	reportFilters report.Filters
)

func resetReportFilterFlags() {
	filterTypes = nil
	filterLanguages = nil
	filterScanModules = nil
	reportFilters = report.EmptyFilters()
}

func addReportFilterFlags(cmd *cobra.Command) {
	fs := cmd.PersistentFlags()

	fs.Bool("level-high", false, "Include high severity issues")
	fs.Bool("level-medium", false, "Include medium severity issues")
	fs.Bool("level-low", false, "Include low severity issues")
	fs.Bool("level-potential", false, "Include potential severity issues")

	fs.Bool("status-undefined", false, "Include issues with undefined confirmation status")
	fs.Bool("status-confirmed", false, "Include confirmed issues")
	fs.Bool("status-confirmed-auto", false, "Include auto-confirmed issues")
	fs.Bool("status-rejected", false, "Include rejected issues")

	fs.Bool("mode-entry-point", false, "Include entry-point scan mode issues")
	fs.Bool("mode-public-methods", false, "Include public-methods scan mode issues")
	fs.Bool("mode-root-function", false, "Include root-function scan mode issues")
	fs.Bool("mode-others", false, "Include other scan mode issues")

	fs.Bool("found-this-scan", false, "Include issues found in this scan")
	fs.Bool("found-prev-scan", false, "Include issues found in a previous scan")

	fs.Bool("conditional", false, "Include conditional issues")
	fs.Bool("non-conditional", false, "Include non-conditional issues")
	fs.Bool("suppressed", false, "Include suppressed issues")
	fs.Bool("non-suppressed", false, "Include non-suppressed issues")

	fs.Bool("suspected", false, "Include suspected issues")
	fs.Bool("second-level", false, "Include second-level issues")
	fs.Bool("only-favorite", false, "Include only favorite issues")

	fs.StringArrayVar(&filterTypes, "type", nil, "Vulnerability type filter (repeatable)")
	fs.StringArrayVar(&filterScanModules, "scan-module", nil, fmt.Sprintf("Scan module filter (repeatable): %s", strings.Join(report.AllowedScanModules, ", ")))
	fs.StringArrayVar(&filterLanguages, "language", nil, fmt.Sprintf("Language filter (repeatable): %s", strings.Join(report.AllowedLanguages, ", ")))
}

func buildReportFiltersFromFlags(cmd *cobra.Command) (report.Filters, error) {
	f := report.Filters{
		Apply:       true,
		Types:       append([]string(nil), filterTypes...),
		Languages:   append([]string(nil), filterLanguages...),
		ScanModules: append([]string(nil), filterScanModules...),
	}

	setTrueIfChanged := func(name string, dest **bool) {
		flag := cmd.Flag(name)
		if flag != nil && flag.Changed {
			t := true
			*dest = &t
		}
	}

	setTrueIfChanged("level-high", &f.LevelHigh)
	setTrueIfChanged("level-medium", &f.LevelMedium)
	setTrueIfChanged("level-low", &f.LevelLow)
	setTrueIfChanged("level-potential", &f.LevelPotential)

	setTrueIfChanged("status-undefined", &f.StatusUndefined)
	setTrueIfChanged("status-confirmed", &f.StatusConfirmed)
	setTrueIfChanged("status-confirmed-auto", &f.StatusConfirmedAuto)
	setTrueIfChanged("status-rejected", &f.StatusRejected)

	setTrueIfChanged("mode-entry-point", &f.ModeEntryPoint)
	setTrueIfChanged("mode-public-methods", &f.ModePublicMethods)
	setTrueIfChanged("mode-root-function", &f.ModeRootFunction)
	setTrueIfChanged("mode-others", &f.ModeOthers)

	setTrueIfChanged("found-this-scan", &f.FoundThisScan)
	setTrueIfChanged("found-prev-scan", &f.FoundPrevScan)

	setTrueIfChanged("conditional", &f.Conditional)
	setTrueIfChanged("non-conditional", &f.NonConditional)
	setTrueIfChanged("suppressed", &f.Suppressed)
	setTrueIfChanged("non-suppressed", &f.NonSuppressed)

	setTrueIfChanged("suspected", &f.Suspected)
	setTrueIfChanged("second-level", &f.SecondLevel)
	setTrueIfChanged("only-favorite", &f.OnlyFavorite)

	f.Normalize()

	if !f.HasAny() {
		return report.Filters{}, validation.NewError("at least one filter flag is required")
	}

	if err := f.ValidateArrays(); err != nil {
		return report.Filters{}, err
	}

	// Ensure empty slices (not nil) for arrays when useFilters=true.
	if f.Types == nil {
		f.Types = []string{}
	}
	if f.Languages == nil {
		f.Languages = []string{}
	}
	if f.ScanModules == nil {
		f.ScanModules = []string{}
	}

	return f, nil
}
