package di

import (
	"github.com/POSIdev-community/aictl/internal/core/usecase/set/project/exclusions"
	"github.com/POSIdev-community/aictl/internal/core/usecase/set/project/policies"
	setPolicyCheck "github.com/POSIdev-community/aictl/internal/core/usecase/set/project/policycheck"
	"github.com/POSIdev-community/aictl/internal/core/usecase/set/project/settings"
	setPresenter "github.com/POSIdev-community/aictl/internal/presenter/set"
)

func buildSetCmd(a *adapters) (*setPresenter.CmdSet, error) {
	settingsUC, err := settings.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return nil, err
	}

	policiesUC, err := policies.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return nil, err
	}

	policyCheckUC, err := setPolicyCheck.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return nil, err
	}

	exclusionsUC, err := exclusions.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return nil, err
	}

	persistentPreRunESetCmd := setPresenter.NewPersistentPreRunESetCmd(a.cfg)
	persistentPreRunESetProjectCmd := setPresenter.NewPersistentPreRunESetProjectCmd(a.cfg, persistentPreRunESetCmd)

	cmdSettings := setPresenter.NewSetProjectSettingsCmd(settingsUC)
	cmdPolicies := setPresenter.NewSetProjectPoliciesCmd(policiesUC)
	cmdPolicyCheck := setPresenter.NewSetProjectPolicyCheckCmd(policyCheckUC)
	cmdExclusions := setPresenter.NewSetProjectExclusionsCmd(exclusionsUC)
	cmdProject := setPresenter.NewSetProjectCmd(persistentPreRunESetProjectCmd, cmdSettings, cmdPolicies, cmdPolicyCheck, cmdExclusions)

	return setPresenter.NewSetCmd(persistentPreRunESetCmd, cmdProject), nil
}
