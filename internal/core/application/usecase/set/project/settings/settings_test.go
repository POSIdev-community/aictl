package settings

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/application/usecase/set/project/settings/mocks"
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
)

var (
	okAIProj       = `{"GoSettings": {"CustomParameters": "+z"}}`
	emptySettings  = settings.ScanSettings{}
	filledSettings = settings.ScanSettings{
		ProjectName: "test",
		Languages:   []string{"go", "java"},
		GoSettings: settings.GoSettings{
			LaunchParameters: "-v",
		},
		JavaSettings: settings.JavaSettings{
			LaunchParameters: "-v",
		},
	}
)

func TestUseCase_Execute(t *testing.T) {
	t.Parallel()

	t.Run("update default settings set for project", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		projectID := uuid.New()
		scanID := uuid.New()

		updatedSettings := filledSettings
		updatedSettings.GoSettings.LaunchParameters = "+z"

		aiAdapter := mocks.NewAI(t)
		aiAdapter.On("GetScanAiproj", ctx, projectID, scanID).Return(okAIProj, nil).Once()
		aiAdapter.On("GetDefaultSettings", ctx).Return(filledSettings, nil).Once()
		aiAdapter.On("SetProjectSettings", ctx, projectID, &updatedSettings).Return(nil).Once()

		cliAdapter := mocks.NewCLI(t)

		uc, err := NewUseCase(aiAdapter, cliAdapter)
		require.NoError(t, err)

		require.NoError(t, uc.Execute(ctx, projectID, scanID))
	})

	t.Run("empty default settings", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		projectID := uuid.New()
		scanID := uuid.New()

		updatedSettings := emptySettings
		updatedSettings.GoSettings.LaunchParameters = "+z"

		aiAdapter := mocks.NewAI(t)
		aiAdapter.On("GetScanAiproj", ctx, projectID, scanID).Return(okAIProj, nil).Once()
		aiAdapter.On("GetDefaultSettings", ctx).Return(emptySettings, nil).Once()
		aiAdapter.On("SetProjectSettings", ctx, projectID, &updatedSettings).Return(nil).Once()

		cliAdapter := mocks.NewCLI(t)

		uc, err := NewUseCase(aiAdapter, cliAdapter)
		require.NoError(t, err)

		require.NoError(t, uc.Execute(ctx, projectID, scanID))
	})

	t.Run("bad aiproj", func(t *testing.T) {
		t.Parallel()

		const badAIProj = `{"Trash": "foo-bar"}`

		ctx := t.Context()
		projectID := uuid.New()
		scanID := uuid.New()

		aiAdapter := mocks.NewAI(t)
		aiAdapter.On("GetScanAiproj", ctx, projectID, scanID).Return(badAIProj, nil).Once()

		cliAdapter := mocks.NewCLI(t)

		uc, err := NewUseCase(aiAdapter, cliAdapter)
		require.NoError(t, err)

		err = uc.Execute(ctx, projectID, scanID)
		require.Error(t, err)
		assert.ErrorContains(t, err, "unknown field")
	})
}
