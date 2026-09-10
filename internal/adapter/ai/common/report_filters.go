package common

import (
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/internal/core/domain/version"
)

var (
	minVersion60, _ = version.NewVersion("6.0.0")
	minVersion61, _ = version.NewVersion("6.1.0")
)

// ValidateReportFiltersForVersion rejects filter values that require a newer AI server.
// Error messages name the minimum version (e.g. ">= 6.0.0").
func ValidateReportFiltersForVersion(f report.Filters, serverVer version.Version) error {
	if !f.Apply {
		return nil
	}

	for _, m := range f.ScanModules {
		switch m {
		case "SecretDetection", "MaliciousCodeDetection":
			if serverVer.Less(minVersion60) {
				return validation.NewError(fmt.Sprintf(
					"scan module %q requires AI server version >= 6.0.0 (current %s)",
					m, serverVer.String()))
			}
		}
	}

	for _, lang := range f.Languages {
		switch lang {
		case "OneC":
			if serverVer.Less(minVersion60) {
				return validation.NewError(fmt.Sprintf(
					"language %q requires AI server version >= 6.0.0 (current %s)",
					lang, serverVer.String()))
			}
		case "Dart":
			if serverVer.Less(minVersion61) {
				return validation.NewError(fmt.Sprintf(
					"language %q requires AI server version >= 6.1.0 (current %s)",
					lang, serverVer.String()))
			}
		}
	}

	return nil
}
