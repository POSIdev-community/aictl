package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/POSIdev-community/aictl/internal/core/domain/branch"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/scan"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanstage"
	"github.com/POSIdev-community/aictl/internal/core/port"
	. "github.com/POSIdev-community/aictl/pkg/clientai"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/POSIdev-community/aictl/pkg/logger"
	"github.com/google/uuid"
	"io"
	"os"
	"sort"
	"strings"
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
