package port

import (
	"context"

	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/scan"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanstage"
	"github.com/google/uuid"
)

const (
	SarifReportType = "SARIF"
	PlainReportType = "PLAIN_REPORT"
)

type Ai interface {
	DeleteProject(ctx context.Context, projectId uuid.UUID) error

	GetProjects(ctx context.Context) ([]project.Project, error)

	GetScan(ctx context.Context, projectId, scanId uuid.UUID) (*scan.Scan, error)
	GetScanAiproj(ctx context.Context, projectId, scanSettingsId uuid.UUID) (string, error)
	GetScanStage(ctx context.Context, projectId, scanId uuid.UUID) (scanstage.ScanStage, error)
	GetScans(ctx context.Context, branchId uuid.UUID) ([]scan.Scan, error)

	StartScanBranch(ctx context.Context, branchId uuid.UUID) (uuid.UUID, error)
	StartScanProject(ctx context.Context, projectId uuid.UUID) (uuid.UUID, error)

	StopScan(ctx context.Context, scanResultId uuid.UUID) error

	UpdateSources(ctx context.Context, projectId, branchId uuid.UUID, sourcePath string) error
}
