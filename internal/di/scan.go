package di

import (
	"github.com/POSIdev-community/aictl/internal/core/usecase/scan/await"
	startBranch "github.com/POSIdev-community/aictl/internal/core/usecase/scan/start/branch"
	startProject "github.com/POSIdev-community/aictl/internal/core/usecase/scan/start/project"
	"github.com/POSIdev-community/aictl/internal/core/usecase/scan/stop"
	scanPresenter "github.com/POSIdev-community/aictl/internal/presenter/scan"
)

func buildScanCmd(a *adapters) (*scanPresenter.CmdScan, error) {
	awaitUC, err := await.NewUseCase(a.ai, a.cli, a.cfg, await.DefaultPollInterval)
	if err != nil {
		return nil, err
	}
	cmdAwait := scanPresenter.NewScanAwaitCmd(a.cfg, awaitUC)

	branchUC, err := startBranch.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return nil, err
	}
	cmdStartBranch := scanPresenter.NewScanStartBranchCmd(a.cfg, branchUC)

	projectUC, err := startProject.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return nil, err
	}
	cmdStartProject := scanPresenter.NewScanStartProjectCmd(a.cfg, projectUC)

	persistentPreRunEScanCmd := scanPresenter.NewPersistentPreRunEScanCmd(a.cfg)
	persistentPreRunEScanStartCmd := scanPresenter.NewPersistentPreRunEScanStartCmd(persistentPreRunEScanCmd)

	cmdStart := scanPresenter.NewScanStartCmd(persistentPreRunEScanStartCmd, cmdStartBranch, cmdStartProject)

	stopUC, err := stop.NewUseCase(a.ai, a.cli)
	if err != nil {
		return nil, err
	}
	cmdStop := scanPresenter.NewScanStopCmd(stopUC)

	return scanPresenter.NewScanCmd(persistentPreRunEScanCmd, cmdAwait, cmdStart, cmdStop), nil
}
