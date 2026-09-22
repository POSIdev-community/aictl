package common

import (
	"context"
	"io"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/branch"
	"github.com/POSIdev-community/aictl/internal/core/domain/license"
	"github.com/POSIdev-community/aictl/internal/core/domain/policystate"
	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/queue"
	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/POSIdev-community/aictl/internal/core/domain/scafeeds"
	"github.com/POSIdev-community/aictl/internal/core/domain/scan"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanagent"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanstage"
	"github.com/POSIdev-community/aictl/internal/core/domain/scantype"
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/internal/core/domain/statistic"
	"github.com/POSIdev-community/aictl/pkg/gitignore"
)

type ClientAi interface {
	GetDefaultSettings(ctx context.Context) (settings.ScanSettings, error)
	GetProjectSettings(ctx context.Context, projectId uuid.UUID) (settings.ScanSettings, error)
	SetProjectSettings(ctx context.Context, projectId uuid.UUID, settings *settings.ScanSettings) error
	GetProjectPolicies(ctx context.Context, projectId uuid.UUID) (io.ReadCloser, error)
	SetProjectPolicies(ctx context.Context, projectId uuid.UUID, rawJSON []byte) error
	GetProjectExclusions(ctx context.Context, projectId uuid.UUID) (string, error)
	SetProjectExclusions(ctx context.Context, projectId uuid.UUID, exclusions string) error
	UpdateProjectLanguages(ctx context.Context, projectId uuid.UUID) error
	GetScanAgents(ctx context.Context) ([]scanagent.ScanAgent, error)
	CreateBranch(ctx context.Context, projectId uuid.UUID, branchName, scanTargetPath string, exclusions gitignore.Exclusions, tempDir string) (*uuid.UUID, error)
	CreateProject(ctx context.Context, projectName string) (*uuid.UUID, error)
	CreateSbomProject(ctx context.Context, projectName string) (*uuid.UUID, error)
	DeleteProject(ctx context.Context, projectId uuid.UUID) error
	ExistsProject(ctx context.Context, projectName string) (bool, error)
	GetProjectId(ctx context.Context, projectName string) (*uuid.UUID, error)
	GetProjectByName(ctx context.Context, projectName string) (*project.Project, error)
	GetProjects(ctx context.Context) ([]project.Project, error)
	GetProject(ctx context.Context, projectId uuid.UUID) (*project.Project, error)
	UpdateSbom(ctx context.Context, projectId uuid.UUID, sbomPath string) error
	UpdateScaFeeds(ctx context.Context, path, version string) error
	GetScaFeeds(ctx context.Context, statuses []scafeeds.Status) ([]scafeeds.Package, error)
	DownloadScaFeeds(ctx context.Context, version string) (io.ReadCloser, string, error)
	RollbackScaFeeds(ctx context.Context) (scafeeds.Package, error)
	StartScanSbom(ctx context.Context, projectId uuid.UUID, scanLabel string) (uuid.UUID, error)
	GetDefaultTemplateId(ctx context.Context, reportType report.ReportType) (uuid.UUID, error)
	GetCustomTemplateId(ctx context.Context, reportName string) (uuid.UUID, error)
	GetReportTemplates(ctx context.Context, localization string) ([]report.Template, error)
	GetReport(ctx context.Context, projectId, scanResultId, templateId uuid.UUID, includeComments, includeDFD, includeGlossary bool, l10n string, filters report.Filters) (io.ReadCloser, error)
	GetSbom(ctx context.Context, projectId, scanResultId uuid.UUID) (io.ReadCloser, error)
	GetScanLogs(ctx context.Context, projectId, scanResultId uuid.UUID) (io.ReadCloser, error)
	GetScanErrors(ctx context.Context, projectId, scanResultId uuid.UUID) ([]string, error)
	GetBranches(ctx context.Context, projectId uuid.UUID) ([]branch.Branch, error)
	GetBranch(ctx context.Context, branchId uuid.UUID) (*branch.Branch, error)
	GetScans(ctx context.Context, branchId uuid.UUID) ([]scan.Scan, error)
	GetLastScan(ctx context.Context, branchId uuid.UUID) (*scan.Scan, error)
	GetScan(ctx context.Context, projectId, scanId uuid.UUID) (*scan.Scan, error)
	GetProjectAiproj(ctx context.Context, projectId uuid.UUID) (io.ReadCloser, error)
	GetScanAiproj(ctx context.Context, projectId, scanSettingsId uuid.UUID) (io.ReadCloser, error)
	GetScanStage(ctx context.Context, projectId, scanId uuid.UUID) (scanstage.ScanStage, error)
	GetScanPolicyState(ctx context.Context, projectId, scanId uuid.UUID) (policystate.State, error)
	GetScanItem(ctx context.Context, id uuid.UUID) (queue.Item, error)
	GetScanQueue(ctx context.Context) ([]queue.Entry, error)
	GetActiveScans(ctx context.Context) ([]queue.Entry, error)
	StartScanBranch(ctx context.Context, branchId uuid.UUID, scanLabel string, scanType scantype.Type) (uuid.UUID, error)
	StartScanProject(ctx context.Context, projectId uuid.UUID, scanLabel string, scanType scantype.Type) (uuid.UUID, error)
	StopScan(ctx context.Context, scanResultId uuid.UUID) error
	UpdateSources(ctx context.Context, projectId, branchId uuid.UUID, scanTargetPath string, exclusions gitignore.Exclusions, tempDir string) error
	GetHealthcheck(ctx context.Context) (bool, error)
	// CheckLicense validates IsValid and returns the mapped license for caching.
	CheckLicense(ctx context.Context) (*license.License, error)
	GetScanStatistic(ctx context.Context, projectId, scanResultId uuid.UUID) (*statistic.Statistic, error)
	GetScanIssues(ctx context.Context, projectId, scanResultId uuid.UUID) ([]statistic.Issue, error)
}
