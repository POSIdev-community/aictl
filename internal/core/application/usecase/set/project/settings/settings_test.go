package settings

import (
	"testing"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/application/usecase/set/project/settings/mocks"
	"github.com/POSIdev-community/aictl/internal/core/domain/aiproj"
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

		updatedSettings := filledSettings
		updatedSettings.GoSettings.LaunchParameters = "+z"

		aiAdapter := mocks.NewAI(t)
		aiAdapter.On("InitializeWithRetry", ctx).Return(nil).Once()
		aiAdapter.On("GetDefaultSettings", ctx).Return(filledSettings, nil).Once()
		aiAdapter.On("SetProjectSettings", ctx, projectID, &updatedSettings).Return(nil).Once()

		cliAdapter := mocks.NewCLI(t)

		cfg := config.NewConfig(config.Uri{}, "", true, projectID, uuid.New())

		uc, err := NewUseCase(aiAdapter, cliAdapter, cfg)
		require.NoError(t, err)

		aiProj, err := aiproj.FromString(okAIProj)
		require.NoError(t, err)

		require.NoError(t, uc.Execute(ctx, &aiProj))
	})

	t.Run("empty default settings", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		projectID := uuid.New()

		updatedSettings := emptySettings
		updatedSettings.GoSettings.LaunchParameters = "+z"

		aiAdapter := mocks.NewAI(t)
		aiAdapter.On("InitializeWithRetry", ctx).Return(nil).Once()
		aiAdapter.On("GetDefaultSettings", ctx).Return(emptySettings, nil).Once()
		aiAdapter.On("SetProjectSettings", ctx, projectID, &updatedSettings).Return(nil).Once()

		cliAdapter := mocks.NewCLI(t)

		cfg := config.NewConfig(config.Uri{}, "", true, projectID, uuid.New())

		uc, err := NewUseCase(aiAdapter, cliAdapter, cfg)
		require.NoError(t, err)

		aiProj, err := aiproj.FromString(okAIProj)
		require.NoError(t, err)

		require.NoError(t, uc.Execute(ctx, &aiProj))
	})

}
