package di

import (
	checkAiproj "github.com/POSIdev-community/aictl/internal/core/usecase/check/aiproj"
	checkPresenter "github.com/POSIdev-community/aictl/internal/presenter/check"
)

func buildCheckCmd(a *adapters) (*checkPresenter.CmdCheck, error) {
	uc, err := checkAiproj.NewUseCase(a.cli)
	if err != nil {
		return nil, err
	}

	persistentPreRunE := checkPresenter.NewPersistentPreRunECheckCmd()
	cmdAiproj := checkPresenter.NewCheckAiprojCmd(uc)

	return checkPresenter.NewCheckCmd(persistentPreRunE, cmdAiproj), nil
}
