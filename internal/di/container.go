package di

import (
	"github.com/POSIdev-community/aictl/internal/adapter/ai"
	"github.com/POSIdev-community/aictl/internal/adapter/cli"
	configAdapter "github.com/POSIdev-community/aictl/internal/adapter/config"
	configClear "github.com/POSIdev-community/aictl/internal/core/application/usecase/config/clear"
	"github.com/POSIdev-community/aictl/internal/core/application/usecase/config/set"
	"github.com/POSIdev-community/aictl/internal/core/application/usecase/config/show"
	"github.com/POSIdev-community/aictl/internal/core/application/usecase/config/unset"
	"github.com/POSIdev-community/aictl/internal/core/application/usecase/create/agenttoken"
	"github.com/POSIdev-community/aictl/internal/core/application/usecase/create/branch"
	"github.com/POSIdev-community/aictl/internal/core/application/usecase/create/project"
	"github.com/POSIdev-community/aictl/internal/core/application/usecase/delete/projects"
	getProjects "github.com/POSIdev-community/aictl/internal/core/application/usecase/get/projects"
	"github.com/POSIdev-community/aictl/internal/core/application/usecase/get/scan"
	"github.com/POSIdev-community/aictl/internal/core/application/usecase/get/scan/aiproj"
	"github.com/POSIdev-community/aictl/internal/core/application/usecase/get/scan/report"
	defaultreport "github.com/POSIdev-community/aictl/internal/core/application/usecase/get/scan/report/defaultreport"
	"github.com/POSIdev-community/aictl/internal/core/application/usecase/get/scan/state"
	"github.com/POSIdev-community/aictl/internal/core/application/usecase/scan/await"
	startBranch "github.com/POSIdev-community/aictl/internal/core/application/usecase/scan/start/branch"
	startProject "github.com/POSIdev-community/aictl/internal/core/application/usecase/scan/start/project"
	"github.com/POSIdev-community/aictl/internal/core/application/usecase/scan/stop"
	"github.com/POSIdev-community/aictl/internal/core/application/usecase/set/project/settings"
	"github.com/POSIdev-community/aictl/internal/core/application/usecase/update/sources"
	"github.com/POSIdev-community/aictl/internal/core/application/usecase/update/sources/git"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter"
	"github.com/POSIdev-community/aictl/internal/presenter/context"
	"github.com/POSIdev-community/aictl/internal/presenter/create"
	deletePresenter "github.com/POSIdev-community/aictl/internal/presenter/delete"
	"github.com/POSIdev-community/aictl/internal/presenter/get"
	scanPresenter "github.com/POSIdev-community/aictl/internal/presenter/scan"
	setPresenter "github.com/POSIdev-community/aictl/internal/presenter/set"
	"github.com/POSIdev-community/aictl/internal/presenter/update"
)

func InitializeCmd(cfg *config.Config) (*presenter.CmdRoot, error) {
	// Adapters
	cfgAdapter := configAdapter.NewContextAdapter()
	cliAdapter := cli.NewAdapter()
	aiAdapter, err := ai.NewAdapter(cfg)
	if err != nil {
		return nil, err
	}

	// Context commands
	cmdContext, err := buildContextCmd(cfgAdapter, cliAdapter, cfg)
	if err != nil {
		return nil, err
	}

	// Create commands
	cmdCreate, err := buildCreateCmd(aiAdapter, cliAdapter, cfg)
	if err != nil {
		return nil, err
	}

	// Delete commands
	cmdDelete, err := buildDeleteCmd(aiAdapter, cliAdapter, cfg)
	if err != nil {
		return nil, err
	}

	// Get commands
	cmdGet, err := buildGetCmd(aiAdapter, cliAdapter, cfg)
	if err != nil {
		return nil, err
	}

	// Scan commands
	cmdScan, err := buildScanCmd(aiAdapter, cliAdapter, cfg)
	if err != nil {
		return nil, err
	}

	// Set commands
	cmdSet, err := buildSetCmd(aiAdapter, cliAdapter, cfg)
	if err != nil {
		return nil, err
	}

	// Update commands
	cmdUpdate, err := buildUpdateCmd(aiAdapter, cliAdapter, cfg)
	if err != nil {
		return nil, err
	}

	return presenter.NewRootCmd(cmdContext, cmdCreate, cmdDelete, cmdGet, cmdScan, cmdSet, cmdUpdate), nil
}

func buildContextCmd(cfgAdapter *configAdapter.Adapter, cliAdapter *cli.Adapter, cfg *config.Config) (*context.CmdContext, error) {
	clearUC, err := configClear.NewUseCase(cfgAdapter, cliAdapter)
	if err != nil {
		return nil, err
	}
	cmdClear := context.NewConfigClearCommand(clearUC)

	setUC, err := set.NewUseCase(cfgAdapter, cfg)
	if err != nil {
		return nil, err
	}
	cmdSet := context.NewConfigSetCommand(cfg, setUC)

	showUC, err := show.NewUseCase(cfgAdapter, cliAdapter, cfg)
	if err != nil {
		return nil, err
	}
	cmdShow := context.NewConfigShowCommand(showUC)

	unsetUC, err := unset.NewUseCase(cfgAdapter, cfg)
	if err != nil {
		return nil, err
	}
	cmdUnset := context.NewConfigUnsetCommand(unsetUC)

	return context.NewContextCmd(cmdClear, cmdSet, cmdShow, cmdUnset), nil
}

func buildCreateCmd(aiAdapter *ai.Adapter, cliAdapter *cli.Adapter, cfg *config.Config) (*create.CmdCreate, error) {
	branchUC, err := branch.NewUseCase(aiAdapter, cliAdapter)
	if err != nil {
		return nil, err
	}
	cmdBranch := create.NewCreateBranchCmd(cfg, branchUC)

	projectUC, err := project.NewUseCase(aiAdapter, cliAdapter)
	if err != nil {
		return nil, err
	}
	cmdProject := create.NewCreateProjectCmd(projectUC)

	agentTokenUC, err := agenttoken.NewUseCase(aiAdapter, cliAdapter, cfg)
	if err != nil {
		return nil, err
	}
	cmdAgentToken := create.NewCreateAgentTokenCmd(cfg, agentTokenUC)

	return create.NewCreateCmd(cfg, cmdBranch, cmdProject, cmdAgentToken), nil
}

func buildDeleteCmd(aiAdapter *ai.Adapter, cliAdapter *cli.Adapter, cfg *config.Config) (*deletePresenter.CmdDelete, error) {
	projectsUC, err := projects.NewUseCase(aiAdapter, cliAdapter)
	if err != nil {
		return nil, err
	}
	cmdProjects := deletePresenter.NewDeleteProjectsCommand(projectsUC)

	return deletePresenter.NewDeleteCmd(cfg, cmdProjects), nil
}

func buildGetCmd(aiAdapter *ai.Adapter, cliAdapter *cli.Adapter, cfg *config.Config) (*get.CmdGet, error) {
	// Projects
	projectsUC, err := getProjects.NewUseCase(aiAdapter, cliAdapter)
	if err != nil {
		return nil, err
	}
	cmdProjects := get.NewGetProjectsCmd(projectsUC)

	// Scan
	scanUC, err := scan.NewUseCase(aiAdapter, cliAdapter, cfg)
	if err != nil {
		return nil, err
	}

	aiprojUC, err := aiproj.NewUseCase(aiAdapter, cliAdapter, cfg)
	if err != nil {
		return nil, err
	}
	cmdAiproj := get.NewGetScanAiprojCmd(aiprojUC)

	cmdLogs := get.NewGetScanLogsCmd()

	defaultReportUC, err := defaultreport.NewUseCase(aiAdapter, cliAdapter, cfg)
	if err != nil {
		return nil, err
	}

	cmdReportAutocheck := get.NewGetScanReportAutocheckCmd(defaultReportUC)
	cmdReportGitlab := get.NewGetScanReportGitlabCmd(defaultReportUC)
	cmdReportJson := get.NewGetScanReportJsonCmd(defaultReportUC)
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

	persistentPreRunEGetCmd := get.NewPersistentPreRunEGetCmd(cfg)
	persistentPreRunEGetScanCmd := get.NewPersistentPreRunEGetScanCmd(cfg, persistentPreRunEGetCmd)
	persistentPreRunEGetScanReportCmd := get.NewPersistentPreRunEGetScanReportCmd(persistentPreRunEGetScanCmd)

	customReportUC, err := report.NewUseCase(aiAdapter, cliAdapter, cfg)
	if err != nil {
		return nil, err
	}

	cmdReport := get.NewGetScanReportCmd(customReportUC, persistentPreRunEGetScanReportCmd, cmdReportAutocheck, cmdReportGitlab,
		cmdReportJson, cmdReportMarkdown, cmdReportNist, cmdReportOud4, cmdReportOwasp, cmdReportOwaspm,
		cmdReportPcidss, cmdReportPlain, cmdReportSans, cmdReportSarif, cmdReportXml)

	cmdSbom := get.NewGetScanSbomCmd()

	stateUC, err := state.NewUseCase(aiAdapter, cliAdapter, cfg)
	if err != nil {
		return nil, err
	}
	cmdState := get.NewGetScanStateCmd(stateUC)

	cmdScan := get.NewGetScanCmd(persistentPreRunEGetScanCmd, scanUC, cmdAiproj, cmdLogs, cmdReport, cmdSbom, cmdState)

	return get.NewGetCmd(persistentPreRunEGetCmd, cmdProjects, cmdScan), nil
}

func buildScanCmd(aiAdapter *ai.Adapter, cliAdapter *cli.Adapter, cfg *config.Config) (*scanPresenter.CmdScan, error) {
	awaitUC, err := await.NewUseCase(aiAdapter, cliAdapter, cfg)
	if err != nil {
		return nil, err
	}
	cmdAwait := scanPresenter.NewScanAwaitCmd(cfg, awaitUC)

	branchUC, err := startBranch.NewUseCase(aiAdapter, cliAdapter, cfg)
	if err != nil {
		return nil, err
	}
	cmdStartBranch := scanPresenter.NewScanStartBranchCmd(cfg, branchUC)

	projectUC, err := startProject.NewUseCase(aiAdapter, cliAdapter, cfg)
	if err != nil {
		return nil, err
	}
	cmdStartProject := scanPresenter.NewScanStartProjectCmd(cfg, projectUC)

	persistentPreRunEScanCmd := scanPresenter.NewPersistentPreRunEScanCmd(cfg)
	persistentPreRunEScanStartCmd := scanPresenter.NewPersistentPreRunEScanStartCmd(persistentPreRunEScanCmd)

	cmdStart := scanPresenter.NewScanStartCmd(persistentPreRunEScanStartCmd, cmdStartBranch, cmdStartProject)

	stopUC, err := stop.NewUseCase(aiAdapter, cliAdapter)
	if err != nil {
		return nil, err
	}
	cmdStop := scanPresenter.NewScanStopCmd(stopUC)

	return scanPresenter.NewScanCmd(persistentPreRunEScanCmd, cmdAwait, cmdStart, cmdStop), nil
}

func buildSetCmd(aiAdapter *ai.Adapter, cliAdapter *cli.Adapter, cfg *config.Config) (*setPresenter.CmdSet, error) {
	settingsUC, err := settings.NewUseCase(aiAdapter, cliAdapter, cfg)
	if err != nil {
		return nil, err
	}
	cmdSettings := setPresenter.NewSetProjectSettingsCmd(settingsUC)
	cmdProject := setPresenter.NewSetProjectCmd(cfg, cmdSettings)

	return setPresenter.NewSetCmd(cfg, cmdProject), nil
}

func buildUpdateCmd(aiAdapter *ai.Adapter, cliAdapter *cli.Adapter, cfg *config.Config) (*update.CmdUpdate, error) {
	sourcesUC, err := sources.NewUseCase(aiAdapter, cliAdapter, cfg)
	if err != nil {
		return nil, err
	}

	gitUC, err := git.NewUseCase(aiAdapter, cliAdapter)
	if err != nil {
		return nil, err
	}
	cmdGit := update.NewUpdateSourcesGitCmd(gitUC)

	cmdSources := update.NewUpdateSourcesCmd(cfg, sourcesUC, cmdGit)

	return update.NewUpdateCmd(cfg, cmdSources), nil
}
