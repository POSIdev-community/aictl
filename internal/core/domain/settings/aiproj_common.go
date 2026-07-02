package settings

import (
	"fmt"

	domainaiproj "github.com/POSIdev-community/aictl/internal/core/domain/aiproj"
)

func (s *ScanSettings) UpdateFromParsed(parsed *domainaiproj.Result) error {
	if parsed == nil {
		return fmt.Errorf("parsed aiproj is nil")
	}

	switch parsed.Version {
	case "1.9":
		if parsed.V19 == nil {
			return fmt.Errorf("aiproj 1.9 payload is nil")
		}

		return s.UpdateFromV19(parsed.V19)
	case "1.10":
		if parsed.V110 == nil {
			return fmt.Errorf("aiproj 1.10 payload is nil")
		}

		return s.UpdateFromV110(parsed.V110)
	case "1.11":
		if parsed.V111 == nil {
			return fmt.Errorf("aiproj 1.11 payload is nil")
		}

		return s.UpdateFromV111(parsed.V111)
	default:
		return fmt.Errorf("unsupported aiproj version: %s", parsed.Version)
	}
}

func applyScanModules(s *ScanSettings, modules []string) error {
	if len(modules) == 0 {
		return nil
	}

	s.WhiteBoxSettings.StaticCodeAnalysisEnabled = false
	s.WhiteBoxSettings.PatternMatchingEnabled = false
	s.WhiteBoxSettings.SearchForVulnerableComponentsEnabled = false
	s.WhiteBoxSettings.SearchWithScaEnabled = false
	s.WhiteBoxSettings.SearchForConfigurationFlawsEnabled = false
	s.WhiteBoxSettings.SecretDetectionEnabled = false
	s.WhiteBoxSettings.SearchForMaliciousCodeEnabled = false

	for _, module := range modules {
		switch module {
		case string(domainaiproj.ScanModuleStaticCodeAnalysis):
			s.WhiteBoxSettings.StaticCodeAnalysisEnabled = true
		case string(domainaiproj.ScanModulePatternMatching):
			s.WhiteBoxSettings.PatternMatchingEnabled = true
		case string(domainaiproj.ScanModuleComponents):
			s.WhiteBoxSettings.SearchForVulnerableComponentsEnabled = true
		case string(domainaiproj.ScanModuleSoftwareCompositionAnalysis):
			s.WhiteBoxSettings.SearchWithScaEnabled = true
		case string(domainaiproj.ScanModuleConfiguration):
			s.WhiteBoxSettings.SearchForConfigurationFlawsEnabled = true
		case string(domainaiproj.ScanModuleSecretDetection):
			s.WhiteBoxSettings.SecretDetectionEnabled = true
		case string(domainaiproj.ScanModuleMaliciousCodeDetection):
			s.WhiteBoxSettings.SearchForMaliciousCodeEnabled = true
		case string(domainaiproj.ScanModuleBlackBox):
			s.BlackBoxEnabled = true
		default:
			return fmt.Errorf("unsupported module: %s", module)
		}
	}

	return nil
}

func applyLanguages(s *ScanSettings, languages []string) {
	if len(languages) == 0 {
		return
	}

	s.Languages = make([]string, len(languages))
	for i, lang := range languages {
		switch lang {
		case string(domainaiproj.ProgrammingLanguageCSharpWindowsLinux):
			s.Languages[i] = "CSharp"
		case string(domainaiproj.ProgrammingLanguageCSharpWindows):
			s.Languages[i] = "CSharpWinOnly"
		default:
			s.Languages[i] = lang
		}
	}
}
