package license

import (
	"slices"
	"strings"

	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

// Scan module types gated by enterprise license (same set as infr-agent).
const (
	ScanModuleSoftwareCompositionAnalysis = "SoftwareCompositionAnalysis"
	ScanModuleComponents                  = "Components"
	ScanModuleMaliciousCodeDetection      = "MaliciousCodeDetection"
)

const (
	errTextLanguagesPrefix = "the following programming languages are not included in the license: "
	errNoModulesRemain     = "license: no licensed scan modules remain after filtering"
)

// LicensedModule describes one licensed scan capability and the ScanModule values it covers.
type LicensedModule struct {
	ID              string
	Enabled         bool
	ScanModuleTypes []string
}

// License is the CLI view of GET /api/license (without IsValid — that is checked at Initialize).
type License struct {
	Languages       []string
	LicensedModules []LicensedModule // nil means legacy format (SCA and MOLOT enabled).
}

// LicensedScanModuleTypes returns scan modules subject to license filtering.
func LicensedScanModuleTypes() []string {
	return []string{
		ScanModuleSoftwareCompositionAnalysis,
		ScanModuleComponents,
		ScanModuleMaliciousCodeDetection,
	}
}

// AllowedScanModules returns scan modules permitted by the current license.
func (l *License) AllowedScanModules() []string {
	if l == nil || l.LicensedModules == nil {
		return LicensedScanModuleTypes()
	}

	allowed := make([]string, 0, len(LicensedScanModuleTypes()))
	for _, module := range l.LicensedModules {
		if !module.Enabled {
			continue
		}
		allowed = append(allowed, module.ScanModuleTypes...)
	}

	return allowed
}

// FilterLicensedScanModules removes license-gated modules that are not allowed.
func FilterLicensedScanModules(lic *License, modules []string) (filtered, removed []string) {
	if len(modules) == 0 {
		return nil, nil
	}

	allowed := lic.AllowedScanModules()
	licensable := LicensedScanModuleTypes()

	filtered = make([]string, 0, len(modules))
	for _, module := range modules {
		if !slices.Contains(licensable, module) {
			filtered = append(filtered, module)
			continue
		}
		if slices.Contains(allowed, module) {
			filtered = append(filtered, module)
			continue
		}
		removed = append(removed, module)
	}

	return filtered, removed
}

// CheckLanguages ensures settings languages are covered by the license.
// Empty settings languages pass. Empty license allowlist rejects any non-empty settings list.
func CheckLanguages(lic *License, settingsLanguages []string) error {
	if len(settingsLanguages) == 0 {
		return nil
	}

	var allowedList []string
	if lic != nil {
		allowedList = lic.Languages
	}

	allowed := make(map[string]struct{}, len(allowedList))
	for _, lg := range allowedList {
		allowed[lg] = struct{}{}
	}

	seen := make(map[string]struct{})
	invalid := make([]string, 0)

	for _, lang := range settingsLanguages {
		if _, ok := allowed[lang]; ok {
			continue
		}
		if _, dup := seen[lang]; dup {
			continue
		}
		seen[lang] = struct{}{}
		invalid = append(invalid, lang)
	}

	if len(invalid) == 0 {
		return nil
	}

	return validation.NewMessageError(errTextLanguagesPrefix + strings.Join(invalid, ", "))
}

// EnabledLicensableModules returns license-gated modules currently enabled in white-box settings.
func EnabledLicensableModules(wb settings.WhiteBoxSettings) []string {
	out := make([]string, 0, 3)
	if wb.SearchWithScaEnabled {
		out = append(out, ScanModuleSoftwareCompositionAnalysis)
	}
	if wb.SearchForVulnerableComponentsEnabled {
		out = append(out, ScanModuleComponents)
	}
	if wb.SearchForMaliciousCodeEnabled {
		out = append(out, ScanModuleMaliciousCodeDetection)
	}

	return out
}

// DisableModules sets the corresponding WhiteBoxSettings flags to false.
func DisableModules(wb *settings.WhiteBoxSettings, modules []string) {
	if wb == nil {
		return
	}
	for _, m := range modules {
		switch m {
		case ScanModuleSoftwareCompositionAnalysis:
			wb.SearchWithScaEnabled = false
		case ScanModuleComponents:
			wb.SearchForVulnerableComponentsEnabled = false
		case ScanModuleMaliciousCodeDetection:
			wb.SearchForMaliciousCodeEnabled = false
		}
	}
}

// HasAnyWhiteBoxModule reports whether any white-box analysis flag is enabled (BlackBox ignored).
func HasAnyWhiteBoxModule(wb settings.WhiteBoxSettings) bool {
	return wb.StaticCodeAnalysisEnabled ||
		wb.PatternMatchingEnabled ||
		wb.SearchForConfigurationFlawsEnabled ||
		wb.SecretDetectionEnabled ||
		wb.SearchWithScaEnabled ||
		wb.SearchForVulnerableComponentsEnabled ||
		wb.SearchForMaliciousCodeEnabled
}

// PrepareSettings checks languages (unless skipLanguageCheck) and plans module disables.
// On "no modules remain" it returns ErrNoModulesRemain and does not require persisting settings.
// When removed is empty and err is nil, settings are unchanged.
func PrepareSettings(lic *License, scanSettings settings.ScanSettings, skipLanguageCheck bool) (updated settings.ScanSettings, removed []string, err error) {
	if !skipLanguageCheck {
		if err := CheckLanguages(lic, scanSettings.Languages); err != nil {
			return scanSettings, nil, err
		}
	}

	enabled := EnabledLicensableModules(scanSettings.WhiteBoxSettings)
	_, removed = FilterLicensedScanModules(lic, enabled)

	updated = scanSettings
	DisableModules(&updated.WhiteBoxSettings, removed)

	if !HasAnyWhiteBoxModule(updated.WhiteBoxSettings) {
		return scanSettings, removed, validation.NewMessageError(errNoModulesRemain)
	}

	return updated, removed, nil
}
