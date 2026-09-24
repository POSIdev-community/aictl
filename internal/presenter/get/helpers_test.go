package get

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/regexfilter"
	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/POSIdev-community/aictl/internal/core/domain/scafeeds"
	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

func resetGetPackageFlags() {
	projectIdFlag = ""
	scanId = uuid.Nil
	branchId = uuid.Nil
	customReportName = ""
	outPath = ""
	forceRewriteOutPath = false
	includeComments = false
	includeDFD = false
	includeGlossary = false
	l10n = "en"
	resetReportFilterFlags()
}

type noopHealthcheckUC struct{}

func (noopHealthcheckUC) Execute(context.Context) error { return nil }

type noopVersionUC struct{}

func (noopVersionUC) Execute(context.Context) error { return nil }

type noopAgentsUC struct{}

func (noopAgentsUC) Execute(context.Context, bool) error { return nil }

type noopProjectsUC struct{}

func (noopProjectsUC) Execute(context.Context, regexfilter.RegexFilter, bool) error { return nil }

type noopBranchesUC struct{}

func (noopBranchesUC) Execute(context.Context, regexfilter.RegexFilter, bool) error { return nil }

type noopBranchUC struct{}

func (noopBranchUC) Execute(context.Context, uuid.UUID) error { return nil }

type noopScansUC struct{}

func (noopScansUC) Execute(context.Context, regexfilter.RegexFilter, bool, bool) error {
	return nil
}

type noopScansSbomProjectUC struct{}

func (noopScansSbomProjectUC) Execute(context.Context, regexfilter.RegexFilter, bool, bool) error {
	return nil
}

type noopScanUC struct{}

func (noopScanUC) Execute(context.Context, uuid.UUID) error { return nil }

type noopQueueUC struct{}

func (noopQueueUC) Execute(context.Context, *uuid.UUID) error { return nil }

type noopScanningUC struct{}

func (noopScanningUC) Execute(context.Context, *uuid.UUID) error { return nil }

type noopReportTemplatesUC struct{}

func (noopReportTemplatesUC) Execute(context.Context, regexfilter.RegexFilter, bool, string) error {
	return nil
}

type noopScaFeedsListUC struct{}

func (noopScaFeedsListUC) Execute(context.Context, []scafeeds.Status) error { return nil }

type noopScaFeedsDownloadUC struct{}

func (noopScaFeedsDownloadUC) Execute(context.Context, string, string) error { return nil }

type noopProjectAiprojUC struct{}

func (noopProjectAiprojUC) Execute(context.Context, string) error { return nil }

type noopProjectSettingsUC struct{}

func (noopProjectSettingsUC) Execute(context.Context, bool) error { return nil }

type noopProjectPoliciesUC struct{}

func (noopProjectPoliciesUC) Execute(context.Context) error { return nil }

type noopProjectPolicyCheckUC struct{}

func (noopProjectPolicyCheckUC) Execute(context.Context) error { return nil }

type noopProjectExclusionsUC struct{}

func (noopProjectExclusionsUC) Execute(context.Context) error { return nil }

type noopScanAiprojUC struct{}

func (noopScanAiprojUC) Execute(context.Context, uuid.UUID, string) error { return nil }

type noopScanLogsUC struct{}

func (noopScanLogsUC) Execute(context.Context, uuid.UUID, string) error { return nil }

type noopScanErrorsUC struct{}

func (noopScanErrorsUC) Execute(context.Context, uuid.UUID) error { return nil }

type noopScanSbomUC struct{}

func (noopScanSbomUC) Execute(context.Context, uuid.UUID, string) error { return nil }

type noopScanStateUC struct{}

func (noopScanStateUC) Execute(context.Context, uuid.UUID, bool) error { return nil }

type noopScanStatisticUC struct{}

func (noopScanStatisticUC) Execute(context.Context, uuid.UUID, string, bool, bool) error { return nil }

type noopDefaultReportUC struct{}

func (noopDefaultReportUC) Execute(
	context.Context, uuid.UUID, report.ReportType, string, bool, bool, bool, string, report.Filters,
) error {
	return nil
}

type noopCustomReportUC struct{}

func (noopCustomReportUC) Execute(
	context.Context, uuid.UUID, string, string, bool, bool, bool, string, report.Filters,
) error {
	return nil
}

type getUCs struct {
	healthcheck        UseCaseGetHealthcheck
	projects           UseCaseGetProjects
	branches           UseCaseGetBranches
	branch             UseCaseGetBranch
	scans              UseCaseGetScans
	scansSbomProject   UseCaseGetScansSbomProject
	scan               UseCaseGetScan
	agents             UseCaseGetAgents
	version            UseCaseGetVersion
	queue              UseCaseGetQueue
	scanning           UseCaseGetScanning
	reportTemplates    UseCaseGetReportTemplates
	projectAiproj      UseCaseGetProjectAiproj
	projectSettings    UseCaseGetProjectSettings
	projectPolicies    UseCaseGetProjectPolicies
	projectPolicyCheck UseCaseGetProjectPolicyCheck
	projectExcl        UseCaseGetProjectExclusions
	scanAiproj         UseCaseGetScanAiproj
	scanLogs           UseCaseGetScanLogs
	scanErrors         UseCaseGetScanErrors
	scanSbom           UseCaseGetScanSbom
	scanState          UseCaseGetScanState
	scanStatistic      UseCaseGetScanStatistic
	customReport       UseCaseGetScanReport
	defaultReport      noopDefaultReportUC
}

func defaultGetUCs() getUCs {
	return getUCs{
		healthcheck:        noopHealthcheckUC{},
		projects:           noopProjectsUC{},
		branches:           noopBranchesUC{},
		branch:             noopBranchUC{},
		scans:              noopScansUC{},
		scansSbomProject:   noopScansSbomProjectUC{},
		scan:               noopScanUC{},
		agents:             noopAgentsUC{},
		version:            noopVersionUC{},
		queue:              noopQueueUC{},
		scanning:           noopScanningUC{},
		reportTemplates:    noopReportTemplatesUC{},
		projectAiproj:      noopProjectAiprojUC{},
		projectSettings:    noopProjectSettingsUC{},
		projectPolicies:    noopProjectPoliciesUC{},
		projectPolicyCheck: noopProjectPolicyCheckUC{},
		projectExcl:        noopProjectExclusionsUC{},
		scanAiproj:         noopScanAiprojUC{},
		scanLogs:           noopScanLogsUC{},
		scanErrors:         noopScanErrorsUC{},
		scanSbom:           noopScanSbomUC{},
		scanState:          noopScanStateUC{},
		scanStatistic:      noopScanStatisticUC{},
		customReport:       noopCustomReportUC{},
		defaultReport:      noopDefaultReportUC{},
	}
}

func buildGetRoot(t *testing.T, cfg *config.Config, ucs getUCs) *CmdGet {
	t.Helper()
	require.NotNil(t, cfg)

	dr := ucs.defaultReport
	reportPreRun := NewPersistentPreRunEGetScanReportCmd(NewPersistentPreRunEGetScanCmd(cfg, NewPersistentPreRunEGetCmd(cfg)))
	cmdReportWithFilters := NewGetScanReportWithFiltersCmd(
		ucs.customReport,
		dr,
		NewPersistentPreRunEGetScanReportWithFiltersCmd(reportPreRun),
	)
	cmdReport := NewGetScanReportCmd(
		ucs.customReport,
		reportPreRun,
		cmdReportWithFilters,
		NewGetScanReportAutocheckCmd(dr),
		NewGetScanReportGitlabCmd(dr),
		NewGetScanReportJsonCmd(dr),
		NewGetScanReportJsonV2Cmd(dr),
		NewGetScanReportMarkdownCmd(dr),
		NewGetScanReportNistCmd(dr),
		NewGetScanReportOud4Cmd(dr),
		NewGetScanReportOwaspCmd(dr),
		NewGetScanReportOwaspmCmd(dr),
		NewGetScanReportPcidssCmd(dr),
		NewGetScanReportPlainCmd(dr),
		NewGetScanReportSansCmd(dr),
		NewGetScanReportSarifCmd(dr),
		NewGetScanReportXmlCmd(dr),
	)

	cmdScan := NewGetScanCmd(
		NewPersistentPreRunEGetScanCmd(cfg, NewPersistentPreRunEGetCmd(cfg)),
		ucs.scan,
		NewGetScanAiprojCmd(ucs.scanAiproj),
		NewGetScanLogsCmd(ucs.scanLogs),
		NewGetScanErrorsCmd(ucs.scanErrors),
		cmdReport,
		NewGetScanSbomCmd(ucs.scanSbom),
		NewGetScanStateCmd(ucs.scanState),
		NewGetScanStatisticCmd(ucs.scanStatistic),
	)

	cmdProject := NewGetProjectCmd(
		NewPersistentPreRunEGetProjectCmd(cfg, NewPersistentPreRunEGetCmd(cfg)),
		NewGetProjectAiprojCmd(ucs.projectAiproj),
		NewGetProjectSettingsCmd(ucs.projectSettings),
		NewGetProjectPoliciesCmd(ucs.projectPolicies),
		NewGetProjectPolicyCheckCmd(ucs.projectPolicyCheck),
		NewGetProjectExclusionsCmd(ucs.projectExcl),
	)

	return NewGetCmd(
		NewPersistentPreRunEGetCmd(cfg),
		NewGetHealthcheckCmd(ucs.healthcheck),
		NewGetProjectsCmd(ucs.projects),
		cmdProject,
		NewGetBranchesCmd(cfg, ucs.branches),
		NewGetBranchCmd(NewPersistentPreRunEGetBranchCmd(NewPersistentPreRunEGetCmd(cfg)), ucs.branch),
		NewGetScansCmd(cfg, ucs.scans, NewGetScansSbomProjectCmd(cfg, ucs.scansSbomProject)),
		cmdScan,
		NewGetAgentsCmd(ucs.agents),
		NewGetVersionCmd(ucs.version),
		NewGetQueueCmd(ucs.queue),
		NewGetScanningCmd(ucs.scanning),
		NewGetReportTemplatesCmd(ucs.reportTemplates),
		NewGetScaFeedsCmd(noopScaFeedsListUC{}, noopScaFeedsDownloadUC{}),
	)
}

func mustGetCfg(t *testing.T) *config.Config {
	t.Helper()
	return cmdtest.MustCfg(t)
}
