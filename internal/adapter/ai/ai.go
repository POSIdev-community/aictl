package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/POSIdev-community/aictl/internal/core/domain/version"
	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/branch"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/scan"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanstage"
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
	. "github.com/POSIdev-community/aictl/pkg/clientai"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/POSIdev-community/aictl/pkg/logger"
)

type Adapter struct {
	aiClient *AiClient
	cfg      *config.Config
}

func NewAdapter(cfg *config.Config) (*Adapter, error) {
	aiClient, err := NewAiClient()
	if err != nil {
		return nil, fmt.Errorf("new ai client: %w", err)
	}

	return &Adapter{aiClient, cfg}, nil
}

func (a *Adapter) Initialize(ctx context.Context) error {
	err := a.aiClient.Initialize(ctx, a.cfg)
	if err != nil {
		return fmt.Errorf("initialize ai client: %w", err)
	}

	err = a.CheckMinVersion(ctx)
	if err != nil {
		return fmt.Errorf("check min version: %w", err)
	}

	err = a.CheckLicense(ctx)
	if err != nil {
		return fmt.Errorf("check license: %w", err)
	}

	return nil
}

func (a *Adapter) InitializeWithRetry(ctx context.Context) error {
	err := a.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("initialize ai adapter: %w", err)
	}

	a.aiClient.AddJwtRetry()

	return nil
}

func (a *Adapter) CheckMinVersion(ctx context.Context) error {
	v, err := a.GetVersion(ctx)
	if err != nil {
		return fmt.Errorf("get version: %w", err)
	}

	minVersion, _ := version.NewVersion("5.0.0")
	if v.Less(&minVersion) {
		return fmt.Errorf("version less than 5.0.0")
	}

	return nil
}

func (a *Adapter) GetDefaultSettings(ctx context.Context) (settings.ScanSettings, error) {
	res, err := a.aiClient.GetApiProjectsDefaultSettingsWithResponse(ctx, a.aiClient.AddJWTToHeader)
	if err != nil {
		return settings.ScanSettings{}, fmt.Errorf("get projects default settings: %w", err)
	}

	statusCode := res.StatusCode()
	responseBody := string(res.Body)
	if err = CheckResponseByModel(statusCode, responseBody, nil); err != nil {
		return settings.ScanSettings{}, fmt.Errorf("get projects default settings: %w", err)
	}

	defaultSettings := *res.JSON200

	return settings.ScanSettings{
		ProjectName: getOrDefault(defaultSettings.Name, ""),
		Languages: func() []string {
			if defaultSettings.Languages == nil {
				return nil
			}

			res := make([]string, len(*defaultSettings.Languages))
			for i := range *defaultSettings.Languages {
				res[i] = string((*defaultSettings.Languages)[i])
			}

			return res
		}(),
		WhiteBoxSettings: settings.WhiteBoxSettings{
			StaticCodeAnalysisEnabled:            getOrDefault(defaultSettings.WhiteBox.StaticCodeAnalysisEnabled, false),
			PatternMatchingEnabled:               getOrDefault(defaultSettings.WhiteBox.PatternMatchingEnabled, false),
			SearchForVulnerableComponentsEnabled: getOrDefault(defaultSettings.WhiteBox.SearchForVulnerableComponentsEnabled, false),
			SearchForConfigurationFlawsEnabled:   getOrDefault(defaultSettings.WhiteBox.SearchForConfigurationFlawsEnabled, false),
			SearchWithScaEnabled:                 getOrDefault(defaultSettings.WhiteBox.SearchWithScaEnabled, false),
		},
		DotNetSettings: settings.DotNetSettings{
			ProjectType:                           string(getOrDefault(defaultSettings.DotNetSettings.ProjectType, "")),
			SolutionFile:                          getOrDefault(defaultSettings.DotNetSettings.SolutionFile, ""),
			WebSiteFolder:                         getOrDefault(defaultSettings.DotNetSettings.WebSiteFolder, ""),
			LaunchParameters:                      getOrDefault(defaultSettings.DotNetSettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: getOrDefault(defaultSettings.DotNetSettings.UseAvailablePublicAndProtectedMethods, false),
			DownloadDependencies:                  getOrDefault(defaultSettings.DotNetSettings.DownloadDependencies, false),
		},
		GoSettings: settings.GoSettings{
			LaunchParameters:                      getOrDefault(defaultSettings.GoSettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: getOrDefault(defaultSettings.GoSettings.UseAvailablePublicAndProtectedMethods, false),
		},
		JavaScriptSettings: settings.JavaScriptSettings{
			LaunchParameters:                      getOrDefault(defaultSettings.JavaScriptSettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: getOrDefault(defaultSettings.JavaScriptSettings.UseAvailablePublicAndProtectedMethods, false),
			DownloadDependencies:                  getOrDefault(defaultSettings.JavaScriptSettings.DownloadDependencies, false),
			UseTaintAnalysis:                      getOrDefault(defaultSettings.JavaScriptSettings.UseTaintAnalysis, false),
			UseJsaAnalysis:                        getOrDefault(defaultSettings.JavaScriptSettings.UseJsaAnalysis, false),
		},
		JavaSettings: settings.JavaSettings{
			Parameters:                            getOrDefault(defaultSettings.JavaSettings.Parameters, ""),
			UnpackUserPackages:                    getOrDefault(defaultSettings.JavaSettings.UnpackUserPackages, false),
			UserPackagePrefixes:                   getOrDefault(defaultSettings.JavaSettings.UserPackagePrefixes, ""),
			Version:                               string(getOrDefault(defaultSettings.JavaSettings.Version, "")),
			LaunchParameters:                      getOrDefault(defaultSettings.JavaSettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: getOrDefault(defaultSettings.JavaSettings.UseAvailablePublicAndProtectedMethods, false),
			DownloadDependencies:                  getOrDefault(defaultSettings.JavaSettings.DownloadDependencies, false),
			DependenciesPath:                      getOrDefault(defaultSettings.JavaSettings.DependenciesPath, ""),
		},
		PhpSettings: settings.PhpSettings{
			LaunchParameters:                      getOrDefault(defaultSettings.PhpSettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: getOrDefault(defaultSettings.PhpSettings.UseAvailablePublicAndProtectedMethods, false),
			DownloadDependencies:                  getOrDefault(defaultSettings.PhpSettings.DownloadDependencies, false),
		},
		PmTaintSettings: settings.PmTaintSettings{
			LaunchParameters:                      getOrDefault(defaultSettings.PmTaintSettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: getOrDefault(defaultSettings.PmTaintSettings.UseAvailablePublicAndProtectedMethods, false),
		},
		PythonSettings: settings.PythonSettings{
			LaunchParameters:                      getOrDefault(defaultSettings.PythonSettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: getOrDefault(defaultSettings.PythonSettings.UseAvailablePublicAndProtectedMethods, false),
			DownloadDependencies:                  getOrDefault(defaultSettings.PythonSettings.DownloadDependencies, false),
			DependenciesPath:                      getOrDefault(defaultSettings.PythonSettings.DependenciesPath, ""),
		},
		RubySettings: settings.RubySettings{
			LaunchParameters:                      getOrDefault(defaultSettings.RubySettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: getOrDefault(defaultSettings.RubySettings.UseAvailablePublicAndProtectedMethods, false),
		},
		PygrepSettings: settings.PygrepSettings{
			RulesDirPath:     getOrDefault(defaultSettings.PygrepSettings.RulesDirPath, ""),
			LaunchParameters: getOrDefault(defaultSettings.PygrepSettings.LaunchParameters, ""),
		},
		ScaSettings: settings.ScaSettings{
			LaunchParameters:       getOrDefault(defaultSettings.ScaSettings.LaunchParameters, ""),
			BuildDependenciesGraph: getOrDefault(defaultSettings.ScaSettings.BuildDependenciesGraph, false),
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

	res, err := a.aiClient.PutApiProjectsProjectIdSettingsWithResponse(ctx, projectId, projectSettings, a.aiClient.AddJWTToHeader)
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

func (a *Adapter) CreateBranch(ctx context.Context, projectId uuid.UUID, branchName, scanTargetPath string) (*uuid.UUID, error) {
	useStubSources := scanTargetPath == ""
	if useStubSources {
		var err error
		scanTargetPath, err = createStubScanTarget()
		if err != nil {
			return nil, err
		}
	}

	archivePath, err := prepareArchive(scanTargetPath)
	if archivePath != scanTargetPath {
		defer func() {
			_ = os.Remove(archivePath)
		}()
	}
	if err != nil {
		return nil, err
	}

	body, contentType, err := prepareMultipartBody(ctx, archivePath, !useStubSources, MultipartField{Key: "Name", Value: branchName})
	if err != nil {
		return nil, err
	}

	response, err := a.aiClient.PostApiStoreProjectProjectIdBranchesArchiveWithBodyWithResponse(ctx, projectId, contentType, body, a.aiClient.AddJWTToHeader)
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

	createProjectResponse, err := a.aiClient.PostApiProjectsBaseWithResponse(ctx, projectBaseModel, a.aiClient.AddJWTToHeader)
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
	response, err := a.aiClient.DeleteApiProjectsProjectId(ctx, projectId, a.aiClient.AddJWTToHeader)
	if err != nil {
		return fmt.Errorf("ai adapter delete project: %w", err)
	}

	if err = CheckResponse(response, "project"); err != nil {
		return fmt.Errorf("ai adapter delete project: %w", err)
	}

	return nil
}

func (a *Adapter) ExistsProject(ctx context.Context, projectName string) (bool, error) {
	response, err := a.aiClient.GetApiProjectsNameExistsWithResponse(ctx, projectName, a.aiClient.AddJWTToHeader)
	if err != nil {
		return false, fmt.Errorf("ai adapter get project name exists: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	errorModel := response.JSON400
	if err = CheckResponseByModel(statusCode, body, errorModel); err != nil {
		return false, err
	}

	boolValueTrue, err := strconv.ParseBool(body)
	if err != nil {
		return false, err
	}

	return boolValueTrue, nil
}

func (a *Adapter) GetProjectId(ctx context.Context, projectName string) (*uuid.UUID, error) {
	response, err := a.aiClient.GetApiProjectsNameNameWithResponse(ctx, projectName, a.aiClient.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get project name exists: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	errorModel := response.JSON400

	if statusCode == http.StatusBadRequest && errorModel != nil && *errorModel.ErrorCode == ApiErrorTypePROJECTNOTFOUND {
		return nil, nil
	}

	if err = CheckResponseByModel(statusCode, body, errorModel); err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("ai adapter get project name exists: %w", err)
	}

	proj := *response.JSON200

	return proj.Id, nil
}

func (a *Adapter) GetProjects(ctx context.Context) ([]project.Project, error) {
	log := logger.FromContext(ctx)

	log.StdErrf("Send get projects request")

	response, err := a.aiClient.GetApiProjectsWithResponse(ctx, nil, a.aiClient.AddJWTToHeader)
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
	response, err := a.aiClient.GetApiProjectsProjectIdWithResponse(ctx, projectId, a.aiClient.AddJWTToHeader)
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

func (a *Adapter) GetDefaultTemplateId(ctx context.Context, reportType report.ReportType) (uuid.UUID, error) {
	localeId := "ru-Ru"
	params := GetApiReportsTemplatesTypeParams{
		LocaleId: &localeId,
	}

	var aiReportType ReportType
	switch reportType {
	case report.AutoCheck:
		aiReportType = ReportTypeAutoCheck
	case report.Custom:
		aiReportType = ReportTypeCustom
	case report.Gitlab:
		aiReportType = ReportTypeGitlab
	case report.Json:
		aiReportType = ReportTypeJson
	case report.Markdown:
		aiReportType = ReportTypeMd
	case report.Nist:
		aiReportType = ReportTypeNist
	case report.Oud4:
		aiReportType = ReportTypeOud4
	case report.Owasp:
		aiReportType = ReportTypeOwasp
	case report.Owaspm:
		aiReportType = ReportTypeOwaspm
	case report.Pcidss:
		aiReportType = ReportTypePcidss
	case report.PlainReport:
		aiReportType = ReportTypePlainReport
	case report.Sans:
		aiReportType = ReportTypeSans
	case report.Sarif:
		aiReportType = ReportTypeSarif
	case report.Xml:
		aiReportType = ReportTypeXml
	default:
		return uuid.UUID{}, fmt.Errorf("invalid reportType: %s", reportType)
	}

	response, err := a.aiClient.GetApiReportsTemplatesType(ctx, aiReportType, &params, a.aiClient.AddJWTToHeader)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("ai adapter get default template id: %w", err)
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

func (a *Adapter) GetCustomTemplateId(ctx context.Context, reportName string) (uuid.UUID, error) {
	localeId := "ru-RU"
	params := GetApiReportsUserTemplatesNameParams{
		LocaleId: &localeId,
	}

	response, err := a.aiClient.GetApiReportsUserTemplatesName(ctx, reportName, &params, a.aiClient.AddJWTToHeader)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("ai adapter get custom template id: %w", err)
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

	var model ReportTemplateSimpleModel
	if err := json.Unmarshal(bodyBytes, &model); err != nil {
		return uuid.UUID{}, err
	}

	return *model.Id, nil
}

func (a *Adapter) GetReport(ctx context.Context, projectId, scanResultId, templateId uuid.UUID, includeComments, includeDFD, includeGlossary bool, l10n string) (io.ReadCloser, error) {
	useFilters := false
	sessionId := uuid.New()

	model := ReportGenerateModel{
		LocaleId: &l10n,
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

	response, err := a.aiClient.PostApiReportsGenerate(ctx, model, a.aiClient.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter generate report: %w", err)
	}

	if err = CheckResponse(response, "report"); err != nil {
		return nil, fmt.Errorf("ai adapter generate report: %w", err)
	}

	return response.Body, nil
}

func (a *Adapter) GetSbom(ctx context.Context, projectId, scanResultId uuid.UUID) (io.ReadCloser, error) {
	response, err := a.aiClient.GetApiStoreProjectIdSbomsScanResultId(ctx, projectId, scanResultId, a.aiClient.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get sbom: %w", err)
	}

	if err = CheckResponse(response, "sbom"); err != nil {
		return nil, fmt.Errorf("ai adapter get sbom: %w", err)
	}

	return response.Body, nil
}

func (a *Adapter) GetScanLogs(ctx context.Context, projectId, scanResultId uuid.UUID) (io.ReadCloser, error) {
	response, err := a.aiClient.GetApiStoreProjectIdLogsScanResultId(ctx, projectId, scanResultId, a.aiClient.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get scan logs: %w", err)
	}

	if err = CheckResponse(response, "logs"); err != nil {
		return nil, fmt.Errorf("ai adapter get scan logs: %w", err)
	}

	return response.Body, nil
}

func (a *Adapter) GetBranches(ctx context.Context, projectId uuid.UUID) ([]branch.Branch, error) {
	getBranchesResponse, err := a.aiClient.GetApiProjectsProjectIdBranchesWithResponse(ctx, projectId, a.aiClient.AddJWTToHeader)
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
		if model.Description == nil {
			description := ""
			model.Description = &description
		}
		branches[i] = branch.NewBranch(*model.Id, *model.Name, *model.Description, *model.IsWorking)
	}

	return branches, nil
}

func (a *Adapter) GetLastScan(ctx context.Context, branchId uuid.UUID) (*scan.Scan, error) {
	response, err := a.aiClient.GetApiBranchesBranchIdScanResultsLastWithResponse(ctx, branchId, a.aiClient.AddJWTToHeader)
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
	response, err := a.aiClient.GetApiProjectsProjectIdScanResultsScanResultIdWithResponse(ctx, projectId, scanId, a.aiClient.AddJWTToHeader)
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

func (a *Adapter) GetScanAiproj(ctx context.Context, projectId, scanSettingsId uuid.UUID) (io.ReadCloser, error) {
	response, err := a.aiClient.GetApiProjectsProjectIdScanSettingsScanSettingsIdAiproj(ctx, projectId, scanSettingsId, a.aiClient.AddJWTToHeader)
	if err != nil {
		return nil, errs.NewNotFoundError("get scan aiproj")
	}

	if err = CheckResponse(response, "aiproj"); err != nil {
		return nil, fmt.Errorf("ai adapter get aiproj: %w", err)
	}

	return response.Body, nil
}

func (a *Adapter) GetScanStage(ctx context.Context, projectId, scanId uuid.UUID) (scanstage.ScanStage, error) {
	response, err := a.aiClient.GetApiProjectsProjectIdScanResultsScanResultIdProgressWithResponse(ctx, projectId, scanId, a.aiClient.AddJWTToHeader)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return scanstage.ScanStage{}, err
		}

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

func (a *Adapter) GetScanQueue(ctx context.Context) ([]uuid.UUID, error) {
	response, err := a.aiClient.GetApiScansWithResponse(ctx, a.aiClient.AddJWTToHeader)
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

func (a *Adapter) StartScanBranch(ctx context.Context, branchId uuid.UUID, scanLabel string) (uuid.UUID, error) {
	scanType := Incremental
	params := StartScanModel{
		ScanType:  &scanType,
		ScanLabel: &scanLabel,
	}

	response, err := a.aiClient.PostApiScansBranchesBranchIdStartWithResponse(ctx, branchId, params, a.aiClient.AddJWTToHeader)
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

func (a *Adapter) StartScanProject(ctx context.Context, projectId uuid.UUID, scanLabel string) (uuid.UUID, error) {
	scanType := Incremental
	params := StartScanModel{
		ScanType:  &scanType,
		ScanLabel: &scanLabel,
	}

	response, err := a.aiClient.PostApiScansProjectIdStartWithResponse(ctx, projectId, params, a.aiClient.AddJWTToHeader)
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
	response, err := a.aiClient.PostApiScansScanResultIdStopWithResponse(ctx, scanResultId, a.aiClient.AddJWTToHeader)
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
	log := logger.FromContext(ctx)

	archivePath, err := prepareArchive(scanTargetPath)
	if archivePath != scanTargetPath {
		defer func() {
			_ = os.Remove(archivePath)
		}()
	}
	if err != nil {
		return err
	}

	if fi, err := os.Stat(archivePath); err == nil {
		log.StdErrf("archive prepared, size: %.1f MB", float64(fi.Size())/(1024*1024))
	}

	body, contentType, err := prepareMultipartBody(ctx, archivePath, true)
	if err != nil {
		return err
	}
	defer func() {
		_ = body.Close()
	}()

	archived := true
	params := PostApiStoreProjectIdBranchesBranchIdSourcesParams{Archived: &archived}

	response, err := a.aiClient.PostApiStoreProjectIdBranchesBranchIdSourcesWithBodyWithResponse(ctx, projectId, branchId, &params, contentType, body, a.aiClient.AddJWTToHeader)
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

func (a *Adapter) GetVersion(ctx context.Context) (version.Version, error) {
	response, err := a.aiClient.GetApiVersionsPackageCurrentWithResponse(ctx, a.aiClient.AddJWTToHeader)
	if err != nil {
		return version.Version{}, fmt.Errorf("ai get version: %w", err)
	}

	statusCode := response.StatusCode()
	responseBody := string(response.Body)
	if err = CheckResponseByModel(statusCode, responseBody, nil); err != nil {
		return version.Version{}, fmt.Errorf("ai get version: %w", err)
	}

	v, err := version.NewVersion(responseBody)

	return v, nil
}

func (a *Adapter) GetHealthcheck(ctx context.Context) (bool, error) {
	response, err := a.aiClient.GetHealthSummaryWithResponse(ctx, a.aiClient.AddJWTToHeader)
	if err != nil {
		return false, fmt.Errorf("ai get version: %w", err)
	}

	statusCode := response.StatusCode()
	responseBody := string(response.Body)
	if err = CheckResponseByModel(statusCode, responseBody, nil); err != nil {
		return false, fmt.Errorf("ai get version: %w", err)
	}

	health := true
	for _, service := range *response.JSON200.Services {
		health = health && *service.Status == "Healthy"
	}

	return health, nil
}

func (a *Adapter) CheckLicense(ctx context.Context) error {
	response, err := a.aiClient.GetApiLicenseWithResponse(ctx, a.aiClient.AddJWTToHeader)
	if err != nil {
		return fmt.Errorf("ai check license: %w", err)
	}

	statusCode := response.StatusCode()
	responseBody := string(response.Body)
	if err = CheckResponseByModel(statusCode, responseBody, nil); err != nil {
		return fmt.Errorf("ai check license: %w", err)
	}

	if !*response.JSON200.IsValid {
		return fmt.Errorf("license is invalid")
	}

	return nil
}
