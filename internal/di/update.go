package di

import (
	updateProjectLanguages "github.com/POSIdev-community/aictl/internal/core/usecase/update/project/languages"
	updateProjectSettings "github.com/POSIdev-community/aictl/internal/core/usecase/update/project/settings"
	"github.com/POSIdev-community/aictl/internal/core/usecase/update/sbom"
	"github.com/POSIdev-community/aictl/internal/core/usecase/update/scafeeds"
	scaFeedsRollback "github.com/POSIdev-community/aictl/internal/core/usecase/update/scafeeds/rollback"
	"github.com/POSIdev-community/aictl/internal/core/usecase/update/sources"
	"github.com/POSIdev-community/aictl/internal/presenter/update"
)

func buildUpdateCmd(a *adapters) (*update.CmdUpdate, error) {
	sourcesUC, err := sources.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return nil, err
	}

	languagesUC, err := updateProjectLanguages.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return nil, err
	}

	cmdSources := update.NewUpdateSourcesCmd(a.cfg, sourcesUC, languagesUC)

	sbomUC, err := sbom.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return nil, err
	}
	cmdSbom := update.NewUpdateSbomCmd(a.cfg, sbomUC)

	scaFeedsUC, err := scafeeds.NewUseCase(a.ai, a.cli)
	if err != nil {
		return nil, err
	}
	rollbackUC, err := scaFeedsRollback.NewUseCase(a.ai, a.cli)
	if err != nil {
		return nil, err
	}
	cmdScaFeeds := update.NewUpdateScaFeedsCmd(scaFeedsUC, rollbackUC)

	cmdProject, err := buildUpdateProjectCmd(a, languagesUC)
	if err != nil {
		return nil, err
	}

	return update.NewUpdateCmd(a.cfg, cmdSources, cmdSbom, cmdScaFeeds, cmdProject), nil
}

func buildUpdateProjectCmd(a *adapters, languagesUC *updateProjectLanguages.UseCase) (update.CmdUpdateProject, error) {
	settingsUC, err := updateProjectSettings.NewUseCase(a.ai, a.cli, a.cfg)
	if err != nil {
		return update.CmdUpdateProject{}, err
	}
	cmdSettings := update.NewUpdateProjectSettingsCmd(settingsUC)

	cmdLanguages := update.NewUpdateProjectLanguagesCmd(languagesUC)

	persistentPreRunEUpdateCmd := update.NewPersistentPreRunEUpdateCmd(a.cfg)
	persistentPreRunEUpdateProjectCmd := update.NewPersistentPreRunEUpdateProjectCmd(a.cfg, persistentPreRunEUpdateCmd)

	return update.NewUpdateProjectCmd(persistentPreRunEUpdateProjectCmd, cmdSettings, cmdLanguages), nil
}
