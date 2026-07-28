package v6_x

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanagent"
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/pkg/clientai/v6_x"
)

func (a *ClientAI61) GetProjectSettings(ctx context.Context, projectId uuid.UUID) (settings.ScanSettings, error) {
	res, err := a.GetApiProjectsProjectIdSettingsWithResponse(ctx, projectId, a.AddJWTToHeader)
	if err != nil {
		return settings.ScanSettings{}, fmt.Errorf("get project settings request: %w", err)
	}

	statusCode := res.StatusCode()
	responseBody := string(res.Body)
	if err = CheckResponseByModel(statusCode, responseBody, res.JSON400); err != nil {
		return settings.ScanSettings{}, fmt.Errorf("get project settings: %w", err)
	}

	if res.JSON200 == nil {
		return settings.ScanSettings{}, fmt.Errorf("get project settings: empty response")
	}

	result := mapProjectSettingsFromModel(*res.JSON200)

	preferredRes, err := a.GetApiProjectsProjectIdPreferredAgentsSettingsWithResponse(ctx, projectId, a.AddJWTToHeader)
	if err != nil {
		return settings.ScanSettings{}, fmt.Errorf("get preferred agents settings request: %w", err)
	}

	statusCode = preferredRes.StatusCode()
	responseBody = string(preferredRes.Body)
	if err = CheckResponseByModel(statusCode, responseBody, preferredRes.JSON400); err != nil {
		return settings.ScanSettings{}, fmt.Errorf("get preferred agents settings: %w", err)
	}

	if preferredRes.JSON200 != nil {
		result.PreferredAgentsSettings = mapPreferredAgentsFromModel(*preferredRes.JSON200)
	}

	return result, nil
}

func (a *ClientAI61) GetScanAgents(ctx context.Context) ([]scanagent.ScanAgent, error) {
	response, err := a.GetAllWithResponse(ctx, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get scan agents request: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	if err = CheckResponseByModel(statusCode, body, response.JSON400); err != nil {
		return nil, fmt.Errorf("ai adapter get scan agents: %w", err)
	}

	if response.JSON200 == nil {
		return []scanagent.ScanAgent{}, nil
	}

	models := *response.JSON200
	agents := make([]scanagent.ScanAgent, 0, len(models))
	for _, model := range models {
		agents = append(agents, mapAgentToScanAgent(model))
	}

	return agents, nil
}

func (a *ClientAI61) putPreferredAgentsSettings(ctx context.Context, projectId uuid.UUID, s *settings.ScanSettings) error {
	preferredAgents := make([]openapi_types.UUID, len(s.PreferredAgentsSettings.PreferredAgents))
	for i, id := range s.PreferredAgentsSettings.PreferredAgents {
		preferredAgents[i] = openapi_types.UUID(id)
	}

	preferredOnly := s.PreferredAgentsSettings.PreferredAgentsOnly
	body := v6_x.PreferredAgentsSettingsModel{
		PreferredAgents:     &preferredAgents,
		PreferredAgentsOnly: &preferredOnly,
	}

	res, err := a.PutApiProjectsProjectIdPreferredAgentsSettingsWithResponse(ctx, projectId, body, a.AddJWTToHeader)
	if err != nil {
		return fmt.Errorf("put preferred agents settings request: %w", err)
	}

	statusCode := res.StatusCode()
	responseBody := string(res.Body)
	if err = CheckResponseByModel(statusCode, responseBody, res.JSON400); err != nil {
		return fmt.Errorf("put preferred agents settings: %w", err)
	}

	return nil
}

func priorityFromSettings(s *settings.ScanSettings) v6_x.Priority {
	if s.Priority == "" {
		return v6_x.PriorityLow
	}

	return v6_x.Priority(s.Priority)
}

func mapPreferredAgentsFromModel(model v6_x.PreferredAgentsSettingsModel) settings.PreferredAgentsSettings {
	result := settings.PreferredAgentsSettings{
		PreferredAgentsOnly: common.GetOrDefault(model.PreferredAgentsOnly, false),
	}

	if model.PreferredAgents != nil {
		agents := make([]uuid.UUID, len(*model.PreferredAgents))
		for i, id := range *model.PreferredAgents {
			agents[i] = uuid.UUID(id)
		}
		result.PreferredAgents = agents
	}

	return result
}

func mapAgentToScanAgent(model v6_x.Agent) scanagent.ScanAgent {
	status := string(model.Status)
	if model.Disabled {
		status += " (disabled)"
	}
	if model.Stuck {
		status += " (stuck)"
	}

	return scanagent.ScanAgent{
		Id:              uuid.UUID(model.Id),
		Name:            model.Name,
		Status:          status,
		Version:         model.Version,
		OperatingSystem: string(model.OperatingSystem),
		Disabled:        model.Disabled,
		Stuck:           model.Stuck,
	}
}

func mapProjectSettingsFromModel(model v6_x.ProjectSettingsModel) settings.ScanSettings {
	result := settings.ScanSettings{
		ProjectName: common.GetOrDefault(model.ProjectName, ""),
		Languages: func() []string {
			if model.Languages == nil {
				return nil
			}

			res := make([]string, len(*model.Languages))
			for i := range *model.Languages {
				res[i] = string((*model.Languages)[i])
			}

			return res
		}(),
		SkipGitIgnoreFiles: common.GetOrDefault(model.SkipGitIgnoreFiles, false),
	}

	if model.Priority != nil {
		result.Priority = settings.Priority(*model.Priority)
	}

	if model.WhiteBoxSettings != nil {
		wb := model.WhiteBoxSettings
		result.WhiteBoxSettings = settings.WhiteBoxSettings{
			StaticCodeAnalysisEnabled:            common.GetOrDefault(wb.StaticCodeAnalysisEnabled, false),
			PatternMatchingEnabled:               common.GetOrDefault(wb.PatternMatchingEnabled, false),
			SearchForVulnerableComponentsEnabled: common.GetOrDefault(wb.SearchForVulnerableComponentsEnabled, false),
			SearchForConfigurationFlawsEnabled:   common.GetOrDefault(wb.SearchForConfigurationFlawsEnabled, false),
			SearchWithScaEnabled:                 common.GetOrDefault(wb.SearchWithScaEnabled, false),
			SecretDetectionEnabled:               common.GetOrDefault(wb.SecretDetectionEnabled, false),
			SearchForMaliciousCodeEnabled:        common.GetOrDefault(wb.SearchForMaliciousCodeEnabled, false),
		}
	}

	if model.DotNetSettings != nil {
		dn := model.DotNetSettings
		result.DotNetSettings = settings.DotNetSettings{
			ProjectType:                           string(common.GetOrDefault(dn.ProjectType, "")),
			SolutionFile:                          common.GetOrDefault(dn.SolutionFile, ""),
			LaunchParameters:                      common.GetOrDefault(dn.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: common.GetOrDefault(dn.UseAvailablePublicAndProtectedMethods, false),
			DownloadDependencies:                  common.GetOrDefault(dn.DownloadDependencies, false),
		}
	}

	if model.GoSettings != nil {
		gs := model.GoSettings
		result.GoSettings = settings.GoSettings{
			LaunchParameters:                      common.GetOrDefault(gs.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: common.GetOrDefault(gs.UseAvailablePublicAndProtectedMethods, false),
		}
	}

	if model.JavaScriptSettings != nil {
		js := model.JavaScriptSettings
		result.JavaScriptSettings = settings.JavaScriptSettings{
			LaunchParameters:                      common.GetOrDefault(js.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: common.GetOrDefault(js.UseAvailablePublicAndProtectedMethods, false),
			DownloadDependencies:                  common.GetOrDefault(js.DownloadDependencies, false),
			UseTaintAnalysis:                      common.GetOrDefault(js.UseTaintAnalysis, false),
			UseJsaAnalysis:                        common.GetOrDefault(js.UseJsaAnalysis, false),
		}
	}

	if model.JavaSettings != nil {
		js := model.JavaSettings
		result.JavaSettings = settings.JavaSettings{
			Parameters:                            common.GetOrDefault(js.Parameters, ""),
			UnpackUserPackages:                    common.GetOrDefault(js.UnpackUserPackages, false),
			UserPackagePrefixes:                   common.GetOrDefault(js.UserPackagePrefixes, ""),
			Version:                               string(common.GetOrDefault(js.Version, "")),
			LaunchParameters:                      common.GetOrDefault(js.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: common.GetOrDefault(js.UseAvailablePublicAndProtectedMethods, false),
			DownloadDependencies:                  common.GetOrDefault(js.DownloadDependencies, false),
			DependenciesPath:                      common.GetOrDefault(js.DependenciesPath, ""),
		}
	}

	if model.PhpSettings != nil {
		ps := model.PhpSettings
		result.PhpSettings = settings.PhpSettings{
			LaunchParameters:                      common.GetOrDefault(ps.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: common.GetOrDefault(ps.UseAvailablePublicAndProtectedMethods, false),
			DownloadDependencies:                  common.GetOrDefault(ps.DownloadDependencies, false),
		}
	}

	if model.PmTaintSettings != nil {
		pm := model.PmTaintSettings
		result.PmTaintSettings = settings.PmTaintSettings{
			LaunchParameters:                      common.GetOrDefault(pm.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: common.GetOrDefault(pm.UseAvailablePublicAndProtectedMethods, false),
		}
	}

	if model.PythonSettings != nil {
		ps := model.PythonSettings
		result.PythonSettings = settings.PythonSettings{
			LaunchParameters:                      common.GetOrDefault(ps.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: common.GetOrDefault(ps.UseAvailablePublicAndProtectedMethods, false),
			DownloadDependencies:                  common.GetOrDefault(ps.DownloadDependencies, false),
			DependenciesPath:                      common.GetOrDefault(ps.DependenciesPath, ""),
		}
	}

	if model.RubySettings != nil {
		rs := model.RubySettings
		result.RubySettings = settings.RubySettings{
			LaunchParameters:                      common.GetOrDefault(rs.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: common.GetOrDefault(rs.UseAvailablePublicAndProtectedMethods, false),
		}
	}

	if model.ScaSettings != nil {
		ss := model.ScaSettings
		result.ScaSettings = settings.ScaSettings{
			LaunchParameters:       common.GetOrDefault(ss.LaunchParameters, ""),
			BuildDependenciesGraph: common.GetOrDefault(ss.BuildDependenciesGraph, false),
		}
	}

	return result
}
