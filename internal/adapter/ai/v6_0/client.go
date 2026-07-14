package v6_0

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	"github.com/POSIdev-community/aictl/internal/core/apperror"
	"github.com/POSIdev-community/aictl/internal/core/domain/branch"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/queue"
	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/POSIdev-community/aictl/internal/core/domain/scan"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanstage"
	"github.com/POSIdev-community/aictl/internal/core/domain/scantype"
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/internal/core/domain/statistic"
	"github.com/POSIdev-community/aictl/internal/core/domain/version"
	"github.com/POSIdev-community/aictl/pkg/clientai/v6_0"
	"github.com/POSIdev-community/aictl/pkg/gitignore"
	"github.com/POSIdev-community/aictl/pkg/logger"
	"github.com/google/uuid"
)

type ClientAI60 struct {
	*v6_0.ClientWithResponses
	jwtClient *v6_0.ClientWithResponses

	*common.BaseClient
}

func NewAiClient(base *common.BaseClient) *ClientAI60 {
	return &ClientAI60{
		BaseClient: base,
	}
}

func (a *ClientAI60) Initialize(ctx context.Context, cfg *config.Config) error {
	client, err := v6_0.NewClientWithResponses(cfg.UriString(), v6_0.WithHTTPClient(a.HttpClient))
	if err != nil {
		return fmt.Errorf("new client: %w", err)
	}
	a.ClientWithResponses = client

	a.jwtClient, err = v6_0.NewClientWithResponses(cfg.UriString(), v6_0.WithHTTPClient(a.JwtHttpClient))
	if err != nil {
		return fmt.Errorf("new jwt client: %w", err)
	}

	transport := &http.Transport{}
	if cfg.TLSSkip() {
		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{}
		}

		transport.TLSClientConfig.InsecureSkipVerify = true
	}

	a.HttpClient.Transport = transport.Clone()
	a.JwtHttpClient.Transport = transport.Clone()

	if err := a.getJWT(ctx, cfg); err != nil {
		return fmt.Errorf("update jwt: %w", err)
	}

	a.Initialized = true

	return nil
}

func (a *ClientAI60) AddJwtRetry() {
	a.HttpClient.Transport = common.NewRetryRoundTripper(a.HttpClient.Transport, http.StatusUnauthorized, a.refreshJWT)

	a.WithRetry = true
}

func (a *ClientAI60) getJWT(ctx context.Context, cfg *config.Config) error {
	if a.Initialized {
		return nil
	}

	response, err := a.jwtClient.GetApiAuthSigninWithResponse(ctx, func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Access-Token", cfg.Token())

		return nil
	})
	if err != nil {
		return fmt.Errorf("get api auth signin: %w", err)
	}

	if err = CheckResponseByModel(response.StatusCode(), string(response.Body), response.JSON400); err != nil {
		return err
	}

	a.AccessToken = *response.JSON200.AccessToken
	a.RefreshToken = *response.JSON200.RefreshToken

	return nil
}

func (a *ClientAI60) refreshJWT(ctx context.Context, req *http.Request) error {
	err := a.DoJWTRefresh(func() error {
		log := logger.FromContext(ctx)

		response, err := a.jwtClient.GetApiAuthRefreshTokenWithResponse(ctx, func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+a.RefreshToken)

			return nil
		})
		if err != nil {
			return fmt.Errorf("get api auth refresh token: %w", err)
		}

		if err = CheckResponse(response.HTTPResponse, "jwt refresh"); err != nil {
			return err
		}

		if response.JSON200 == nil || response.JSON200.AccessToken == nil {
			log.StdErrf("Got empty access token")

			return fmt.Errorf("no access token")
		}

		a.AccessToken = *response.JSON200.AccessToken

		return nil
	})
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+a.AccessToken)

	return nil
}

func (a *ClientAI60) GetDefaultSettings(ctx context.Context) (settings.ScanSettings, error) {
	res, err := a.GetApiProjectsDefaultSettingsWithResponse(ctx, a.AddJWTToHeader)
	if err != nil {
		return settings.ScanSettings{}, fmt.Errorf("get projects default settings request: %w", err)
	}

	statusCode := res.StatusCode()
	responseBody := string(res.Body)
	if err = CheckResponseByModel(statusCode, responseBody, nil); err != nil {
		return settings.ScanSettings{}, fmt.Errorf("get projects default settings: %w", err)
	}

	defaultSettings := *res.JSON200

	return settings.ScanSettings{
		ProjectName: common.GetOrDefault(defaultSettings.Name, ""),
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
			StaticCodeAnalysisEnabled:            common.GetOrDefault(defaultSettings.WhiteBox.StaticCodeAnalysisEnabled, false),
			PatternMatchingEnabled:               common.GetOrDefault(defaultSettings.WhiteBox.PatternMatchingEnabled, false),
			SearchForVulnerableComponentsEnabled: common.GetOrDefault(defaultSettings.WhiteBox.SearchForVulnerableComponentsEnabled, false),
			SearchForConfigurationFlawsEnabled:   common.GetOrDefault(defaultSettings.WhiteBox.SearchForConfigurationFlawsEnabled, false),
			SearchWithScaEnabled:                 common.GetOrDefault(defaultSettings.WhiteBox.SearchWithScaEnabled, false),
		},
		DotNetSettings: settings.DotNetSettings{
			ProjectType:                           string(common.GetOrDefault(defaultSettings.DotNetSettings.ProjectType, "")),
			SolutionFile:                          common.GetOrDefault(defaultSettings.DotNetSettings.SolutionFile, ""),
			LaunchParameters:                      common.GetOrDefault(defaultSettings.DotNetSettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: common.GetOrDefault(defaultSettings.DotNetSettings.UseAvailablePublicAndProtectedMethods, false),
			DownloadDependencies:                  common.GetOrDefault(defaultSettings.DotNetSettings.DownloadDependencies, false),
		},
		GoSettings: settings.GoSettings{
			LaunchParameters:                      common.GetOrDefault(defaultSettings.GoSettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: common.GetOrDefault(defaultSettings.GoSettings.UseAvailablePublicAndProtectedMethods, false),
		},
		JavaScriptSettings: settings.JavaScriptSettings{
			LaunchParameters:                      common.GetOrDefault(defaultSettings.JavaScriptSettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: common.GetOrDefault(defaultSettings.JavaScriptSettings.UseAvailablePublicAndProtectedMethods, false),
			DownloadDependencies:                  common.GetOrDefault(defaultSettings.JavaScriptSettings.DownloadDependencies, false),
			UseTaintAnalysis:                      common.GetOrDefault(defaultSettings.JavaScriptSettings.UseTaintAnalysis, false),
			UseJsaAnalysis:                        common.GetOrDefault(defaultSettings.JavaScriptSettings.UseJsaAnalysis, false),
		},
		JavaSettings: settings.JavaSettings{
			Parameters:                            common.GetOrDefault(defaultSettings.JavaSettings.Parameters, ""),
			UnpackUserPackages:                    common.GetOrDefault(defaultSettings.JavaSettings.UnpackUserPackages, false),
			UserPackagePrefixes:                   common.GetOrDefault(defaultSettings.JavaSettings.UserPackagePrefixes, ""),
			Version:                               string(common.GetOrDefault(defaultSettings.JavaSettings.Version, "")),
			LaunchParameters:                      common.GetOrDefault(defaultSettings.JavaSettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: common.GetOrDefault(defaultSettings.JavaSettings.UseAvailablePublicAndProtectedMethods, false),
			DownloadDependencies:                  common.GetOrDefault(defaultSettings.JavaSettings.DownloadDependencies, false),
			DependenciesPath:                      common.GetOrDefault(defaultSettings.JavaSettings.DependenciesPath, ""),
		},
		PhpSettings: settings.PhpSettings{
			LaunchParameters:                      common.GetOrDefault(defaultSettings.PhpSettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: common.GetOrDefault(defaultSettings.PhpSettings.UseAvailablePublicAndProtectedMethods, false),
			DownloadDependencies:                  common.GetOrDefault(defaultSettings.PhpSettings.DownloadDependencies, false),
		},
		PmTaintSettings: settings.PmTaintSettings{
			LaunchParameters:                      common.GetOrDefault(defaultSettings.PmTaintSettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: common.GetOrDefault(defaultSettings.PmTaintSettings.UseAvailablePublicAndProtectedMethods, false),
		},
		PythonSettings: settings.PythonSettings{
			LaunchParameters:                      common.GetOrDefault(defaultSettings.PythonSettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: common.GetOrDefault(defaultSettings.PythonSettings.UseAvailablePublicAndProtectedMethods, false),
			DownloadDependencies:                  common.GetOrDefault(defaultSettings.PythonSettings.DownloadDependencies, false),
			DependenciesPath:                      common.GetOrDefault(defaultSettings.PythonSettings.DependenciesPath, ""),
		},
		RubySettings: settings.RubySettings{
			LaunchParameters:                      common.GetOrDefault(defaultSettings.RubySettings.LaunchParameters, ""),
			UseAvailablePublicAndProtectedMethods: common.GetOrDefault(defaultSettings.RubySettings.UseAvailablePublicAndProtectedMethods, false),
		},
		ScaSettings: settings.ScaSettings{
			LaunchParameters:       common.GetOrDefault(defaultSettings.ScaSettings.LaunchParameters, ""),
			BuildDependenciesGraph: common.GetOrDefault(defaultSettings.ScaSettings.BuildDependenciesGraph, false),
		},
	}, err
}

func (a *ClientAI60) SetProjectSettings(ctx context.Context, projectId uuid.UUID, settings *settings.ScanSettings) error {
	if settings == nil {
		return nil
	}

	priority := priorityFromSettings(settings)
	projectSettings := v6_0.PutApiProjectsProjectIdSettingsJSONRequestBody{
		ProjectName: &settings.ProjectName,
		Priority:    &priority,
		Languages: func() *[]v6_0.LegacyProgrammingLanguageGroup {
			if settings.Languages == nil {
				return nil
			}
			res := make([]v6_0.LegacyProgrammingLanguageGroup, len(settings.Languages))
			for i := range settings.Languages {
				res[i] = v6_0.LegacyProgrammingLanguageGroup(settings.Languages[i])
			}
			return &res
		}(),
		WhiteBoxSettings: &v6_0.WhiteBoxSettingsModel{
			StaticCodeAnalysisEnabled:            &settings.WhiteBoxSettings.StaticCodeAnalysisEnabled,
			PatternMatchingEnabled:               &settings.WhiteBoxSettings.PatternMatchingEnabled,
			SearchForVulnerableComponentsEnabled: &settings.WhiteBoxSettings.SearchForVulnerableComponentsEnabled,
			SearchForConfigurationFlawsEnabled:   &settings.WhiteBoxSettings.SearchForConfigurationFlawsEnabled,
			SearchWithScaEnabled:                 &settings.WhiteBoxSettings.SearchWithScaEnabled,
			SecretDetectionEnabled:               &settings.WhiteBoxSettings.SecretDetectionEnabled,
			SearchForMaliciousCodeEnabled:        &settings.WhiteBoxSettings.SearchForMaliciousCodeEnabled,
		},
		DotNetSettings: &v6_0.DotNetSettingsModel{
			ProjectType:                           common.Reference(v6_0.DotNetProjectType(settings.DotNetSettings.ProjectType)),
			SolutionFile:                          &settings.DotNetSettings.SolutionFile,
			LaunchParameters:                      &settings.DotNetSettings.LaunchParameters,
			UseAvailablePublicAndProtectedMethods: &settings.DotNetSettings.UseAvailablePublicAndProtectedMethods,
			DownloadDependencies:                  &settings.DotNetSettings.DownloadDependencies,
		},
		GoSettings: &v6_0.GoSettingsModel{
			LaunchParameters:                      &settings.GoSettings.LaunchParameters,
			UseAvailablePublicAndProtectedMethods: &settings.GoSettings.UseAvailablePublicAndProtectedMethods,
		},
		JavaScriptSettings: &v6_0.JavaScriptSettingsModel{
			LaunchParameters:                      &settings.JavaScriptSettings.LaunchParameters,
			UseAvailablePublicAndProtectedMethods: &settings.JavaScriptSettings.UseAvailablePublicAndProtectedMethods,
			DownloadDependencies:                  &settings.JavaScriptSettings.DownloadDependencies,
			UseTaintAnalysis:                      &settings.JavaScriptSettings.UseTaintAnalysis,
			UseJsaAnalysis:                        &settings.JavaScriptSettings.UseJsaAnalysis,
		},
		JavaSettings: &v6_0.JavaSettingsModel{
			Parameters:                            &settings.JavaSettings.Parameters,
			UnpackUserPackages:                    &settings.JavaSettings.UnpackUserPackages,
			UserPackagePrefixes:                   &settings.JavaSettings.UserPackagePrefixes,
			Version:                               common.Reference(v6_0.JavaVersions(settings.JavaSettings.Version)),
			LaunchParameters:                      &settings.JavaSettings.LaunchParameters,
			UseAvailablePublicAndProtectedMethods: &settings.JavaSettings.UseAvailablePublicAndProtectedMethods,
			DownloadDependencies:                  &settings.JavaSettings.DownloadDependencies,
			DependenciesPath:                      &settings.JavaSettings.DependenciesPath,
		},
		PhpSettings: &v6_0.PhpSettingsModel{
			LaunchParameters:                      &settings.PhpSettings.LaunchParameters,
			UseAvailablePublicAndProtectedMethods: &settings.PhpSettings.UseAvailablePublicAndProtectedMethods,
			DownloadDependencies:                  &settings.PhpSettings.DownloadDependencies,
		},
		PmTaintSettings: &v6_0.PmTaintBaseSettingsModel{
			LaunchParameters:                      &settings.PmTaintSettings.LaunchParameters,
			UseAvailablePublicAndProtectedMethods: &settings.PmTaintSettings.UseAvailablePublicAndProtectedMethods,
		},
		PythonSettings: &v6_0.PythonSettingsModel{
			LaunchParameters:                      &settings.PythonSettings.LaunchParameters,
			UseAvailablePublicAndProtectedMethods: &settings.PythonSettings.UseAvailablePublicAndProtectedMethods,
			DownloadDependencies:                  &settings.PythonSettings.DownloadDependencies,
			DependenciesPath:                      &settings.PythonSettings.DependenciesPath,
		},
		RubySettings: &v6_0.RubySettingsModel{
			LaunchParameters:                      &settings.RubySettings.LaunchParameters,
			UseAvailablePublicAndProtectedMethods: &settings.RubySettings.UseAvailablePublicAndProtectedMethods,
		},
		ScaSettings: &v6_0.ScaSettingsModel{
			LaunchParameters:       &settings.ScaSettings.LaunchParameters,
			BuildDependenciesGraph: &settings.ScaSettings.BuildDependenciesGraph,
		},
	}

	res, err := a.PutApiProjectsProjectIdSettingsWithResponse(ctx, projectId, projectSettings, a.AddJWTToHeader)
	if err != nil {
		return fmt.Errorf("put project settings request: %w", err)
	}

	statusCode := res.StatusCode()
	responseBody := string(res.Body)
	errorModel := res.JSON400
	if err = CheckResponseByModel(statusCode, responseBody, errorModel); err != nil {
		return fmt.Errorf("put project settings: %w", err)
	}

	if err := a.putPreferredAgentsSettings(ctx, projectId, settings); err != nil {
		return err
	}

	if settings.HasBlackBoxSettings() {
		if err := a.setBlackBoxSettings(ctx, projectId, settings); err != nil {
			return err
		}
	}

	return nil
}

func (a *ClientAI60) CreateBranch(ctx context.Context, projectId uuid.UUID, branchName, scanTargetPath string, exclusions gitignore.Exclusions) (*uuid.UUID, error) {
	useStubSources := scanTargetPath == ""
	if useStubSources {
		var err error
		scanTargetPath, err = common.CreateStubScanTarget()
		if err != nil {
			return nil, err
		}
	}

	archivePath, err := common.PrepareArchive(ctx, scanTargetPath, exclusions)
	if archivePath != scanTargetPath {
		defer func() {
			_ = os.Remove(archivePath)
		}()
	}
	if err != nil {
		return nil, err
	}

	body, contentType, err := common.PrepareMultipartBody(ctx, archivePath, !useStubSources, common.MultipartField{Key: "Name", Value: branchName})
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = body.Close()
	}()

	response, err := a.PostApiStoreProjectProjectIdBranchesArchiveWithBodyWithResponse(ctx, projectId, contentType, body, a.AddJWTToHeader)
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

func (a *ClientAI60) CreateProject(ctx context.Context, projectName string) (*uuid.UUID, error) {
	projectUrl := "http://localhost"

	patternMatchingEnabled := true
	searchForConfigurationFlawsEnabled := true
	searchForVulnerableComponentsEnabled := true
	searchWithScaEnabled := false
	staticCodeAnalysisEnabled := true
	preferredAgentsOnly := false
	preferredAgents := []uuid.UUID{}
	priority := v6_0.PriorityLow

	projectBaseModel := v6_0.PostApiProjectsBaseJSONRequestBody{
		Name:       &projectName,
		ProjectUrl: &projectUrl,
		WhiteBox: &v6_0.WhiteBoxSettingsModel{
			PatternMatchingEnabled:               &patternMatchingEnabled,
			SearchForConfigurationFlawsEnabled:   &searchForConfigurationFlawsEnabled,
			SearchForVulnerableComponentsEnabled: &searchForVulnerableComponentsEnabled,
			SearchWithScaEnabled:                 &searchWithScaEnabled,
			StaticCodeAnalysisEnabled:            &staticCodeAnalysisEnabled,
		},
		Id: &uuid.UUID{},
		Languages: &[]v6_0.LegacyProgrammingLanguageGroup{
			v6_0.LegacyProgrammingLanguageGroupGo,
		},
		PreferredAgentsSettings: &v6_0.PreferredAgentsSettings{
			PreferredAgents:     &preferredAgents,
			PreferredAgentsOnly: preferredAgentsOnly,
		},
		Priority: &priority,
	}

	createProjectResponse, err := a.PostApiProjectsBaseWithResponse(ctx, projectBaseModel, a.AddJWTToHeader)
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

func (a *ClientAI60) DeleteProject(ctx context.Context, projectId uuid.UUID) error {
	response, err := a.DeleteApiProjectsProjectId(ctx, projectId, a.AddJWTToHeader)
	if err != nil {
		return fmt.Errorf("ai adapter delete project request: %w", err)
	}

	if err = CheckResponse(response, "project"); err != nil {
		return fmt.Errorf("ai adapter delete project: %w", err)
	}

	return nil
}

func (a *ClientAI60) ExistsProject(ctx context.Context, projectName string) (bool, error) {
	response, err := a.GetApiProjectsNameExistsWithResponse(ctx, projectName, a.AddJWTToHeader)
	if err != nil {
		return false, fmt.Errorf("ai adapter get project name exists request: %w", err)
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

func (a *ClientAI60) GetProjectId(ctx context.Context, projectName string) (*uuid.UUID, error) {
	response, err := a.GetApiProjectsNameNameWithResponse(ctx, projectName, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get project name exists request: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	errorModel := response.JSON400

	if statusCode == http.StatusBadRequest && errorModel != nil && *errorModel.ErrorCode == v6_0.ApiErrorTypePROJECTNOTFOUND {
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

func (a *ClientAI60) GetProjects(ctx context.Context) ([]project.Project, error) {
	log := logger.FromContext(ctx)

	log.StdErrf("Send get projects request")

	response, err := a.GetApiProjectsWithResponse(ctx, nil, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get projects request request: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	if err = CheckResponseByModel(statusCode, body, nil); err != nil {
		return nil, fmt.Errorf("ai adapter get projects response: %w", err)
	}

	models := *response.JSON200
	projects := make([]project.Project, 0, len(models))

	for _, model := range models {
		if *model.ProjectType != v6_0.Permanent {
			continue
		}

		p := project.NewProject(*model.Id, *model.Name)
		projects = append(projects, p)
	}

	return projects, nil
}

func (a *ClientAI60) GetProject(ctx context.Context, projectId uuid.UUID) (*project.Project, error) {
	response, err := a.GetApiProjectsProjectIdWithResponse(ctx, projectId, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get projects request: %w", err)
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

func (a *ClientAI60) GetDefaultTemplateId(ctx context.Context, reportType report.ReportType) (uuid.UUID, error) {
	localeId := "ru-Ru"
	params := v6_0.GetApiReportsTemplatesTypeParams{
		LocaleId: &localeId,
	}

	var aiReportType v6_0.ReportType
	switch reportType {
	case report.AutoCheck:
		aiReportType = v6_0.ReportTypeAutoCheck
	case report.Custom:
		aiReportType = v6_0.ReportTypeCustom
	case report.Gitlab:
		aiReportType = v6_0.ReportTypeGitlab
	case report.Json:
		aiReportType = v6_0.ReportTypeJson
	case report.JsonV2:
		return uuid.UUID{}, fmt.Errorf("jsonv2 report is not supported on Application Inspector 6.0")
	case report.Markdown:
		aiReportType = v6_0.ReportTypeMd
	case report.Nist:
		aiReportType = v6_0.ReportTypeNist
	case report.Oud4:
		aiReportType = v6_0.ReportTypeOud4
	case report.Owasp:
		aiReportType = v6_0.ReportTypeOwasp
	case report.Owaspm:
		aiReportType = v6_0.ReportTypeOwaspm
	case report.Pcidss:
		aiReportType = v6_0.ReportTypePcidss
	case report.PlainReport:
		aiReportType = v6_0.ReportTypePlainReport
	case report.Sans:
		aiReportType = v6_0.ReportTypeSans
	case report.Sarif:
		aiReportType = v6_0.ReportTypeSarif
	case report.Xml:
		aiReportType = v6_0.ReportTypeXml
	default:
		return uuid.UUID{}, fmt.Errorf("invalid reportType: %s", reportType)
	}

	response, err := a.GetApiReportsTemplatesType(ctx, aiReportType, &params, a.AddJWTToHeader)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("ai adapter get default template id request: %w", err)
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

func (a *ClientAI60) GetCustomTemplateId(ctx context.Context, reportName string) (uuid.UUID, error) {
	localeId := "ru-RU"
	params := v6_0.GetApiReportsUserTemplatesNameParams{
		LocaleId: &localeId,
	}

	response, err := a.GetApiReportsUserTemplatesName(ctx, reportName, &params, a.AddJWTToHeader)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("ai adapter get custom template id request: %w", err)
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

func (a *ClientAI60) GetReport(ctx context.Context, projectId, scanResultId, templateId uuid.UUID, includeComments, includeDFD, includeGlossary bool, l10n string) (io.ReadCloser, error) {
	useFilters := false
	sessionId := uuid.New()

	model := v6_0.ReportGenerateModel{
		LocaleId: &l10n,
		Parameters: &v6_0.UserReportParametersModel{
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

	response, err := a.PostApiReportsGenerate(ctx, model, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter generate report request: %w", err)
	}

	if err = CheckResponse(response, "report"); err != nil {
		return nil, fmt.Errorf("ai adapter generate report: %w", err)
	}

	return response.Body, nil
}

func (a *ClientAI60) GetSbom(ctx context.Context, projectId, scanResultId uuid.UUID) (io.ReadCloser, error) {
	response, err := a.GetApiStoreProjectIdSbomsScanResultId(ctx, projectId, scanResultId, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get sbom: %w", err)
	}

	if err = CheckResponse(response, "sbom"); err != nil {
		return nil, fmt.Errorf("ai adapter get sbom: %w", err)
	}

	return response.Body, nil
}

func (a *ClientAI60) GetScanLogs(ctx context.Context, projectId, scanResultId uuid.UUID) (io.ReadCloser, error) {
	response, err := a.GetApiStoreProjectIdLogsScanResultId(ctx, projectId, scanResultId, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get scan logs request: %w", err)
	}

	if err = CheckResponse(response, "logs"); err != nil {
		return nil, fmt.Errorf("ai adapter get scan logs: %w", err)
	}

	return response.Body, nil
}

func (a *ClientAI60) GetBranches(ctx context.Context, projectId uuid.UUID) ([]branch.Branch, error) {
	getBranchesResponse, err := a.GetApiProjectsProjectIdBranchesWithResponse(ctx, projectId, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get branch request: %w", err)
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

func (a *ClientAI60) GetScans(ctx context.Context, branchId uuid.UUID) ([]scan.Scan, error) {
	response, err := a.GetApiBranchesBranchIdScanResultsWithResponse(ctx, branchId, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get scan results request: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	errorModel := response.JSON400
	if err = CheckResponseByModel(statusCode, body, errorModel); err != nil {
		return nil, fmt.Errorf("ai adapter get scans: %w", err)
	}

	models := *response.JSON200
	scans := make([]scan.Scan, 0, len(models))
	for _, model := range models {
		if model.Id == nil || model.SettingsId == nil {
			continue
		}

		scans = append(scans, scan.NewScan(uuid.UUID(*model.Id), uuid.UUID(*model.SettingsId), model.ScanDate, model.ScanLabel))
	}

	return scans, nil
}

func (a *ClientAI60) GetLastScan(ctx context.Context, branchId uuid.UUID) (*scan.Scan, error) {
	response, err := a.GetApiBranchesBranchIdScanResultsLastWithResponse(ctx, branchId, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get last scan result request: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	errorModel := response.JSON400
	if err = CheckResponseByModel(statusCode, body, errorModel); err != nil {
		return nil, fmt.Errorf("ai adapter get last scan result: %w", err)
	}

	model := response.JSON200
	if model.Id == nil || model.SettingsId == nil {
		return nil, apperror.NewEmptyResponseError("last scan result")
	}

	scanResult := scan.NewScan(uuid.UUID(*model.Id), uuid.UUID(*model.SettingsId), model.ScanDate, model.ScanLabel)

	return &scanResult, nil
}

func (a *ClientAI60) GetScan(ctx context.Context, projectId, scanId uuid.UUID) (*scan.Scan, error) {
	response, err := a.GetApiProjectsProjectIdScanResultsScanResultIdWithResponse(ctx, projectId, scanId, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get scan request: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	errorModel := response.JSON400
	if err = CheckResponseByModel(statusCode, body, errorModel); err != nil {
		return nil, fmt.Errorf("ai adapter get scan aiproj: %w", err)
	}

	model := response.JSON200
	if model.Id == nil || model.SettingsId == nil {
		return nil, apperror.NewEmptyResponseError("scan")
	}

	scanResult := scan.NewScan(uuid.UUID(*model.Id), uuid.UUID(*model.SettingsId), model.ScanDate, model.ScanLabel)

	return &scanResult, nil
}

func (a *ClientAI60) GetProjectAiproj(ctx context.Context, projectId uuid.UUID) (io.ReadCloser, error) {
	response, err := a.GetApiProjectsProjectIdAiproj(ctx, projectId, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get project aiproj request: %w", err)
	}

	if err = CheckResponse(response, "aiproj"); err != nil {
		return nil, fmt.Errorf("ai adapter get project aiproj: %w", err)
	}

	return response.Body, nil
}

func (a *ClientAI60) GetScanAiproj(ctx context.Context, projectId, scanSettingsId uuid.UUID) (io.ReadCloser, error) {
	response, err := a.GetApiProjectsProjectIdScanSettingsScanSettingsIdAiproj(ctx, projectId, scanSettingsId, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get aiproj request: %w", err)
	}

	if err = CheckResponse(response, "aiproj"); err != nil {
		return nil, fmt.Errorf("ai adapter get aiproj: %w", err)
	}

	return response.Body, nil
}

func scanProgressFromModel(model *v6_0.ScanProgressModel) (scanstage.ScanStage, bool) {
	if model == nil || model.Stage == nil {
		return scanstage.ScanStage{}, false
	}

	stage := scanstage.ScanStage{
		Stage: string(*model.Stage),
	}
	if model.Value != nil {
		stage.Value = *model.Value
	}
	if model.SubStage != nil {
		stage.SubStage = *model.SubStage
	}

	return stage, true
}

func (a *ClientAI60) GetScanStage(ctx context.Context, projectId, scanId uuid.UUID) (scanstage.ScanStage, error) {
	response, err := a.GetApiProjectsProjectIdScanResultsScanResultIdProgressWithResponse(ctx, projectId, scanId, a.AddJWTToHeader)
	if err == nil {
		statusCode := response.StatusCode()
		body := string(response.Body)
		errorModel := response.JSON400
		if err = CheckResponseByModel(statusCode, body, errorModel); err == nil {
			if stage, ok := scanProgressFromModel(response.JSON200); ok {
				return stage, nil
			}
		}
	}

	scanResponse, err := a.GetApiProjectsProjectIdScanResultsScanResultIdWithResponse(ctx, projectId, scanId, a.AddJWTToHeader)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return scanstage.ScanStage{}, err
		}

		return scanstage.ScanStage{}, fmt.Errorf("ai adapter get scan result request: %w", err)
	}

	statusCode := scanResponse.StatusCode()
	body := string(scanResponse.Body)
	errorModel := scanResponse.JSON400
	if err = CheckResponseByModel(statusCode, body, errorModel); err != nil {
		return scanstage.ScanStage{}, fmt.Errorf("ai adapter get scan result: %w", err)
	}

	if stage, ok := scanProgressFromModel(response.JSON200); ok {
		return stage, nil
	}

	return scanstage.ScanStage{}, apperror.NewEmptyResponseError("scan progress")
}

func (a *ClientAI60) GetScanItem(ctx context.Context, id uuid.UUID) (queue.Item, error) {
	if item, found, err := a.findScanQueueItem(ctx, id); err != nil {
		return queue.Item{}, err
	} else if found {
		return item, nil
	}

	return a.getScanQueueItemByID(ctx, id)
}

func (a *ClientAI60) getScanQueueItemByID(ctx context.Context, id uuid.UUID) (queue.Item, error) {
	response, err := a.GetItemWithResponse(ctx, id, a.AddJWTToHeader)
	if err != nil {
		return queue.Item{}, fmt.Errorf("ai adapter get scan queue item request: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	if err = CheckResponseByModel(statusCode, body, response.JSON400); err != nil {
		if apperror.IsApiErrorCode(err, string(v6_0.ApiErrorTypeEMPTYSCANRESULT)) {
			time.Sleep(1 * time.Second)
			return a.getScanQueueItemByID(ctx, id)
		}
		if apperror.IsApiErrorCode(err, string(v6_0.ApiErrorTypeQUEUEITEMNOTFOUND)) ||
			apperror.IsApiErrorCode(err, string(v6_0.ApiErrorTypeQUEUEITEMALREADYASSIGNEDTOAGENT)) {
			return queue.Item{ScanId: id}, nil
		}

		return queue.Item{}, fmt.Errorf("ai adapter get scan queue item: %w", err)
	}

	if response.JSON200 == nil {
		return queue.Item{}, apperror.NewEmptyResponseError("scan queue item")
	}

	dto := response.JSON200

	return queue.Item{
		Place:  int(dto.Position),
		ScanId: dto.ScanResultId,
	}, nil
}

func (a *ClientAI60) findScanQueueItem(ctx context.Context, id uuid.UUID) (queue.Item, bool, error) {
	response, err := a.GetAllItemsWithResponse(ctx, a.AddJWTToHeader)
	if err != nil {
		return queue.Item{}, false, fmt.Errorf("ai adapter get scan queue items request: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	if err = CheckResponseByModel(statusCode, body, response.JSON400); err != nil {
		return queue.Item{}, false, fmt.Errorf("ai adapter get scan queue items: %w", err)
	}

	if response.JSON200 == nil {
		return queue.Item{}, false, apperror.NewEmptyResponseError("scan queue items")
	}

	models := *response.JSON200
	sort.Slice(models, func(i, j int) bool {
		return models[i].QueueDateTime < models[j].QueueDateTime
	})

	for i, model := range models {
		if id == model.Id || id == model.ScanResultId {
			return queue.Item{
				Place:  i + 1,
				OutOf:  len(models),
				ScanId: model.ScanResultId,
			}, true, nil
		}
	}

	return queue.Item{}, false, nil
}

func (a *ClientAI60) createScanQueueItem(ctx context.Context, branchId uuid.UUID, scanLabel string, scanType scantype.Type) (uuid.UUID, error) {
	scope, err := toScope(scanType)
	if err != nil {
		return uuid.UUID{}, err
	}

	params := v6_0.CreateQueueItem{
		BranchId:  branchId,
		Scope:     scope,
		ScanLabel: &scanLabel,
	}

	response, err := a.CreateItemWithResponse(ctx, params, a.AddJWTToHeader)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("create scan queue item request: %w", err)
	}

	statusCode := response.StatusCode()
	responseBody := string(response.Body)
	if err := CheckResponseByModel(statusCode, responseBody, response.JSON400); err != nil {
		return uuid.UUID{}, fmt.Errorf("ai adapter start scan: %w", err)
	}

	if response.JSON200 == nil {
		return uuid.UUID{}, apperror.NewEmptyResponseError("scan result id")
	}

	return *response.JSON200, nil
}

func (a *ClientAI60) StartScanBranch(ctx context.Context, branchId uuid.UUID, scanLabel string, scanType scantype.Type) (uuid.UUID, error) {
	taskId, err := a.createScanQueueItem(ctx, branchId, scanLabel, scanType)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("ai adapter start scan: %w", err)
	}

	item, err := a.GetScanItem(ctx, taskId)
	if err != nil {
		return uuid.UUID{}, err
	}
	if item.ScanId != uuid.Nil && item.ScanId != taskId {
		return item.ScanId, nil
	}
	if item.OutOf > 0 {
		return item.ScanId, nil
	}

	lastScan, lastErr := a.GetLastScan(ctx, branchId)
	if lastErr != nil {
		return uuid.UUID{}, fmt.Errorf("scan queue item unavailable, get last scan: %w", lastErr)
	}

	return lastScan.Id, nil
}

func (a *ClientAI60) StartScanProject(ctx context.Context, projectId uuid.UUID, scanLabel string, scanType scantype.Type) (uuid.UUID, error) {
	branches, err := a.GetBranches(ctx, projectId)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("get branches for project scan: %w", err)
	}
	if len(branches) == 0 {
		return uuid.UUID{}, fmt.Errorf("no branches found for project")
	}

	branchId := branches[0].Id
	for _, b := range branches {
		if b.IsWorking {
			branchId = b.Id
			break
		}
	}

	return a.createScanQueueItem(ctx, branchId, scanLabel, scanType)
}

func toScope(scanType scantype.Type) (v6_0.Scope, error) {
	switch scanType {
	case scantype.Incremental:
		return v6_0.ScopeIncremental, nil
	case scantype.Full:
		return v6_0.ScopeFull, nil
	default:
		return "", fmt.Errorf("invalid scan type: %d", scanType)
	}
}

func (a *ClientAI60) StopScan(ctx context.Context, scanResultId uuid.UUID) error {
	response, err := a.StopScanWithResponse(ctx, scanResultId, a.AddJWTToHeader)
	if err != nil {
		return fmt.Errorf("ai adapter stop scan request: %w", err)
	}

	statusCode := response.StatusCode()
	responseBody := string(response.Body)
	if err = CheckResponseByModel(statusCode, responseBody, response.JSON400); err != nil {
		return fmt.Errorf("ai adapter stop scan: %w", err)
	}

	return nil
}

func (a *ClientAI60) UpdateSources(ctx context.Context, projectId, branchId uuid.UUID, scanTargetPath string, exclusions gitignore.Exclusions) error {
	log := logger.FromContext(ctx)

	archivePath, err := common.PrepareArchive(ctx, scanTargetPath, exclusions)
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

	body, contentType, err := common.PrepareMultipartBody(ctx, archivePath, true)
	if err != nil {
		return err
	}
	defer func() {
		_ = body.Close()
	}()

	archived := true
	params := v6_0.PostApiStoreProjectIdBranchesBranchIdSourcesParams{Archived: &archived}

	response, err := a.PostApiStoreProjectIdBranchesBranchIdSourcesWithBodyWithResponse(ctx, projectId, branchId, &params, contentType, body, a.AddJWTToHeader)
	if err != nil {
		return fmt.Errorf("ai update sources request: %w", err)
	}

	statusCode := response.StatusCode()
	responseBody := string(response.Body)
	errorModel := response.JSON400
	if err = CheckResponseByModel(statusCode, responseBody, errorModel); err != nil {
		return fmt.Errorf("ai update sources post sources: %w", err)
	}

	return nil
}

func (a *ClientAI60) GetVersion(ctx context.Context) (version.Version, error) {
	response, err := a.GetApiVersionsPackageCurrentWithResponse(ctx, a.AddJWTToHeader)
	if err != nil {
		return version.Version{}, fmt.Errorf("ai get version request: %w", err)
	}

	statusCode := response.StatusCode()
	responseBody := string(response.Body)
	if err = CheckResponseByModel(statusCode, responseBody, nil); err != nil {
		return version.Version{}, fmt.Errorf("ai get version: %w", err)
	}

	v, err := version.NewVersion(responseBody)
	if err != nil {
		return version.Version{}, fmt.Errorf("new version: %w", err)
	}

	return v, nil
}

func (a *ClientAI60) GetHealthcheck(ctx context.Context) (bool, error) {
	response, err := a.GetHealthSummaryWithResponse(ctx, a.AddJWTToHeader)
	if err != nil {
		return false, fmt.Errorf("ai get version request: %w", err)
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

func (a *ClientAI60) CheckLicense(ctx context.Context) error {
	response, err := a.GetApiLicenseWithResponse(ctx, a.AddJWTToHeader)
	if err != nil {
		return fmt.Errorf("ai check license request: %w", err)
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

func (a *ClientAI60) GetScanStatistic(ctx context.Context, projectId, scanResultId uuid.UUID) (*statistic.Statistic, error) {
	response, err := a.GetApiProjectsProjectIdScanResultsScanResultIdStatisticWithResponse(ctx, projectId, scanResultId, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai get scan statistic request: %w", err)
	}

	statusCode := response.StatusCode()
	responseBody := string(response.Body)
	if err = CheckResponseByModel(statusCode, responseBody, nil); err != nil {
		return nil, fmt.Errorf("ai get scan statistic: %w", err)
	}

	model := response.JSON200

	return &statistic.Statistic{
		Total:     *model.Total,
		High:      *model.High,
		Medium:    *model.Medium,
		Low:       *model.Low,
		Potential: *model.Potential,
	}, nil
}
