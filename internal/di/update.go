package di

import (
	updateProjectSettings "github.com/POSIdev-community/aictl/internal/core/usecase/update/project/settings"
	"github.com/POSIdev-community/aictl/internal/core/usecase/update/sources"
	"github.com/POSIdev-community/aictl/internal/presenter/update"
)

func buildUpdateCmd(a *adapters) (*update.CmdUpdate, error) {
	sourcesUC, err := sources.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return nil, err
	}

	cmdSources := update.NewUpdateSourcesCmd(a.cfg, sourcesUC)

	cmdProject, err := buildUpdateProjectCmd(a)
	if err != nil {
		return nil, err
	}

	return update.NewUpdateCmd(a.cfg, cmdSources, cmdProject), nil
}

func buildUpdateProjectCmd(a *adapters) (update.CmdUpdateProject, error) {
	settingsUC, err := updateProjectSettings.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return update.CmdUpdateProject{}, err
	}
	cmdSettings := update.NewUpdateProjectSettingsCmd(settingsUC)

	persistentPreRunEUpdateCmd := update.NewPersistentPreRunEUpdateCmd(a.cfg)
	persistentPreRunEUpdateProjectCmd := update.NewPersistentPreRunEUpdateProjectCmd(a.cfg, persistentPreRunEUpdateCmd)

	return update.NewUpdateProjectCmd(persistentPreRunEUpdateProjectCmd, cmdSettings), nil
}
