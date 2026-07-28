package common

import (
	"context"
	"io"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/branch"
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
	"github.com/POSIdev-community/aictl/pkg/gitignore"
)

type ClientAi interface {
	GetDefaultSettings(ctx context.Context) (settings.ScanSettings, error)
	GetProjectSettings(ctx context.Context, projectId uuid.UUID) (settings.ScanSettings, error)
	SetProjectSettings(ctx context.Context, projectId uuid.UUID, settings *settings.ScanSettings) error
	GetProjectPolicies(ctx context.Context, projectId uuid.UUID) (io.ReadCloser, error)
	SetProjectPolicies(ctx context.Context, projectId uuid.UUID, rawJSON []byte) error
	GetScanAgents(ctx context.Context) ([]scanagent.ScanAgent, error)
	CreateBranch(ctx context.Context, projectId uuid.UUID, branchName, scanTargetPath string, exclusions gitignore.Exclusions) (*uuid.UUID, error)
	CreateProject(ctx context.Context, projectName string) (*uuid.UUID, error)
	DeleteProject(ctx context.Context, projectId uuid.UUID) error
	ExistsProject(ctx context.Context, projectName string) (bool, error)
	GetProjectId(ctx context.Context, projectName string) (*uuid.UUID, error)
	GetProjects(ctx context.Context) ([]project.Project, error)
	GetProject(ctx context.Context, projectId uuid.UUID) (*project.Project, error)
	GetDefaultTemplateId(ctx context.Context, reportType report.ReportType) (uuid.UUID, error)
	GetCustomTemplateId(ctx context.Context, reportName string) (uuid.UUID, error)
	GetReport(ctx context.Context, projectId, scanResultId, templateId uuid.UUID, includeComments, includeDFD, includeGlossary bool, l10n string) (io.ReadCloser, error)
	GetSbom(ctx context.Context, projectId, scanResultId uuid.UUID) (io.ReadCloser, error)
	GetScanLogs(ctx context.Context, projectId, scanResultId uuid.UUID) (io.ReadCloser, error)
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
	StartScanBranch(ctx context.Context, branchId uuid.UUID, scanLabel string, scanType scantype.Type) (uuid.UUID, error)
	StartScanProject(ctx context.Context, projectId uuid.UUID, scanLabel string, scanType scantype.Type) (uuid.UUID, error)
	StopScan(ctx context.Context, scanResultId uuid.UUID) error
	UpdateSources(ctx context.Context, projectId, branchId uuid.UUID, scanTargetPath string, exclusions gitignore.Exclusions) error
	GetHealthcheck(ctx context.Context) (bool, error)
	CheckLicense(ctx context.Context) error
	GetScanStatistic(ctx context.Context, projectId, scanResultId uuid.UUID) (*statistic.Statistic, error)
}
