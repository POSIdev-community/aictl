package di

import (
	getBranch "github.com/POSIdev-community/aictl/internal/core/usecase/get/branch"
	getBranches "github.com/POSIdev-community/aictl/internal/core/usecase/get/branches"
	getHealthchech "github.com/POSIdev-community/aictl/internal/core/usecase/get/healthcheck"
	projectAiproj "github.com/POSIdev-community/aictl/internal/core/usecase/get/project/aiproj"
	getProjectSettings "github.com/POSIdev-community/aictl/internal/core/usecase/get/project/settings"
	getProjects "github.com/POSIdev-community/aictl/internal/core/usecase/get/projects"
	"github.com/POSIdev-community/aictl/internal/core/usecase/get/scan"
	scanAiproj "github.com/POSIdev-community/aictl/internal/core/usecase/get/scan/aiproj"
	"github.com/POSIdev-community/aictl/internal/core/usecase/get/scan/logs"
	"github.com/POSIdev-community/aictl/internal/core/usecase/get/scan/report"
	defaultreport "github.com/POSIdev-community/aictl/internal/core/usecase/get/scan/report/defaultreport"
	"github.com/POSIdev-community/aictl/internal/core/usecase/get/scan/sbom"
	"github.com/POSIdev-community/aictl/internal/core/usecase/get/scan/state"
	"github.com/POSIdev-community/aictl/internal/core/usecase/get/scan/statistic"
	getScanAgents "github.com/POSIdev-community/aictl/internal/core/usecase/get/scanagents"
	getScans "github.com/POSIdev-community/aictl/internal/core/usecase/get/scans"
	getVersion "github.com/POSIdev-community/aictl/internal/core/usecase/get/version"
	"github.com/POSIdev-community/aictl/internal/presenter/get"
)

func buildGetCmd(a *adapters) (*get.CmdGet, error) {
	healthcheckUC, err := getHealthchech.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return nil, err
	}
	cmdHealthcheck := get.NewGetHealthcheckCmd(healthcheckUC)

	projectsUC, err := getProjects.NewUseCase(a.ai, a.cli)
	if err != nil {
		return nil, err
	}
	cmdProjects := get.NewGetProjectsCmd(projectsUC)

	cmdProject, err := buildGetProjectCmd(a)
	if err != nil {
		return nil, err
	}

	branchesUC, err := getBranches.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return nil, err
	}
	cmdBranches := get.NewGetBranchesCmd(a.cfg, branchesUC)

	cmdBranch, err := buildGetBranchCmd(a)
	if err != nil {
		return nil, err
	}

	scansUC, err := getScans.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return nil, err
	}
	cmdScans := get.NewGetScansCmd(a.cfg, scansUC)

	cmdScan, err := buildGetScanCmd(a)
	if err != nil {
		return nil, err
	}

	versionUC, err := getVersion.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return nil, err
	}
	cmdVersion := get.NewGetVersionCmd(versionUC)

	scanAgentsUC, err := getScanAgents.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return nil, err
	}
	cmdAgents := get.NewGetAgentsCmd(scanAgentsUC)

	persistentPreRunEGetCmd := get.NewPersistentPreRunEGetCmd(a.cfg)

	return get.NewGetCmd(persistentPreRunEGetCmd, cmdHealthcheck, cmdProjects, cmdProject, cmdBranches, cmdBranch, cmdScans, cmdScan, cmdAgents, cmdVersion), nil
}

func buildGetBranchCmd(a *adapters) (get.CmdGetBranch, error) {
	branchUC, err := getBranch.NewUseCase(a.ai, a.cli)
	if err != nil {
		return get.CmdGetBranch{}, err
	}

	persistentPreRunEGetCmd := get.NewPersistentPreRunEGetCmd(a.cfg)
	persistentPreRunEGetBranchCmd := get.NewPersistentPreRunEGetBranchCmd(persistentPreRunEGetCmd)

	return get.NewGetBranchCmd(persistentPreRunEGetBranchCmd, branchUC), nil
}

func buildGetProjectCmd(a *adapters) (get.CmdGetProject, error) {
	projectAiprojUC, err := projectAiproj.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return get.CmdGetProject{}, err
	}
	cmdAiproj := get.NewGetProjectAiprojCmd(projectAiprojUC)

	projectSettingsUC, err := getProjectSettings.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return get.CmdGetProject{}, err
	}
	cmdSettings := get.NewGetProjectSettingsCmd(projectSettingsUC)

	persistentPreRunEGetCmd := get.NewPersistentPreRunEGetCmd(a.cfg)
	persistentPreRunEGetProjectCmd := get.NewPersistentPreRunEGetProjectCmd(a.cfg, persistentPreRunEGetCmd)

	return get.NewGetProjectCmd(persistentPreRunEGetProjectCmd, cmdAiproj, cmdSettings), nil
}

func buildGetScanCmd(a *adapters) (get.CmdGetScan, error) {
	scanUC, err := scan.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return get.CmdGetScan{}, err
	}

	aiprojUC, err := scanAiproj.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return get.CmdGetScan{}, err
	}
	cmdAiproj := get.NewGetScanAiprojCmd(aiprojUC)

	logsUC, err := logs.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return get.CmdGetScan{}, err
	}
	cmdLogs := get.NewGetScanLogsCmd(logsUC)

	cmdReport, err := buildGetScanReportCmd(a)
	if err != nil {
		return get.CmdGetScan{}, err
	}

	getSbomUC, err := sbom.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return get.CmdGetScan{}, err
	}
	cmdSbom := get.NewGetScanSbomCmd(getSbomUC)

	stateUC, err := state.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return get.CmdGetScan{}, err
	}
	cmdState := get.NewGetScanStateCmd(stateUC)

	statisticUC, err := statistic.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return get.CmdGetScan{}, err
	}
	cmdStatistic := get.NewGetScanStatisticCmd(statisticUC)

	persistentPreRunEGetCmd := get.NewPersistentPreRunEGetCmd(a.cfg)
	persistentPreRunEGetScanCmd := get.NewPersistentPreRunEGetScanCmd(a.cfg, persistentPreRunEGetCmd)

	return get.NewGetScanCmd(persistentPreRunEGetScanCmd, scanUC, cmdAiproj, cmdLogs, cmdReport, cmdSbom, cmdState, cmdStatistic), nil
}

func buildGetScanReportCmd(a *adapters) (get.CmdGetScanReport, error) {
	defaultReportUC, err := defaultreport.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return get.CmdGetScanReport{}, err
	}

	cmdReportAutocheck := get.NewGetScanReportAutocheckCmd(defaultReportUC)
	cmdReportGitlab := get.NewGetScanReportGitlabCmd(defaultReportUC)
	cmdReportJson := get.NewGetScanReportJsonCmd(defaultReportUC)
	cmdReportJsonV2 := get.NewGetScanReportJsonV2Cmd(defaultReportUC)
	cmdReportMarkdown := get.NewGetScanReportMarkdownCmd(defaultReportUC)
	cmdReportNist := get.NewGetScanReportNistCmd(defaultReportUC)
	cmdReportOud4 := get.NewGetScanReportOud4Cmd(defaultReportUC)
	cmdReportOwasp := get.NewGetScanReportOwaspCmd(defaultReportUC)
	cmdReportOwaspm := get.NewGetScanReportOwaspmCmd(defaultReportUC)
	cmdReportPcidss := get.NewGetScanReportPcidssCmd(defaultReportUC)
	cmdReportPlain := get.NewGetScanReportPlainCmd(defaultReportUC)
	cmdReportSans := get.NewGetScanReportSansCmd(defaultReportUC)
	cmdReportSarif := get.NewGetScanReportSarifCmd(defaultReportUC)
	cmdReportXml := get.NewGetScanReportXmlCmd(defaultReportUC)

	persistentPreRunEGetCmd := get.NewPersistentPreRunEGetCmd(a.cfg)
	persistentPreRunEGetScanCmd := get.NewPersistentPreRunEGetScanCmd(a.cfg, persistentPreRunEGetCmd)
	persistentPreRunEGetScanReportCmd := get.NewPersistentPreRunEGetScanReportCmd(persistentPreRunEGetScanCmd)

	customReportUC, err := report.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return get.CmdGetScanReport{}, err
	}

	return get.NewGetScanReportCmd(customReportUC, persistentPreRunEGetScanReportCmd, cmdReportAutocheck, cmdReportGitlab,
		cmdReportJson, cmdReportJsonV2, cmdReportMarkdown, cmdReportNist, cmdReportOud4, cmdReportOwasp, cmdReportOwaspm,
		cmdReportPcidss, cmdReportPlain, cmdReportSans, cmdReportSarif, cmdReportXml), nil
}
