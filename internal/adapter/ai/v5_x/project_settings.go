package v5_x

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanagent"
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/pkg/clientai/v5_x"
)

func (a *ClientAI5x) GetProjectSettings(ctx context.Context, projectId uuid.UUID) (settings.ScanSettings, error) {
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

	return mapProjectSettingsFromModel(*res.JSON200), nil
}

func (a *ClientAI5x) GetScanAgents(ctx context.Context) ([]scanagent.ScanAgent, error) {
	response, err := a.GetApiScanAgentsWithResponse(ctx, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get scan agents request: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	if err = CheckResponseByModel(statusCode, body, nil); err != nil {
		return nil, fmt.Errorf("ai adapter get scan agents: %w", err)
	}

	if response.JSON200 == nil {
		return []scanagent.ScanAgent{}, nil
	}

	models := *response.JSON200
	agents := make([]scanagent.ScanAgent, 0, len(models))
	for _, model := range models {
		if model.Id == nil {
			continue
		}

		status := ""
		if model.StatusType != nil {
			status = string(*model.StatusType)
		}

		agents = append(agents, scanagent.ScanAgent{
			Id:              uuid.UUID(*model.Id),
			Name:            common.GetOrDefault(model.Name, ""),
			Status:          status,
			Version:         common.GetOrDefault(model.Version, ""),
			OperatingSystem: common.GetOrDefault(model.OperatingSystem, ""),
		})
	}

	return agents, nil
}

func mapProjectSettingsFromModel(model v5_x.ProjectSettingsModel) settings.ScanSettings {
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

	if model.WhiteBoxSettings != nil {
		wb := model.WhiteBoxSettings
		result.WhiteBoxSettings = settings.WhiteBoxSettings{
			StaticCodeAnalysisEnabled:            common.GetOrDefault(wb.StaticCodeAnalysisEnabled, false),
			PatternMatchingEnabled:               common.GetOrDefault(wb.PatternMatchingEnabled, false),
			SearchForVulnerableComponentsEnabled: common.GetOrDefault(wb.SearchForVulnerableComponentsEnabled, false),
			SearchForConfigurationFlawsEnabled:   common.GetOrDefault(wb.SearchForConfigurationFlawsEnabled, false),
			SearchWithScaEnabled:                 common.GetOrDefault(wb.SearchWithScaEnabled, false),
		}
	}

	if model.DotNetSettings != nil {
		dn := model.DotNetSettings
		result.DotNetSettings = settings.DotNetSettings{
			ProjectType:                           string(common.GetOrDefault(dn.ProjectType, "")),
			SolutionFile:                          common.GetOrDefault(dn.SolutionFile, ""),
			WebSiteFolder:                         common.GetOrDefault(dn.WebSiteFolder, ""),
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

	if model.PygrepSettings != nil {
		pg := model.PygrepSettings
		result.PygrepSettings = settings.PygrepSettings{
			RulesDirPath:     common.GetOrDefault(pg.RulesDirPath, ""),
			LaunchParameters: common.GetOrDefault(pg.LaunchParameters, ""),
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
