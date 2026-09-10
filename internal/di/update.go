package di

import (
	updateProjectLanguages "github.com/POSIdev-community/aictl/internal/core/usecase/update/project/languages"
	updateProjectSettings "github.com/POSIdev-community/aictl/internal/core/usecase/update/project/settings"
	"github.com/POSIdev-community/aictl/internal/core/usecase/update/sbom"
	"github.com/POSIdev-community/aictl/internal/core/usecase/update/sources"
	"github.com/POSIdev-community/aictl/internal/presenter/update"
)

func buildUpdateCmd(a *adapters) (*update.CmdUpdate, error) {
	sourcesUC, err := sources.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return nil, err
	}

	cmdSources := update.NewUpdateSourcesCmd(a.cfg, sourcesUC)

	sbomUC, err := sbom.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return nil, err
	}
	cmdSbom := update.NewUpdateSbomCmd(a.cfg, sbomUC)

	cmdProject, err := buildUpdateProjectCmd(a)
	if err != nil {
		return nil, err
	}

	return update.NewUpdateCmd(a.cfg, cmdSources, cmdSbom, cmdProject), nil
}

func buildUpdateProjectCmd(a *adapters) (update.CmdUpdateProject, error) {
	settingsUC, err := updateProjectSettings.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return update.CmdUpdateProject{}, err
	}
	cmdSettings := update.NewUpdateProjectSettingsCmd(settingsUC)

	languagesUC, err := updateProjectLanguages.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return update.CmdUpdateProject{}, err
	}
	cmdLanguages := update.NewUpdateProjectLanguagesCmd(languagesUC)

	persistentPreRunEUpdateCmd := update.NewPersistentPreRunEUpdateCmd(a.cfg)
	persistentPreRunEUpdateProjectCmd := update.NewPersistentPreRunEUpdateProjectCmd(a.cfg, persistentPreRunEUpdateCmd)

	return update.NewUpdateProjectCmd(persistentPreRunEUpdateProjectCmd, cmdSettings, cmdLanguages), nil
}
