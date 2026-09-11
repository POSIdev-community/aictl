package utils_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	domainlicense "github.com/POSIdev-community/aictl/internal/core/domain/license"
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	usecaseutils "github.com/POSIdev-community/aictl/internal/core/usecase/.utils"
)

type stubLicenseAI struct {
	lic      *domainlicense.License
	licErr   error
	settings settings.ScanSettings
	getErr   error
	setErr   error
	setCalls int
	lastSet  *settings.ScanSettings
}

func (s *stubLicenseAI) GetLicense(context.Context) (*domainlicense.License, error) {
	return s.lic, s.licErr
}

func (s *stubLicenseAI) GetProjectSettings(context.Context, uuid.UUID) (settings.ScanSettings, error) {
	return s.settings, s.getErr
}

func (s *stubLicenseAI) SetProjectSettings(_ context.Context, _ uuid.UUID, st *settings.ScanSettings) error {
	s.setCalls++
	s.lastSet = st

	return s.setErr
}

type stubLicenseCLI struct {
	msgs []string
}

func (s *stubLicenseCLI) ShowTextf(_ context.Context, format string, a ...any) {
	s.msgs = append(s.msgs, fmt.Sprintf(format, a...))
}

func TestApplyLicenseToProjectSettings_DisablesAndPersists(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	projectID := uuid.New()
	ai := &stubLicenseAI{
		lic: &domainlicense.License{
			Languages: []string{"Go"},
			LicensedModules: []domainlicense.LicensedModule{
				{Enabled: false, ScanModuleTypes: []string{domainlicense.ScanModuleSoftwareCompositionAnalysis}},
			},
		},
		settings: settings.ScanSettings{
			Languages: []string{"Go"},
			WhiteBoxSettings: settings.WhiteBoxSettings{
				StaticCodeAnalysisEnabled: true,
				SearchWithScaEnabled:      true,
			},
		},
	}
	cli := &stubLicenseCLI{}

	require.NoError(t, usecaseutils.ApplyLicenseToProjectSettings(ctx, ai, cli, projectID, usecaseutils.ApplyLicenseOptions{}))
	require.Equal(t, 1, ai.setCalls)
	require.False(t, ai.lastSet.WhiteBoxSettings.SearchWithScaEnabled)
	require.True(t, ai.lastSet.WhiteBoxSettings.StaticCodeAnalysisEnabled)
	require.Len(t, cli.msgs, 1)
	require.Contains(t, cli.msgs[0], domainlicense.ScanModuleSoftwareCompositionAnalysis)
}

func TestApplyLicenseToProjectSettings_UnlicensedLanguage(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	ai := &stubLicenseAI{
		lic: &domainlicense.License{Languages: []string{"Java"}},
		settings: settings.ScanSettings{
			Languages: []string{"Go"},
			WhiteBoxSettings: settings.WhiteBoxSettings{
				StaticCodeAnalysisEnabled: true,
			},
		},
	}
	cli := &stubLicenseCLI{}

	err := usecaseutils.ApplyLicenseToProjectSettings(ctx, ai, cli, uuid.New(), usecaseutils.ApplyLicenseOptions{})
	require.Error(t, err)
	var msg *validation.MessageError
	require.ErrorAs(t, err, &msg)
	require.Equal(t, 0, ai.setCalls)
	require.Empty(t, cli.msgs)
}

func TestApplyLicenseToProjectSettings_SbomSkipsLanguages(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	ai := &stubLicenseAI{
		lic: &domainlicense.License{Languages: []string{"Java"}},
		settings: settings.ScanSettings{
			Languages: []string{"Go"},
			WhiteBoxSettings: settings.WhiteBoxSettings{
				StaticCodeAnalysisEnabled: true,
			},
		},
	}
	cli := &stubLicenseCLI{}

	require.NoError(t, usecaseutils.ApplyLicenseToProjectSettings(ctx, ai, cli, uuid.New(), usecaseutils.ApplyLicenseOptions{
		SkipLanguageCheck: true,
	}))
	require.Equal(t, 0, ai.setCalls)
}

func TestApplyLicenseToProjectSettings_NoModulesRemain(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	ai := &stubLicenseAI{
		lic: &domainlicense.License{
			LicensedModules: []domainlicense.LicensedModule{
				{Enabled: false, ScanModuleTypes: []string{domainlicense.ScanModuleSoftwareCompositionAnalysis}},
			},
		},
		settings: settings.ScanSettings{
			WhiteBoxSettings: settings.WhiteBoxSettings{
				SearchWithScaEnabled: true,
			},
		},
	}
	cli := &stubLicenseCLI{}

	err := usecaseutils.ApplyLicenseToProjectSettings(ctx, ai, cli, uuid.New(), usecaseutils.ApplyLicenseOptions{SkipLanguageCheck: true})
	require.Error(t, err)
	require.Contains(t, err.Error(), "no licensed scan modules remain")
	require.Equal(t, 0, ai.setCalls)
	require.Empty(t, cli.msgs)
}

func TestApplyLicenseToProjectSettings_NoChanges(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	ai := &stubLicenseAI{
		lic: &domainlicense.License{Languages: []string{"Go"}}, // legacy modules
		settings: settings.ScanSettings{
			Languages: []string{"Go"},
			WhiteBoxSettings: settings.WhiteBoxSettings{
				SearchWithScaEnabled: true,
			},
		},
	}
	cli := &stubLicenseCLI{}

	require.NoError(t, usecaseutils.ApplyLicenseToProjectSettings(ctx, ai, cli, uuid.New(), usecaseutils.ApplyLicenseOptions{}))
	require.Equal(t, 0, ai.setCalls)
}
