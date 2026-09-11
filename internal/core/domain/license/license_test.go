package license_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/license"
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

func TestLicense_AllowedScanModules_LegacyFormat(t *testing.T) {
	t.Parallel()

	lic := &license.License{LicensedModules: nil}
	require.ElementsMatch(t, license.LicensedScanModuleTypes(), lic.AllowedScanModules())
}

func TestLicense_AllowedScanModules_NewFormat(t *testing.T) {
	t.Parallel()

	lic := &license.License{
		LicensedModules: []license.LicensedModule{
			{
				ID:      "Sca",
				Enabled: true,
				ScanModuleTypes: []string{
					license.ScanModuleSoftwareCompositionAnalysis,
					license.ScanModuleComponents,
				},
			},
			{
				ID:              "MaliciousCode",
				Enabled:         false,
				ScanModuleTypes: []string{license.ScanModuleMaliciousCodeDetection},
			},
		},
	}

	require.ElementsMatch(t, []string{
		license.ScanModuleSoftwareCompositionAnalysis,
		license.ScanModuleComponents,
	}, lic.AllowedScanModules())
}

func TestFilterLicensedScanModules(t *testing.T) {
	t.Parallel()

	lic := &license.License{
		LicensedModules: []license.LicensedModule{
			{
				Enabled: false,
				ScanModuleTypes: []string{
					license.ScanModuleSoftwareCompositionAnalysis,
					license.ScanModuleComponents,
				},
			},
			{
				Enabled:         true,
				ScanModuleTypes: []string{license.ScanModuleMaliciousCodeDetection},
			},
		},
	}

	modules := []string{
		"StaticCodeAnalysis",
		license.ScanModuleSoftwareCompositionAnalysis,
		license.ScanModuleComponents,
		license.ScanModuleMaliciousCodeDetection,
	}

	filtered, removed := license.FilterLicensedScanModules(lic, modules)
	require.ElementsMatch(t, []string{
		"StaticCodeAnalysis",
		license.ScanModuleMaliciousCodeDetection,
	}, filtered)
	require.ElementsMatch(t, []string{
		license.ScanModuleSoftwareCompositionAnalysis,
		license.ScanModuleComponents,
	}, removed)
}

func TestCheckLanguages(t *testing.T) {
	t.Parallel()

	lic := &license.License{Languages: []string{"Java", "Python"}}

	require.NoError(t, license.CheckLanguages(lic, nil))
	require.NoError(t, license.CheckLanguages(lic, []string{"Python"}))

	err := license.CheckLanguages(lic, []string{"Java", "Python", "Go", "Python"})
	require.Error(t, err)
	var msg *validation.MessageError
	require.ErrorAs(t, err, &msg)
	require.Contains(t, err.Error(), "Go")
	require.NotContains(t, err.Error(), "Python, Python")

	err = license.CheckLanguages(&license.License{}, []string{"Go"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "Go")
}

func TestPrepareSettings_DisablesModules(t *testing.T) {
	t.Parallel()

	lic := &license.License{
		Languages: []string{"Go"},
		LicensedModules: []license.LicensedModule{
			{Enabled: false, ScanModuleTypes: []string{license.ScanModuleSoftwareCompositionAnalysis}},
			{Enabled: true, ScanModuleTypes: []string{license.ScanModuleMaliciousCodeDetection}},
		},
	}

	in := settings.ScanSettings{
		Languages: []string{"Go"},
		WhiteBoxSettings: settings.WhiteBoxSettings{
			StaticCodeAnalysisEnabled:     true,
			SearchWithScaEnabled:          true,
			SearchForMaliciousCodeEnabled: true,
		},
	}

	out, removed, err := license.PrepareSettings(lic, in, false)
	require.NoError(t, err)
	require.Equal(t, []string{license.ScanModuleSoftwareCompositionAnalysis}, removed)
	require.False(t, out.WhiteBoxSettings.SearchWithScaEnabled)
	require.True(t, out.WhiteBoxSettings.StaticCodeAnalysisEnabled)
	require.True(t, out.WhiteBoxSettings.SearchForMaliciousCodeEnabled)
	require.True(t, in.WhiteBoxSettings.SearchWithScaEnabled, "input must not be mutated")
}

func TestPrepareSettings_NoModulesRemain(t *testing.T) {
	t.Parallel()

	lic := &license.License{
		LicensedModules: []license.LicensedModule{
			{Enabled: false, ScanModuleTypes: []string{
				license.ScanModuleSoftwareCompositionAnalysis,
				license.ScanModuleComponents,
				license.ScanModuleMaliciousCodeDetection,
			}},
		},
	}

	in := settings.ScanSettings{
		WhiteBoxSettings: settings.WhiteBoxSettings{
			SearchWithScaEnabled: true,
		},
	}

	_, removed, err := license.PrepareSettings(lic, in, true)
	require.Error(t, err)
	require.Equal(t, []string{license.ScanModuleSoftwareCompositionAnalysis}, removed)
	var msg *validation.MessageError
	require.ErrorAs(t, err, &msg)
	require.Equal(t, "license: no licensed scan modules remain after filtering", err.Error())
}

func TestPrepareSettings_SkipLanguageCheck(t *testing.T) {
	t.Parallel()

	lic := &license.License{Languages: []string{"Java"}}
	in := settings.ScanSettings{
		Languages: []string{"Go"},
		WhiteBoxSettings: settings.WhiteBoxSettings{
			StaticCodeAnalysisEnabled: true,
		},
	}

	_, _, err := license.PrepareSettings(lic, in, true)
	require.NoError(t, err)

	_, _, err = license.PrepareSettings(lic, in, false)
	require.Error(t, err)
}

func TestPrepareSettings_BlackBoxDoesNotCount(t *testing.T) {
	t.Parallel()

	lic := &license.License{
		LicensedModules: []license.LicensedModule{
			{Enabled: false, ScanModuleTypes: []string{license.ScanModuleSoftwareCompositionAnalysis}},
		},
	}

	in := settings.ScanSettings{
		BlackBoxEnabled: true,
		WhiteBoxSettings: settings.WhiteBoxSettings{
			SearchWithScaEnabled: true,
		},
	}

	_, _, err := license.PrepareSettings(lic, in, true)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no licensed scan modules remain")
}

func TestHasAnyWhiteBoxModule(t *testing.T) {
	t.Parallel()

	require.False(t, license.HasAnyWhiteBoxModule(settings.WhiteBoxSettings{}))
	require.True(t, license.HasAnyWhiteBoxModule(settings.WhiteBoxSettings{PatternMatchingEnabled: true}))
}
