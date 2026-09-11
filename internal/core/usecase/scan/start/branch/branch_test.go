package branch_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	domainlicense "github.com/POSIdev-community/aictl/internal/core/domain/license"
	"github.com/POSIdev-community/aictl/internal/core/domain/scantype"
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/internal/core/usecase/scan/start/branch"
)

type mockAI struct {
	mock.Mock
}

func (m *mockAI) InitializeWithRetry(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}

func (m *mockAI) GetLicense(ctx context.Context) (*domainlicense.License, error) {
	args := m.Called(ctx)
	lic, _ := args.Get(0).(*domainlicense.License)

	return lic, args.Error(1)
}

func (m *mockAI) GetProjectSettings(ctx context.Context, projectId uuid.UUID) (settings.ScanSettings, error) {
	args := m.Called(ctx, projectId)
	st, _ := args.Get(0).(settings.ScanSettings)

	return st, args.Error(1)
}

func (m *mockAI) SetProjectSettings(ctx context.Context, projectId uuid.UUID, st *settings.ScanSettings) error {
	args := m.Called(ctx, projectId, st)

	return args.Error(0)
}

func (m *mockAI) StartScanBranch(ctx context.Context, branchId uuid.UUID, scanLabel string, scanType scantype.Type) (uuid.UUID, error) {
	args := m.Called(ctx, branchId, scanLabel, scanType)
	id, _ := args.Get(0).(uuid.UUID)

	return id, args.Error(1)
}

type mockCLI struct {
	mock.Mock
}

func (m *mockCLI) ShowTextf(ctx context.Context, format string, a ...any) {
	m.Called(ctx, format, a)
}

func (m *mockCLI) ReturnText(ctx context.Context, text string) {
	m.Called(ctx, text)
}

func TestUseCase_Execute_DisablesModuleThenStarts(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	projectID := uuid.New()
	branchID := uuid.New()
	scanID := uuid.New()
	cfg := config.NewConfig(config.Uri{}, "", true, projectID, branchID)

	ai := &mockAI{}
	cli := &mockCLI{}
	ai.Test(t)
	cli.Test(t)
	t.Cleanup(func() {
		ai.AssertExpectations(t)
		cli.AssertExpectations(t)
	})

	ai.On("InitializeWithRetry", ctx).Return(nil).Once()
	ai.On("GetLicense", ctx).Return(&domainlicense.License{
		Languages: []string{"Go"},
		LicensedModules: []domainlicense.LicensedModule{
			{Enabled: false, ScanModuleTypes: []string{domainlicense.ScanModuleSoftwareCompositionAnalysis}},
		},
	}, nil).Once()
	ai.On("GetProjectSettings", ctx, projectID).Return(settings.ScanSettings{
		Languages: []string{"Go"},
		WhiteBoxSettings: settings.WhiteBoxSettings{
			StaticCodeAnalysisEnabled: true,
			SearchWithScaEnabled:      true,
		},
	}, nil).Once()
	ai.On("SetProjectSettings", ctx, projectID, mock.MatchedBy(func(st *settings.ScanSettings) bool {
		return st != nil && !st.WhiteBoxSettings.SearchWithScaEnabled && st.WhiteBoxSettings.StaticCodeAnalysisEnabled
	})).Return(nil).Once()
	cli.On("ShowTextf", ctx, "license: disabled unlicensed scan module '%s'", mock.Anything).Return().Once()
	cli.On("ShowTextf", ctx, "starting scan, project-id '%v', branch-id '%v'", mock.Anything).Return().Once()
	ai.On("StartScanBranch", ctx, branchID, "label", scantype.Incremental).Return(scanID, nil).Once()
	cli.On("ShowTextf", ctx, "scan started, scan-id '%v'", mock.Anything).Return().Once()
	cli.On("ReturnText", ctx, scanID.String()).Return().Once()

	uc, err := branch.NewUseCase(ai, cli, cfg)
	require.NoError(t, err)
	require.NoError(t, uc.Execute(ctx, "label", scantype.Incremental))
}

func TestUseCase_Execute_UnlicensedLanguageDoesNotStart(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	projectID := uuid.New()
	branchID := uuid.New()
	cfg := config.NewConfig(config.Uri{}, "", true, projectID, branchID)

	ai := &mockAI{}
	cli := &mockCLI{}
	ai.Test(t)
	cli.Test(t)
	t.Cleanup(func() {
		ai.AssertExpectations(t)
		cli.AssertExpectations(t)
	})

	ai.On("InitializeWithRetry", ctx).Return(nil).Once()
	ai.On("GetLicense", ctx).Return(&domainlicense.License{Languages: []string{"Java"}}, nil).Once()
	ai.On("GetProjectSettings", ctx, projectID).Return(settings.ScanSettings{
		Languages: []string{"Go"},
		WhiteBoxSettings: settings.WhiteBoxSettings{
			StaticCodeAnalysisEnabled: true,
		},
	}, nil).Once()

	uc, err := branch.NewUseCase(ai, cli, cfg)
	require.NoError(t, err)

	err = uc.Execute(ctx, "label", scantype.Incremental)
	require.Error(t, err)
	var msg *validation.MessageError
	require.ErrorAs(t, err, &msg)
}
