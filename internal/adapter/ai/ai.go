package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/branch"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/scan"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanstage"
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/internal/core/port"
	. "github.com/POSIdev-community/aictl/pkg/clientai"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/POSIdev-community/aictl/pkg/logger"
)

var _ port.Ai = &Adapter{}

func NewAdapter(ctx context.Context, cfg *config.Config) (*Adapter, error) {
	aiClient, err := NewAiClient(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("new adapter: %w", err)
	}

	return &Adapter{
		aiClient,
	}, nil
}

type Adapter struct {
	aiClient *AiClient
}

func (a *Adapter) GetDefaultSettings(ctx context.Context) (settings.ScanSettings, error) {
	defaultSettings, err := a.aiClient.GetApiProjectsDefaultSettingsWithResponse(ctx)
	if err != nil {
		return settings.ScanSettings{}, fmt.Errorf("get projects default settings: %w", err)
	}

	if defaultSettings == nil || defaultSettings.JSON200 == nil {
		return settings.ScanSettings{}, fmt.Errorf("empty projects default settings")
	}

	result := *defaultSettings.JSON200

	return settings.ScanSettings{
		ProjectName: getOrDefault(result.Name, ""),
		Languages: func() []string {
			if result.Languages == nil {
				return nil
			}

			res := make([]string, len(*result.Languages))
			for i := range *result.Languages {
				res[i] = string((*result.Languages)[i])
			}

			return res
		}(),
		WhiteBoxSettings: settings.WhiteBoxSettings{
			StaticCodeAnalysisEnabled:            getOrDefault(result.WhiteBox.StaticCodeAnalysisEnabled, false),
			PatternMatchingEnabled:               getOrDefault(result.WhiteBox.PatternMatchingEnabled, false),
			SearchForVulnerableComponentsEnabled: getOrDefault(result.WhiteBox.SearchForVulnerableComponentsEnabled, false),
			SearchForConfigurationFlawsEnabled:   getOrDefault(result.WhiteBox.SearchForConfigurationFlawsEnabled, false),
			SearchWithScaEnabled:                 getOrDefault(result.WhiteBox.SearchWithScaEnabled, false),
		},
		DotNetSettings: settings.DotNetSettings{
			ProjectType:                           string(getOrDefault(result.DotNetSettings.ProjectType, "")),
			SolutionFile:                          getOrDefault(result.DotNetSettings.SolutionFile, ""),
			WebSiteFolder:                         getOrDefault(result.DotNetSettings.WebSiteFolder, ""),
			LaunchParameters:                      getOrDefault(result.DotNetSettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: getOrDefault(result.DotNetSettings.UseAvailablePublicAndProtectedMethods, false),
			DownloadDependencies:                  getOrDefault(result.DotNetSettings.DownloadDependencies, false),
		},
		GoSettings: settings.GoSettings{
			LaunchParameters:                      getOrDefault(result.GoSettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: getOrDefault(result.GoSettings.UseAvailablePublicAndProtectedMethods, false),
		},
		JavaScriptSettings: settings.JavaScriptSettings{
			LaunchParameters:                      getOrDefault(result.JavaScriptSettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: getOrDefault(result.JavaScriptSettings.UseAvailablePublicAndProtectedMethods, false),
			DownloadDependencies:                  getOrDefault(result.JavaScriptSettings.DownloadDependencies, false),
			UseTaintAnalysis:                      getOrDefault(result.JavaScriptSettings.UseTaintAnalysis, false),
			UseJsaAnalysis:                        getOrDefault(result.JavaScriptSettings.UseJsaAnalysis, false),
		},
		JavaSettings: settings.JavaSettings{
			Parameters:                            getOrDefault(result.JavaSettings.Parameters, ""),
			UnpackUserPackages:                    getOrDefault(result.JavaSettings.UnpackUserPackages, false),
			UserPackagePrefixes:                   getOrDefault(result.JavaSettings.UserPackagePrefixes, ""),
			Version:                               string(getOrDefault(result.JavaSettings.Version, "")),
			LaunchParameters:                      getOrDefault(result.JavaSettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: getOrDefault(result.JavaSettings.UseAvailablePublicAndProtectedMethods, false),
			DownloadDependencies:                  getOrDefault(result.JavaSettings.DownloadDependencies, false),
			DependenciesPath:                      getOrDefault(result.JavaSettings.DependenciesPath, ""),
		},
		PhpSettings: settings.PhpSettings{
			LaunchParameters:                      getOrDefault(result.PhpSettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: getOrDefault(result.PhpSettings.UseAvailablePublicAndProtectedMethods, false),
			DownloadDependencies:                  getOrDefault(result.PhpSettings.DownloadDependencies, false),
		},
		PmTaintSettings: settings.PmTaintSettings{
			LaunchParameters:                      getOrDefault(result.PmTaintSettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: getOrDefault(result.PmTaintSettings.UseAvailablePublicAndProtectedMethods, false),
		},
		PythonSettings: settings.PythonSettings{
			LaunchParameters:                      getOrDefault(result.PythonSettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: getOrDefault(result.PythonSettings.UseAvailablePublicAndProtectedMethods, false),
			DownloadDependencies:                  getOrDefault(result.PythonSettings.DownloadDependencies, false),
			DependenciesPath:                      getOrDefault(result.PythonSettings.DependenciesPath, ""),
		},
		RubySettings: settings.RubySettings{
			LaunchParameters:                      getOrDefault(result.RubySettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: getOrDefault(result.RubySettings.UseAvailablePublicAndProtectedMethods, false),
		},
		PygrepSettings: settings.PygrepSettings{
			RulesDirPath:     getOrDefault(result.PygrepSettings.RulesDirPath, ""),
			LaunchParameters: getOrDefault(result.PygrepSettings.LaunchParameters, ""),
		},
		ScaSettings: settings.ScaSettings{
			LaunchParameters:       getOrDefault(result.ScaSettings.LaunchParameters, ""),
			BuildDependenciesGraph: getOrDefault(result.ScaSettings.BuildDependenciesGraph, false),
		},
	}, err
}

func (a *Adapter) SetProjectSettings(ctx context.Context, projectId uuid.UUID, settings *settings.ScanSettings) error {
	if settings == nil {
		return nil
	}

	projectSettings := PutApiProjectsProjectIdSettingsJSONRequestBody{
		ProjectName: &settings.ProjectName,
		Languages: func() *[]LegacyProgrammingLanguageGroup {
			if settings.Languages == nil {
				return nil
			}
			res := make([]LegacyProgrammingLanguageGroup, len(settings.Languages))
			for i := range settings.Languages {
				res[i] = LegacyProgrammingLanguageGroup(settings.Languages[i])
			}
			return &res
		}(),
		WhiteBoxSettings: &WhiteBoxSettingsModel{
			StaticCodeAnalysisEnabled:            &settings.WhiteBoxSettings.StaticCodeAnalysisEnabled,
			PatternMatchingEnabled:               &settings.WhiteBoxSettings.PatternMatchingEnabled,
			SearchForVulnerableComponentsEnabled: &settings.WhiteBoxSettings.SearchForVulnerableComponentsEnabled,
			SearchForConfigurationFlawsEnabled:   &settings.WhiteBoxSettings.SearchForConfigurationFlawsEnabled,
			SearchWithScaEnabled:                 &settings.WhiteBoxSettings.SearchWithScaEnabled,
		},
		DotNetSettings: &DotNetSettingsModel{
			ProjectType:                           reference(DotNetProjectType(settings.DotNetSettings.ProjectType)),
			SolutionFile:                          &settings.DotNetSettings.SolutionFile,
			WebSiteFolder:                         &settings.DotNetSettings.WebSiteFolder,
			LaunchParameters:                      &settings.DotNetSettings.LaunchParameters,
			UseAvailablePublicAndProtectedMethods: &settings.DotNetSettings.UseAvailablePublicAndProtectedMethods,
			DownloadDependencies:                  &settings.DotNetSettings.DownloadDependencies,
		},
		GoSettings: &GoSettingsModel{
			LaunchParameters:                      &settings.GoSettings.LaunchParameters,
			UseAvailablePublicAndProtectedMethods: &settings.GoSettings.UseAvailablePublicAndProtectedMethods,
		},
		JavaScriptSettings: &JavaScriptSettingsModel{
			LaunchParameters:                      &settings.JavaScriptSettings.LaunchParameters,
			UseAvailablePublicAndProtectedMethods: &settings.JavaScriptSettings.UseAvailablePublicAndProtectedMethods,
			DownloadDependencies:                  &settings.JavaScriptSettings.DownloadDependencies,
			UseTaintAnalysis:                      &settings.JavaScriptSettings.UseTaintAnalysis,
			UseJsaAnalysis:                        &settings.JavaScriptSettings.UseJsaAnalysis,
		},
		JavaSettings: &JavaSettingsModel{
			Parameters:                            &settings.JavaSettings.Parameters,
			UnpackUserPackages:                    &settings.JavaSettings.UnpackUserPackages,
			UserPackagePrefixes:                   &settings.JavaSettings.UserPackagePrefixes,
			Version:                               reference(JavaVersions(settings.JavaSettings.Version)),
			LaunchParameters:                      &settings.JavaSettings.LaunchParameters,
			UseAvailablePublicAndProtectedMethods: &settings.JavaSettings.UseAvailablePublicAndProtectedMethods,
			DownloadDependencies:                  &settings.JavaSettings.DownloadDependencies,
			DependenciesPath:                      &settings.JavaSettings.DependenciesPath,
		},
		PhpSettings: &PhpSettingsModel{
			LaunchParameters:                      &settings.PhpSettings.LaunchParameters,
			UseAvailablePublicAndProtectedMethods: &settings.PhpSettings.UseAvailablePublicAndProtectedMethods,
			DownloadDependencies:                  &settings.PhpSettings.DownloadDependencies,
		},
		PmTaintSettings: &PmTaintBaseSettingsModel{
			LaunchParameters:                      &settings.PmTaintSettings.LaunchParameters,
			UseAvailablePublicAndProtectedMethods: &settings.PmTaintSettings.UseAvailablePublicAndProtectedMethods,
		},
		PythonSettings: &PythonSettingsModel{
			LaunchParameters:                      &settings.PythonSettings.LaunchParameters,
			UseAvailablePublicAndProtectedMethods: &settings.PythonSettings.UseAvailablePublicAndProtectedMethods,
			DownloadDependencies:                  &settings.PythonSettings.DownloadDependencies,
			DependenciesPath:                      &settings.PythonSettings.DependenciesPath,
		},
		RubySettings: &RubySettingsModel{
			LaunchParameters:                      &settings.RubySettings.LaunchParameters,
			UseAvailablePublicAndProtectedMethods: &settings.RubySettings.UseAvailablePublicAndProtectedMethods,
		},
		PygrepSettings: &PygrepSettingsModel{
			RulesDirPath:     &settings.PygrepSettings.RulesDirPath,
			LaunchParameters: &settings.PygrepSettings.LaunchParameters,
		},
		ScaSettings: &ScaSettingsModel{
			LaunchParameters:       &settings.ScaSettings.LaunchParameters,
			BuildDependenciesGraph: &settings.ScaSettings.BuildDependenciesGraph,
		},
	}

	res, err := a.aiClient.PutApiProjectsProjectIdSettingsWithResponse(ctx, projectId, projectSettings, a.aiClient.AddJwtToHeader)
	if err != nil {
		return fmt.Errorf("put project settings: %w", err)
	}

	statusCode := res.StatusCode()
	responseBody := string(res.Body)
	errorModel := res.JSON400
	if err = CheckResponseByModel(statusCode, responseBody, errorModel); err != nil {
		return fmt.Errorf("put project settings: %w", err)
	}

	return nil
}

func (a *Adapter) CreateBranch(ctx context.Context, projectId uuid.UUID, branchName, scanTarget string) (*uuid.UUID, error) {
	if scanTarget == "" {
		var err error
		scanTarget, err = createStubScanTarget()
		if err != nil {
			return nil, err
		}
	}

	archivePath, err := prepareArchive(scanTarget)
	if err != nil {
		return nil, err
	}

	body, contentType, err := prepareMultipartBody(archivePath, MultipartField{Key: "Name", Value: branchName})
	if err != nil {
		return nil, err
	}

	readCloser := io.NopCloser(body)

	response, err := a.aiClient.PostApiStoreProjectProjectIdBranchesArchiveWithBodyWithResponse(ctx, projectId, contentType, readCloser, a.aiClient.AddJwtToHeader)
	if err != nil {
		return nil, fmt.Errorf("create upload session response error: %w", err)
	}

	statusCode := response.StatusCode()
	responseBody := string(response.Body)
	errorModel := response.JSON400
	if err = CheckResponseByModel(statusCode, responseBody, errorModel); err != nil {
		return nil, err
	}

	id := string(response.Body)
	branchId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return &branchId, nil
}

func (a *Adapter) CreateProject(ctx context.Context, projectName string) (*uuid.UUID, error) {
	projectUrl := "http://localhost"

	patternMatchingEnabled := true
	searchForConfigurationFlawsEnabled := true
	searchForVulnerableComponentsEnabled := true
	searchWithScaEnabled := false
	staticCodeAnalysisEnabled := true

	projectBaseModel := PostApiProjectsBaseJSONRequestBody{
		Name:       &projectName,
		ProjectUrl: &projectUrl,
		WhiteBox: &WhiteBoxSettingsModel{
			PatternMatchingEnabled:               &patternMatchingEnabled,
			SearchForConfigurationFlawsEnabled:   &searchForConfigurationFlawsEnabled,
			SearchForVulnerableComponentsEnabled: &searchForVulnerableComponentsEnabled,
			SearchWithScaEnabled:                 &searchWithScaEnabled,
			StaticCodeAnalysisEnabled:            &staticCodeAnalysisEnabled,
		},
		Id: &uuid.UUID{},
		Languages: &[]LegacyProgrammingLanguageGroup{
			LegacyProgrammingLanguageGroupGo,
		},
	}

	createProjectResponse, err := a.aiClient.PostApiProjectsBaseWithResponse(ctx, projectBaseModel, a.aiClient.AddJwtToHeader)
	if err != nil {
		return nil, fmt.Errorf("create project request error: %w", err)
	}

	statusCode := createProjectResponse.StatusCode()
	body := string(createProjectResponse.Body)
	errorModel := createProjectResponse.JSON400
	if err = CheckResponseByModel(statusCode, body, errorModel); err != nil {
		return nil, err
	}

	projectId, err := uuid.Parse(body)
	if err != nil {
		return nil, err
	}

	return &projectId, nil
}

func (a *Adapter) DeleteProject(ctx context.Context, projectId uuid.UUID) error {
	response, err := a.aiClient.DeleteApiProjectsProjectId(ctx, projectId, a.aiClient.AddJwtToHeader)
	if err != nil {
		return fmt.Errorf("ai adapter delete project: %w", err)
	}

	if err = CheckResponse(response, "project"); err != nil {
		return fmt.Errorf("ai adapter delete project: %w", err)
	}

	return nil
}

func (a *Adapter) GetProjects(ctx context.Context) ([]project.Project, error) {
	log := logger.FromContext(ctx)

	log.Info("Send get projects request")

	response, err := a.aiClient.GetApiProjectsWithResponse(ctx, a.aiClient.AddJwtToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get projects request: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	if err = CheckResponseByModel(statusCode, body, nil); err != nil {
		return nil, fmt.Errorf("ai adapter get projects response: %w", err)
	}

	models := *response.JSON200
	projects := make([]project.Project, 0, len(models))

	for _, model := range models {
		if *model.ProjectType != Permanent {
			continue
		}

		p := project.NewProject(*model.Id, *model.Name)
		projects = append(projects, p)
	}

	return projects, nil
}

func (a *Adapter) GetProject(ctx context.Context, projectId uuid.UUID) (*project.Project, error) {
	response, err := a.aiClient.GetApiProjectsProjectIdWithResponse(ctx, projectId, a.aiClient.AddJwtToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get projects: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	errorModel := response.JSON400
	if err = CheckResponseByModel(statusCode, body, errorModel); err != nil {
		return nil, fmt.Errorf("ai adapter get project: %w", err)
	}

	model := response.JSON200
	p := project.NewProject(*model.Id, *model.Name)

	return &p, nil
}

func (a *Adapter) GetTemplateId(ctx context.Context, reportType string) (uuid.UUID, error) {
	localeId := "ru-Ru"
	params := GetApiReportsTemplatesTypeParams{
		LocaleId: &localeId,
	}

	var aiReportType ReportType
	switch reportType {
	case port.SarifReportType:
		aiReportType = ReportTypeSarif
	case port.PlainReportType:
		aiReportType = ReportTypePlainReport
	}

	response, err := a.aiClient.GetApiReportsTemplatesType(ctx, aiReportType, &params, a.aiClient.AddJwtToHeader)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("ai adapter get template id: %w", err)
	}

	if err = CheckResponse(response, "template"); err != nil {
		return uuid.UUID{}, fmt.Errorf("ai adapter get template id: %w", err)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	defer func() { _ = response.Body.Close() }()
	if err != nil {
		return uuid.UUID{}, err
	}

	if !strings.Contains(response.Header.Get("Content-Type"), "json") && response.StatusCode == 200 {
		return uuid.UUID{}, fmt.Errorf("ai adapter response not 200 and json")
	}

	type ReportTemplateSimpleModel struct {
		Id *uuid.UUID `json:"id,omitempty"`
	}

	var dest ReportTemplateSimpleModel
	if err := json.Unmarshal(bodyBytes, &dest); err != nil {
		return uuid.UUID{}, err
	}

	return *dest.Id, nil
}

func (a *Adapter) GetReport(ctx context.Context, projectId, scanResultId, templateId uuid.UUID) (io.ReadCloser, error) {

	localeId := "ru"
	includeComments := true
	includeDFD := true
	includeGlossary := true
	useFilters := false
	sessionId := uuid.New()

	model := ReportGenerateModel{
		LocaleId: &localeId,
		Parameters: &UserReportParametersModel{
			IncludeComments:  &includeComments,
			IncludeDFD:       &includeDFD,
			IncludeGlossary:  &includeGlossary,
			ReportTemplateId: &templateId,
			UseFilters:       &useFilters,
		},
		ProjectId:    &projectId,
		ScanResultId: &scanResultId,
		SessionId:    &sessionId,
	}

	response, err := a.aiClient.PostApiReportsGenerate(ctx, model, a.aiClient.AddJwtToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter generate report: %w", err)
	}

	if err = CheckResponse(response, "report"); err != nil {
		return nil, fmt.Errorf("ai adapter generate report: %w", err)
	}

	return response.Body, nil
}

func (a *Adapter) GetBranches(ctx context.Context, projectId uuid.UUID) ([]branch.Branch, error) {
	getBranchesResponse, err := a.aiClient.GetApiProjectsProjectIdBranchesWithResponse(ctx, projectId, a.aiClient.AddJwtToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get branch: %w", err)
	}

	statusCode := getBranchesResponse.StatusCode()
	body := string(getBranchesResponse.Body)
	errorModel := getBranchesResponse.JSON400
	if err = CheckResponseByModel(statusCode, body, errorModel); err != nil {
		return nil, fmt.Errorf("ai adapter get branches: %w", err)
	}

	branchModels := *getBranchesResponse.JSON200

	branches := make([]branch.Branch, len(branchModels))
	for i, model := range branchModels {
		branches[i] = branch.NewBranch(*model.Id, *model.Name, *model.Description, *model.IsWorking)
	}

	return branches, nil
}

func (a *Adapter) GetLastScan(ctx context.Context, branchId uuid.UUID) (*scan.Scan, error) {
	response, err := a.aiClient.GetApiBranchesBranchIdScanResultsLastWithResponse(ctx, branchId, a.aiClient.AddJwtToHeader)
	if err != nil {
		return nil, errs.NewNotFoundError("last scan result")
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	errorModel := response.JSON400
	if err = CheckResponseByModel(statusCode, body, errorModel); err != nil {
		return nil, fmt.Errorf("ai adapter get last scan result: %w", err)
	}

	model := response.JSON200
	scanResult := scan.NewScan(*model.Id, *model.SettingsId)

	return scanResult, nil
}

func (a *Adapter) GetScan(ctx context.Context, projectId, scanId uuid.UUID) (*scan.Scan, error) {
	response, err := a.aiClient.GetApiProjectsProjectIdScanResultsScanResultIdWithResponse(ctx, projectId, scanId, a.aiClient.AddJwtToHeader)
	if err != nil {
		return nil, errs.NewNotFoundError("get scan aiproj")
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	errorModel := response.JSON400
	if err = CheckResponseByModel(statusCode, body, errorModel); err != nil {
		return nil, fmt.Errorf("ai adapter get scan aiproj: %w", err)
	}

	model := response.JSON200

	return scan.NewScan(*model.Id, *model.SettingsId), nil
}

func (a *Adapter) GetScanAiproj(ctx context.Context, projectId, scanSettingsId uuid.UUID) (string, error) {
	response, err := a.aiClient.GetApiProjectsProjectIdScanSettingsScanSettingsIdAiprojWithResponse(ctx, projectId, scanSettingsId, a.aiClient.AddJwtToHeader)
	if err != nil {
		return "", errs.NewNotFoundError("get scan aiproj")
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	if err = CheckResponseByModel(statusCode, body, nil); err != nil {
		return "", fmt.Errorf("ai adapter get scan aiproj: %w", err)
	}

	return string(response.Body), nil
}

func (a *Adapter) GetScanStage(ctx context.Context, projectId, scanId uuid.UUID) (scanstage.ScanStage, error) {
	response, err := a.aiClient.GetApiProjectsProjectIdScanResultsScanResultIdProgressWithResponse(ctx, projectId, scanId, a.aiClient.AddJwtToHeader)
	if err != nil {
		return scanstage.ScanStage{}, errs.NewNotFoundError("scan result")
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	errorModel := response.JSON400
	if err = CheckResponseByModel(statusCode, body, errorModel); err != nil {
		return scanstage.ScanStage{}, fmt.Errorf("ai adapter get last scan result: %w", err)
	}

	model := *response.JSON200

	return scanstage.ScanStage{
		Value: *model.Value,
		Stage: string(*model.Stage),
	}, nil
}

func (a *Adapter) GetScans(ctx context.Context, branchId uuid.UUID) ([]scan.Scan, error) {
	response, err := a.aiClient.GetApiBranchesBranchIdScanResultsWithResponse(ctx, branchId, a.aiClient.AddJwtToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get branches: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	errorModel := response.JSON400
	if err = CheckResponseByModel(statusCode, body, errorModel); err != nil {
		return nil, fmt.Errorf("ai adapter get scans: %w", err)
	}

	models := *response.JSON200
	scans := make([]scan.Scan, len(models))
	for i, model := range models {
		scans[i].Id = *model.Id
	}

	return scans, nil
}

func (a *Adapter) GetScanQueue(ctx context.Context) ([]uuid.UUID, error) {
	response, err := a.aiClient.GetApiScansWithResponse(ctx, a.aiClient.AddJwtToHeader)
	if err != nil {
		return nil, err
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	if err = CheckResponseByModel(statusCode, body, nil); err != nil {
		return nil, fmt.Errorf("ai adapter get scans: %w", err)
	}

	models := *response.JSON200
	sort.Slice(models, func(i, j int) bool {
		first := *models[i].QueuingDateTime
		second := *models[j].QueuingDateTime
		return first.Before(second)
	})

	result := make([]uuid.UUID, len(models))
	for i, model := range models {
		result[i] = *model.ScanResultId
	}

	return result, nil
}

func (a *Adapter) StartScanBranch(ctx context.Context, branchId uuid.UUID) (uuid.UUID, error) {
	scanType := Incremental
	params := StartScanModel{
		ScanType: &scanType,
	}

	response, err := a.aiClient.PostApiScansBranchesBranchIdStartWithResponse(ctx, branchId, params, a.aiClient.AddJwtToHeader)
	if err != nil {
		return uuid.UUID{}, err
	}

	statusCode := response.StatusCode()
	responseBody := string(response.Body)
	errorModel := response.JSON400

	if err := CheckResponseByModel(statusCode, responseBody, errorModel); err != nil {
		return uuid.UUID{}, fmt.Errorf("ai adapter start scan: %w", err)
	}

	scanResultId, err := uuid.Parse(responseBody)
	if err != nil {
		return uuid.UUID{}, err
	}

	return scanResultId, nil
}

func (a *Adapter) StartScanProject(ctx context.Context, projectId uuid.UUID) (uuid.UUID, error) {
	scanType := Incremental
	params := StartScanModel{
		ScanType: &scanType,
	}

	response, err := a.aiClient.PostApiScansProjectIdStartWithResponse(ctx, projectId, params, a.aiClient.AddJwtToHeader)
	if err != nil {
		return uuid.UUID{}, err
	}

	statusCode := response.StatusCode()
	responseBody := string(response.Body)
	errorModel := response.JSON400

	if err := CheckResponseByModel(statusCode, responseBody, errorModel); err != nil {
		return uuid.UUID{}, fmt.Errorf("ai adapter start scan: %w", err)
	}

	scanResultId, err := uuid.Parse(responseBody)
	if err != nil {
		return uuid.UUID{}, err
	}

	return scanResultId, nil
}

func (a *Adapter) StopScan(ctx context.Context, scanResultId uuid.UUID) error {
	response, err := a.aiClient.PostApiScansScanResultIdStopWithResponse(ctx, scanResultId, a.aiClient.AddJwtToHeader)
	if err != nil {
		return err
	}

	statusCode := response.StatusCode()
	responseBody := string(response.Body)
	errorModel := response.JSON400
	if err = CheckResponseByModel(statusCode, responseBody, errorModel); err != nil {
		return fmt.Errorf("ai update sources post sources: %w", err)
	}

	return nil
}

func (a *Adapter) UpdateSources(ctx context.Context, projectId, branchId uuid.UUID, scanTargetPath string) error {
	archivePath, err := prepareArchive(scanTargetPath)
	if err != nil {
		return err
	}

	defer os.Remove(archivePath)

	body, contentType, err := prepareMultipartBody(archivePath)
	if err != nil {
		return err
	}

	archived := true
	params := PostApiStoreProjectIdBranchesBranchIdSourcesParams{Archived: &archived}

	response, err := a.aiClient.PostApiStoreProjectIdBranchesBranchIdSourcesWithBodyWithResponse(ctx, projectId, branchId, &params, contentType, body, a.aiClient.AddJwtToHeader)
	if err != nil {
		return fmt.Errorf("ai update sources: %w", err)
	}

	statusCode := response.StatusCode()
	responseBody := string(response.Body)
	errorModel := response.JSON400
	if err = CheckResponseByModel(statusCode, responseBody, errorModel); err != nil {
		return fmt.Errorf("ai update sources post sources: %w", err)
	}

	return nil
}
