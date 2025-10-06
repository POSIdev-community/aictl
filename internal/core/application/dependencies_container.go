package application

import (
	"context"
	"github.com/POSIdev-community/aictl/internal/adapter/ai"
	"github.com/POSIdev-community/aictl/internal/adapter/cli"
	"github.com/POSIdev-community/aictl/internal/adapter/config"

	configClear "github.com/POSIdev-community/aictl/internal/core/application/usecase/config/clear"
	configSet "github.com/POSIdev-community/aictl/internal/core/application/usecase/config/set"
	configShow "github.com/POSIdev-community/aictl/internal/core/application/usecase/config/show"
	configUnset "github.com/POSIdev-community/aictl/internal/core/application/usecase/config/unset"
	createBranch "github.com/POSIdev-community/aictl/internal/core/application/usecase/create/branch"
	createProject "github.com/POSIdev-community/aictl/internal/core/application/usecase/create/project"
	deleteProjects "github.com/POSIdev-community/aictl/internal/core/application/usecase/delete/projects"
	getProjects "github.com/POSIdev-community/aictl/internal/core/application/usecase/get/projects"
	getReports "github.com/POSIdev-community/aictl/internal/core/application/usecase/get/reports"
	getScan "github.com/POSIdev-community/aictl/internal/core/application/usecase/get/scan"
	getScanAiproj "github.com/POSIdev-community/aictl/internal/core/application/usecase/get/scan/aiproj"
	getScanLogs "github.com/POSIdev-community/aictl/internal/core/application/usecase/get/scan/logs"
	getScanReport "github.com/POSIdev-community/aictl/internal/core/application/usecase/get/scan/report"
	getScanReportPlain "github.com/POSIdev-community/aictl/internal/core/application/usecase/get/scan/report/plain"
	getScanReportSarif "github.com/POSIdev-community/aictl/internal/core/application/usecase/get/scan/report/sarif"
	getScanResult "github.com/POSIdev-community/aictl/internal/core/application/usecase/get/scan/result"
	getScanSbom "github.com/POSIdev-community/aictl/internal/core/application/usecase/get/scan/sbom"
	getScanState "github.com/POSIdev-community/aictl/internal/core/application/usecase/get/scan/state"
	getScans "github.com/POSIdev-community/aictl/internal/core/application/usecase/get/scans"
	scanAwait "github.com/POSIdev-community/aictl/internal/core/application/usecase/scan/await"
	scanStartBranch "github.com/POSIdev-community/aictl/internal/core/application/usecase/scan/start/branch"
	scanStartProject "github.com/POSIdev-community/aictl/internal/core/application/usecase/scan/start/project"
	scanStop "github.com/POSIdev-community/aictl/internal/core/application/usecase/scan/stop"
	setSettings "github.com/POSIdev-community/aictl/internal/core/application/usecase/set/settings"
	updateSources "github.com/POSIdev-community/aictl/internal/core/application/usecase/update/sources"
	updateSourcesGit "github.com/POSIdev-community/aictl/internal/core/application/usecase/update/sources/git"
	domainConfig "github.com/POSIdev-community/aictl/internal/core/domain/config"

	"github.com/POSIdev-community/aictl/internal/core/port"
)

type DependenciesContainer struct {
	configAdapter *config.Adapter
}

func NewDependenciesContainer(configAdapter *config.Adapter) *DependenciesContainer {
	return &DependenciesContainer{
		configAdapter,
	}
}

func (c *DependenciesContainer) ConfigClearUseCase() (*configClear.UseCase, error) {
	return getConfigUseCase[configClear.UseCase](c.configAdapter, configClear.NewUseCase)
}

func (c *DependenciesContainer) ConfigSetUseCase() (*configSet.UseCase, error) {
	return getConfigUseCase[configSet.UseCase](c.configAdapter, configSet.NewUseCase)
}

func (c *DependenciesContainer) ConfigShowUseCase() (*configShow.UseCase, error) {
	return getConfigUseCase[configShow.UseCase](c.configAdapter, configShow.NewUseCase)
}

func (c *DependenciesContainer) ConfigUnsetUseCase() (*configUnset.UseCase, error) {
	return getConfigUseCase[configUnset.UseCase](c.configAdapter, configUnset.NewUseCase)
}

func (c *DependenciesContainer) CreateBranchUseCase(ctx context.Context, cfg *domainConfig.Config) (*createBranch.UseCase, error) {
	return getUseCase[createBranch.UseCase](ctx, cfg, createBranch.NewUseCase)
}

func (c *DependenciesContainer) CreateProjectUseCase(ctx context.Context, cfg *domainConfig.Config) (*createProject.UseCase, error) {
	return getUseCase[createProject.UseCase](ctx, cfg, createProject.NewUseCase)
}

func (c *DependenciesContainer) DeleteProjectsUseCase(ctx context.Context, cfg *domainConfig.Config) (*deleteProjects.UseCase, error) {
	return getUseCase[deleteProjects.UseCase](ctx, cfg, deleteProjects.NewUseCase)
}

func (c *DependenciesContainer) GetProjectsUseCase(ctx context.Context, cfg *domainConfig.Config) (*getProjects.UseCase, error) {
	return getUseCase[getProjects.UseCase](ctx, cfg, getProjects.NewUseCase)
}

func (c *DependenciesContainer) GetReportsUseCase(ctx context.Context, cfg *domainConfig.Config) (*getReports.UseCase, error) {
	return getUseCase[getReports.UseCase](ctx, cfg, getReports.NewUseCase)
}

func (c *DependenciesContainer) GetScanUseCase(ctx context.Context, cfg *domainConfig.Config) (*getScan.UseCase, error) {
	return getUseCase[getScan.UseCase](ctx, cfg, getScan.NewUseCase)
}

func (c *DependenciesContainer) GetScanAiprojUseCase(ctx context.Context, cfg *domainConfig.Config) (*getScanAiproj.UseCase, error) {
	return getUseCase[getScanAiproj.UseCase](ctx, cfg, getScanAiproj.NewUseCase)
}

func (c *DependenciesContainer) GetScanLogsUseCase(ctx context.Context, cfg *domainConfig.Config) (*getScanLogs.UseCase, error) {
	return getUseCase[getScanLogs.UseCase](ctx, cfg, getScanLogs.NewUseCase)
}

func (c *DependenciesContainer) GetScanReportUseCase(ctx context.Context, cfg *domainConfig.Config) (*getScanReport.UseCase, error) {
	return getUseCase[getScanReport.UseCase](ctx, cfg, getScanReport.NewUseCase)
}

func (c *DependenciesContainer) GetScanReportPlainUseCase(ctx context.Context, cfg *domainConfig.Config) (*getScanReportPlain.UseCase, error) {
	return getUseCase[getScanReportPlain.UseCase](ctx, cfg, getScanReportPlain.NewUseCase)
}

func (c *DependenciesContainer) GetScanReportSarifUseCase(ctx context.Context, cfg *domainConfig.Config) (*getScanReportSarif.UseCase, error) {
	return getUseCase[getScanReportSarif.UseCase](ctx, cfg, getScanReportSarif.NewUseCase)
}

func (c *DependenciesContainer) GetScanResultUseCase(ctx context.Context, cfg *domainConfig.Config) (*getScanResult.UseCase, error) {
	return getUseCase[getScanResult.UseCase](ctx, cfg, getScanResult.NewUseCase)
}

func (c *DependenciesContainer) GetScanSbomUseCase(ctx context.Context, cfg *domainConfig.Config) (*getScanSbom.UseCase, error) {
	return getUseCase[getScanSbom.UseCase](ctx, cfg, getScanSbom.NewUseCase)
}

func (c *DependenciesContainer) GetScanStateUseCase(ctx context.Context, cfg *domainConfig.Config) (*getScanState.UseCase, error) {
	return getUseCase[getScanState.UseCase](ctx, cfg, getScanState.NewUseCase)
}

func (c *DependenciesContainer) GetScansUseCase(ctx context.Context, cfg *domainConfig.Config) (*getScans.UseCase, error) {
	return getUseCase[getScans.UseCase](ctx, cfg, getScans.NewUseCase)
}

func (c *DependenciesContainer) ScanAwaitUseCase(ctx context.Context, cfg *domainConfig.Config) (*scanAwait.UseCase, error) {
	return getUseCase[scanAwait.UseCase](ctx, cfg, scanAwait.NewUseCase)
}

func (c *DependenciesContainer) ScanStartBranchUseCase(ctx context.Context, cfg *domainConfig.Config) (*scanStartBranch.UseCase, error) {
	return getUseCase[scanStartBranch.UseCase](ctx, cfg, scanStartBranch.NewUseCase)
}

func (c *DependenciesContainer) ScanStartProjectUseCase(ctx context.Context, cfg *domainConfig.Config) (*scanStartProject.UseCase, error) {
	return getUseCase[scanStartProject.UseCase](ctx, cfg, scanStartProject.NewUseCase)
}

func (c *DependenciesContainer) ScanStopUseCase(ctx context.Context, cfg *domainConfig.Config) (*scanStop.UseCase, error) {
	return getUseCase[scanStop.UseCase](ctx, cfg, scanStop.NewUseCase)
}

func (c *DependenciesContainer) SetSettingsUseCase(ctx context.Context, cfg *domainConfig.Config) (*setSettings.UseCase, error) {
	return getUseCase[setSettings.UseCase](ctx, cfg, setSettings.NewUseCase)
}

func (c *DependenciesContainer) UpdateSourcesUseCase(ctx context.Context, cfg *domainConfig.Config) (*updateSources.UseCase, error) {
	return getUseCase[updateSources.UseCase](ctx, cfg, updateSources.NewUseCase)
}

func (c *DependenciesContainer) UpdateSourcesGitUseCase(ctx context.Context, cfg *domainConfig.Config) (*updateSourcesGit.UseCase, error) {
	return getUseCase[updateSourcesGit.UseCase](ctx, cfg, updateSourcesGit.NewUseCase)
}

func getConfigUseCase[T any](
	cfgAdapter port.Config,
	constructor func(port.Config, port.Cli) (*T, error),
) (*T, error) {
	cliAdapter := cli.NewCli()

	useCase, err := constructor(cfgAdapter, cliAdapter)
	if err != nil {
		return nil, err
	}

	return useCase, nil
}

func getUseCase[T any](
	ctx context.Context,
	cfg *domainConfig.Config,
	constructor func(port.Ai, port.Cli) (*T, error),
) (*T, error) {
	aiAdapter, err := ai.NewAdapter(ctx, cfg)
	if err != nil {
		return nil, err
	}

	cliAdapter := cli.NewCli()

	useCase, err := constructor(aiAdapter, cliAdapter)
	if err != nil {
		return nil, err
	}

	return useCase, nil
}
