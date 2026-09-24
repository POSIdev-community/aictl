package settings

import (
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	domainproject "github.com/POSIdev-community/aictl/internal/core/domain/project"
	domainsettings "github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/internal/core/domain/version"
)

var okAIProj = []byte(`{
	"Version": "1.9",
	"ProjectName": "demo",
	"ProgrammingLanguages": ["Go"],
	"ScanModules": ["StaticCodeAnalysis"],
	"GoSettings": {"CustomParameters": "+z"}
}`)

var aiprojWithPolicyCheck = []byte(`{
	"Version": "1.9",
	"ProjectName": "demo",
	"ProgrammingLanguages": ["Go"],
	"ScanModules": ["StaticCodeAnalysis"],
	"UseSecurityPolicies": true
}`)

func TestUseCase_Execute(t *testing.T) {
	t.Parallel()

	t.Run("update default settings set for project", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		projectID := uuid.New()
		serverVersion, err := version.NewVersion("6.1.0")
		require.NoError(t, err)

		aiAdapter := NewMockAI(t)
		aiAdapter.On("InitializeWithRetry", ctx).Return(nil).Once()
		aiAdapter.On("GetProject", ctx, projectID).Return(&domainproject.Project{Id: projectID, Name: "demo", Type: domainproject.TypeSource}, nil).Once()
		aiAdapter.On("GetVersion", ctx).Return(serverVersion, nil).Once()
		aiAdapter.On("GetDefaultSettings", ctx).Return(domainsettings.ScanSettings{
			ProjectName: "test",
			Languages:   []string{"go", "java"},
			JavaSettings: domainsettings.JavaSettings{
				LaunchParameters: "-v",
			},
		}, nil).Once()
		aiAdapter.On("GetProjectSettings", ctx, projectID).Return(domainsettings.ScanSettings{
			Priority: domainsettings.PriorityHigh,
			PreferredAgentsSettings: domainsettings.PreferredAgentsSettings{
				PreferredAgentsOnly: true,
			},
		}, nil).Once()
		aiAdapter.On("SetProjectSettings", ctx, projectID, mock.MatchedBy(func(settings *domainsettings.ScanSettings) bool {
			return settings != nil &&
				settings.ProjectName == "demo" &&
				len(settings.Languages) == 1 &&
				settings.Languages[0] == "Go" &&
				settings.WhiteBoxSettings.StaticCodeAnalysisEnabled &&
				settings.GoSettings.LaunchParameters == "+z" &&
				settings.Priority == domainsettings.PriorityHigh &&
				settings.PreferredAgentsSettings.PreferredAgentsOnly
		})).Return(nil).Once()

		cliAdapter := NewMockCLI(t)

		cfg := config.NewConfig(config.Uri{}, "", true, projectID, uuid.New())

		uc, err := NewUseCase(aiAdapter, cliAdapter, cfg)
		require.NoError(t, err)

		require.NoError(t, uc.Execute(ctx, okAIProj))
	})

	t.Run("empty default settings", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		projectID := uuid.New()
		serverVersion, err := version.NewVersion("6.1.0")
		require.NoError(t, err)

		aiAdapter := NewMockAI(t)
		aiAdapter.On("InitializeWithRetry", ctx).Return(nil).Once()
		aiAdapter.On("GetProject", ctx, projectID).Return(&domainproject.Project{Id: projectID, Name: "demo", Type: domainproject.TypeSource}, nil).Once()
		aiAdapter.On("GetVersion", ctx).Return(serverVersion, nil).Once()
		aiAdapter.On("GetDefaultSettings", ctx).Return(domainsettings.ScanSettings{}, nil).Once()
		aiAdapter.On("GetProjectSettings", ctx, projectID).Return(domainsettings.ScanSettings{}, nil).Once()
		aiAdapter.On("SetProjectSettings", ctx, projectID, mock.MatchedBy(func(settings *domainsettings.ScanSettings) bool {
			return settings != nil && settings.GoSettings.LaunchParameters == "+z"
		})).Return(nil).Once()

		cliAdapter := NewMockCLI(t)

		cfg := config.NewConfig(config.Uri{}, "", true, projectID, uuid.New())

		uc, err := NewUseCase(aiAdapter, cliAdapter, cfg)
		require.NoError(t, err)

		require.NoError(t, uc.Execute(ctx, okAIProj))
	})

	t.Run("applies UseSecurityPolicies from aiproj", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		projectID := uuid.New()
		serverVersion, err := version.NewVersion("6.1.0")
		require.NoError(t, err)

		aiAdapter := NewMockAI(t)
		aiAdapter.On("InitializeWithRetry", ctx).Return(nil).Once()
		aiAdapter.On("GetProject", ctx, projectID).Return(&domainproject.Project{Id: projectID, Name: "demo", Type: domainproject.TypeSource}, nil).Once()
		aiAdapter.On("GetVersion", ctx).Return(serverVersion, nil).Once()
		aiAdapter.On("GetDefaultSettings", ctx).Return(domainsettings.ScanSettings{}, nil).Once()
		aiAdapter.On("GetProjectSettings", ctx, projectID).Return(domainsettings.ScanSettings{}, nil).Once()
		aiAdapter.On("SetProjectSettings", ctx, projectID, mock.Anything).Return(nil).Once()
		aiAdapter.On("GetProjectPolicies", ctx, projectID).Return(io.NopCloser(strings.NewReader(
			`{"checkSecurityPoliciesAccordance":false,"securityPolicies":"[{\"CountToActualize\":1}]"}`,
		)), nil).Once()
		aiAdapter.On("SetProjectPolicies", ctx, projectID, mock.MatchedBy(func(raw []byte) bool {
			var model struct {
				Check    bool   `json:"checkSecurityPoliciesAccordance"`
				Policies string `json:"securityPolicies"`
			}
			if err := json.Unmarshal(raw, &model); err != nil {
				return false
			}
			return model.Check && model.Policies == `[{"CountToActualize":1}]`
		})).Return(nil).Once()

		cliAdapter := NewMockCLI(t)
		cfg := config.NewConfig(config.Uri{}, "", true, projectID, uuid.New())

		uc, err := NewUseCase(aiAdapter, cliAdapter, cfg)
		require.NoError(t, err)
		require.NoError(t, uc.Execute(ctx, aiprojWithPolicyCheck))
	})
}
