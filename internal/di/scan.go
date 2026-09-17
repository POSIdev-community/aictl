package di

import (
	"github.com/POSIdev-community/aictl/internal/core/usecase/scan/await"
	"github.com/POSIdev-community/aictl/internal/core/usecase/scan/checkpolicies"
	startBranch "github.com/POSIdev-community/aictl/internal/core/usecase/scan/start/branch"
	startProject "github.com/POSIdev-community/aictl/internal/core/usecase/scan/start/project"
	startSbom "github.com/POSIdev-community/aictl/internal/core/usecase/scan/start/sbom"
	"github.com/POSIdev-community/aictl/internal/core/usecase/scan/stop"
	scanPresenter "github.com/POSIdev-community/aictl/internal/presenter/scan"
)

func buildScanCmd(a *adapters) (*scanPresenter.CmdScan, error) {
	awaitUC, err := await.NewUseCase(a.ai, a.cli, a.cfg, await.DefaultPollInterval)
	if err != nil {
		return nil, err
	}
	cmdAwait := scanPresenter.NewScanAwaitCmd(a.cfg, awaitUC)

	checkPoliciesUC, err := checkpolicies.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return nil, err
	}
	cmdCheckPolicies := scanPresenter.NewScanCheckPoliciesCmd(a.cfg, checkPoliciesUC)

	branchUC, err := startBranch.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return nil, err
	}
	cmdStartBranch := scanPresenter.NewScanStartBranchCmd(a.cfg, branchUC)
	cmdBranch := scanPresenter.NewScanBranchCmd(a.cfg, branchUC)

	projectUC, err := startProject.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return nil, err
	}
	cmdStartProject := scanPresenter.NewScanStartProjectCmd(a.cfg, projectUC)
	cmdProject := scanPresenter.NewScanProjectCmd(a.cfg, projectUC)

	sbomUC, err := startSbom.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return nil, err
	}
	cmdSbom := scanPresenter.NewScanSbomCmd(a.cfg, sbomUC)

	persistentPreRunEScanCmd := scanPresenter.NewPersistentPreRunEScanCmd(a.cfg)
	persistentPreRunEScanStartCmd := scanPresenter.NewPersistentPreRunEScanStartCmd(persistentPreRunEScanCmd)

	cmdStart := scanPresenter.NewScanStartCmd(persistentPreRunEScanStartCmd, cmdStartBranch, cmdStartProject)

	stopUC, err := stop.NewUseCase(a.ai)
	if err != nil {
		return nil, err
	}
	cmdStop := scanPresenter.NewScanStopCmd(stopUC)

	return scanPresenter.NewScanCmd(
		persistentPreRunEScanCmd,
		cmdAwait,
		cmdCheckPolicies,
		cmdBranch,
		cmdProject,
		cmdSbom,
		cmdStart,
		cmdStop,
	), nil
}
