package settings_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	domainproject "github.com/POSIdev-community/aictl/internal/core/domain/project"
	domainsettings "github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/internal/core/domain/version"
	"github.com/POSIdev-community/aictl/internal/core/usecase/update/project/settings"
)

type fakeAI struct {
	version         version.Version
	current         domainsettings.ScanSettings
	setCalled       bool
	lastSetSettings *domainsettings.ScanSettings
}

func (f *fakeAI) InitializeWithRetry(context.Context) error { return nil }

func (f *fakeAI) GetProject(_ context.Context, projectId uuid.UUID) (*domainproject.Project, error) {
	p := domainproject.NewProject(projectId, "test")
	return &p, nil
}

func (f *fakeAI) GetVersion(context.Context) (version.Version, error) { return f.version, nil }

func (f *fakeAI) GetProjectSettings(context.Context, uuid.UUID) (domainsettings.ScanSettings, error) {
	return f.current, nil
}

func (f *fakeAI) SetProjectSettings(_ context.Context, _ uuid.UUID, s *domainsettings.ScanSettings) error {
	f.setCalled = true
	f.lastSetSettings = s

	return nil
}

type fakeCLI struct{}

func TestUseCase_Execute_Rejects54(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")

	serverVersion, err := version.NewVersion("5.4.0")
	require.NoError(t, err)

	cfg := newTestConfig(projectID)
	uc, err := settings.NewUseCase(&fakeAI{version: serverVersion}, fakeCLI{}, cfg)
	require.NoError(t, err)

	priority := domainsettings.PriorityHigh
	patch := domainsettings.ProjectSettingsPatch{Priority: &priority}

	err = uc.Execute(ctx, patch)
	require.Error(t, err)

	var validationErr *validation.Error
	require.ErrorAs(t, err, &validationErr)
}

func TestUseCase_Execute_PatchAndSet(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")

	serverVersion, err := version.NewVersion("6.1.0")
	require.NoError(t, err)

	ai := &fakeAI{
		version: serverVersion,
		current: domainsettings.ScanSettings{
			ProjectName: "project",
			Priority:    domainsettings.PriorityLow,
		},
	}

	cfg := newTestConfig(projectID)
	uc, err := settings.NewUseCase(ai, fakeCLI{}, cfg)
	require.NoError(t, err)

	priority := domainsettings.PriorityCritical
	patch := domainsettings.ProjectSettingsPatch{Priority: &priority}

	err = uc.Execute(ctx, patch)
	require.NoError(t, err)
	require.True(t, ai.setCalled)
	require.Equal(t, domainsettings.PriorityCritical, ai.lastSetSettings.Priority)
	require.Equal(t, "project", ai.lastSetSettings.ProjectName)
}

func newTestConfig(projectID uuid.UUID) *config.Config {
	cfg := &config.Config{}
	_ = cfg.UpdateProjectId(projectID.String())

	return cfg
}
