package ai

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	"github.com/POSIdev-community/aictl/internal/core/domain/branch"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/policystate"
	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/queue"
	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/POSIdev-community/aictl/internal/core/domain/scan"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanagent"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanstage"
	"github.com/POSIdev-community/aictl/internal/core/domain/scantype"
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/internal/core/domain/statistic"
	"github.com/POSIdev-community/aictl/internal/core/domain/version"
	"github.com/POSIdev-community/aictl/pkg/gitignore"
)

type ClientAi = common.ClientAi

type jwtRetrier interface {
	AddJwtRetry()
}

type Adapter struct {
	baseClient *common.BaseClient

	activeClient  ClientAi
	cfg           *config.Config
	serverVersion version.Version
}

func NewAdapter(cfg *config.Config) (*Adapter, error) {
	baseClient := common.NewBaseClient()

	return &Adapter{baseClient: baseClient, cfg: cfg}, nil
}

func (a *Adapter) InitializeWithRetry(ctx context.Context) error {
	if err := a.Initialize(ctx); err != nil {
		return fmt.Errorf("initialize ai adapter: %w", err)
	}

	if r, ok := a.activeClient.(jwtRetrier); ok {
		r.AddJwtRetry()
	}

	return nil
}

func (a *Adapter) GetDefaultSettings(ctx context.Context) (settings.ScanSettings, error) {
	return a.activeClient.GetDefaultSettings(ctx)
}

func (a *Adapter) GetProjectSettings(ctx context.Context, projectId uuid.UUID) (settings.ScanSettings, error) {
	return a.activeClient.GetProjectSettings(ctx, projectId)
}

func (a *Adapter) SetProjectSettings(ctx context.Context, projectId uuid.UUID, settings *settings.ScanSettings) error {
	return a.activeClient.SetProjectSettings(ctx, projectId, settings)
}

func (a *Adapter) GetProjectPolicies(ctx context.Context, projectId uuid.UUID) (io.ReadCloser, error) {
	return a.activeClient.GetProjectPolicies(ctx, projectId)
}

func (a *Adapter) SetProjectPolicies(ctx context.Context, projectId uuid.UUID, rawJSON []byte) error {
	return a.activeClient.SetProjectPolicies(ctx, projectId, rawJSON)
}

func (a *Adapter) GetProjectExclusions(ctx context.Context, projectId uuid.UUID) (string, error) {
	return a.activeClient.GetProjectExclusions(ctx, projectId)
}

func (a *Adapter) SetProjectExclusions(ctx context.Context, projectId uuid.UUID, exclusions string) error {
	return a.activeClient.SetProjectExclusions(ctx, projectId, exclusions)
}

func (a *Adapter) UpdateProjectLanguages(ctx context.Context, projectId uuid.UUID) error {
	return a.activeClient.UpdateProjectLanguages(ctx, projectId)
}

func (a *Adapter) GetScanAgents(ctx context.Context) ([]scanagent.ScanAgent, error) {
	return a.activeClient.GetScanAgents(ctx)
}

func (a *Adapter) CreateBranch(ctx context.Context, projectId uuid.UUID, branchName, scanTargetPath string, exclusions gitignore.Exclusions, tempDir string) (*uuid.UUID, error) {
	return a.activeClient.CreateBranch(ctx, projectId, branchName, scanTargetPath, exclusions, tempDir)
}

func (a *Adapter) CreateProject(ctx context.Context, projectName string) (*uuid.UUID, error) {
	return a.activeClient.CreateProject(ctx, projectName)
}

func (a *Adapter) CreateSbomProject(ctx context.Context, projectName string) (*uuid.UUID, error) {
	return a.activeClient.CreateSbomProject(ctx, projectName)
}

func (a *Adapter) DeleteProject(ctx context.Context, projectId uuid.UUID) error {
	return a.activeClient.DeleteProject(ctx, projectId)
}

func (a *Adapter) ExistsProject(ctx context.Context, projectName string) (bool, error) {
	return a.activeClient.ExistsProject(ctx, projectName)
}

func (a *Adapter) GetProjectId(ctx context.Context, projectName string) (*uuid.UUID, error) {
	return a.activeClient.GetProjectId(ctx, projectName)
}

func (a *Adapter) GetProjectByName(ctx context.Context, projectName string) (*project.Project, error) {
	return a.activeClient.GetProjectByName(ctx, projectName)
}

func (a *Adapter) GetProjects(ctx context.Context) ([]project.Project, error) {
	return a.activeClient.GetProjects(ctx)
}

func (a *Adapter) GetProject(ctx context.Context, projectId uuid.UUID) (*project.Project, error) {
	return a.activeClient.GetProject(ctx, projectId)
}

func (a *Adapter) UpdateSbom(ctx context.Context, projectId uuid.UUID, sbomPath string) error {
	return a.activeClient.UpdateSbom(ctx, projectId, sbomPath)
}

func (a *Adapter) GetDefaultTemplateId(ctx context.Context, reportType report.ReportType) (uuid.UUID, error) {
	return a.activeClient.GetDefaultTemplateId(ctx, reportType)
}

func (a *Adapter) GetCustomTemplateId(ctx context.Context, reportName string) (uuid.UUID, error) {
	return a.activeClient.GetCustomTemplateId(ctx, reportName)
}

func (a *Adapter) GetReportTemplates(ctx context.Context, localization string) ([]report.Template, error) {
	return a.activeClient.GetReportTemplates(ctx, localization)
}

func (a *Adapter) GetReport(ctx context.Context, projectId, scanResultId, templateId uuid.UUID, includeComments, includeDFD, includeGlossary bool, l10n string, filters report.Filters) (io.ReadCloser, error) {
	if err := common.ValidateReportFiltersForVersion(filters, a.serverVersion); err != nil {
		return nil, err
	}
	return a.activeClient.GetReport(ctx, projectId, scanResultId, templateId, includeComments, includeDFD, includeGlossary, l10n, filters)
}

func (a *Adapter) GetSbom(ctx context.Context, projectId, scanResultId uuid.UUID) (io.ReadCloser, error) {
	return a.activeClient.GetSbom(ctx, projectId, scanResultId)
}

func (a *Adapter) GetScanLogs(ctx context.Context, projectId, scanResultId uuid.UUID) (io.ReadCloser, error) {
	return a.activeClient.GetScanLogs(ctx, projectId, scanResultId)
}

func (a *Adapter) GetScanErrors(ctx context.Context, projectId, scanResultId uuid.UUID) ([]string, error) {
	return a.activeClient.GetScanErrors(ctx, projectId, scanResultId)
}

func (a *Adapter) GetBranches(ctx context.Context, projectId uuid.UUID) ([]branch.Branch, error) {
	return a.activeClient.GetBranches(ctx, projectId)
}

func (a *Adapter) GetBranch(ctx context.Context, branchId uuid.UUID) (*branch.Branch, error) {
	return a.activeClient.GetBranch(ctx, branchId)
}

func (a *Adapter) GetScans(ctx context.Context, branchId uuid.UUID) ([]scan.Scan, error) {
	return a.activeClient.GetScans(ctx, branchId)
}

func (a *Adapter) GetLastScan(ctx context.Context, branchId uuid.UUID) (*scan.Scan, error) {
	return a.activeClient.GetLastScan(ctx, branchId)
}

func (a *Adapter) GetScan(ctx context.Context, projectId, scanId uuid.UUID) (*scan.Scan, error) {
	return a.activeClient.GetScan(ctx, projectId, scanId)
}

func (a *Adapter) GetProjectAiproj(ctx context.Context, projectId uuid.UUID) (io.ReadCloser, error) {
	return a.activeClient.GetProjectAiproj(ctx, projectId)
}

func (a *Adapter) GetScanAiproj(ctx context.Context, projectId, scanSettingsId uuid.UUID) (io.ReadCloser, error) {
	return a.activeClient.GetScanAiproj(ctx, projectId, scanSettingsId)
}

func (a *Adapter) GetScanStage(ctx context.Context, projectId, scanId uuid.UUID) (scanstage.ScanStage, error) {
	return a.activeClient.GetScanStage(ctx, projectId, scanId)
}

func (a *Adapter) GetScanPolicyState(ctx context.Context, projectId, scanId uuid.UUID) (policystate.State, error) {
	return a.activeClient.GetScanPolicyState(ctx, projectId, scanId)
}

func (a *Adapter) GetScanItem(ctx context.Context, id uuid.UUID) (queue.Item, error) {
	return a.activeClient.GetScanItem(ctx, id)
}

func (a *Adapter) GetScanQueue(ctx context.Context) ([]queue.Entry, error) {
	return a.activeClient.GetScanQueue(ctx)
}

func (a *Adapter) GetActiveScans(ctx context.Context) ([]queue.Entry, error) {
	return a.activeClient.GetActiveScans(ctx)
}

func (a *Adapter) StartScanBranch(ctx context.Context, branchId uuid.UUID, scanLabel string, scanType scantype.Type) (uuid.UUID, error) {
	return a.activeClient.StartScanBranch(ctx, branchId, scanLabel, scanType)
}

func (a *Adapter) StartScanProject(ctx context.Context, projectId uuid.UUID, scanLabel string, scanType scantype.Type) (uuid.UUID, error) {
	return a.activeClient.StartScanProject(ctx, projectId, scanLabel, scanType)
}

func (a *Adapter) StartScanSbom(ctx context.Context, projectId uuid.UUID, scanLabel string) (uuid.UUID, error) {
	return a.activeClient.StartScanSbom(ctx, projectId, scanLabel)
}

func (a *Adapter) StopScan(ctx context.Context, scanResultId uuid.UUID) error {
	return a.activeClient.StopScan(ctx, scanResultId)
}

func (a *Adapter) UpdateSources(ctx context.Context, projectId, branchId uuid.UUID, scanTargetPath string, exclusions gitignore.Exclusions, tempDir string) error {
	return a.activeClient.UpdateSources(ctx, projectId, branchId, scanTargetPath, exclusions, tempDir)
}

func (a *Adapter) GetVersion(ctx context.Context) (version.Version, error) {
	return a.serverVersion, nil
}

func (a *Adapter) GetHealthcheck(ctx context.Context) (bool, error) {
	return a.activeClient.GetHealthcheck(ctx)
}

func (a *Adapter) CheckLicense(ctx context.Context) error {
	return a.activeClient.CheckLicense(ctx)
}

func (a *Adapter) GetScanStatistic(ctx context.Context, projectId, scanResultId uuid.UUID) (*statistic.Statistic, error) {
	return a.activeClient.GetScanStatistic(ctx, projectId, scanResultId)
}
