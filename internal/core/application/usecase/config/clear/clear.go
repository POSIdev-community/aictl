package clear

import (
	"github.com/POSIdev-community/aictl/pkg/errs"
)

type CFG interface {
	ClearCurrentContext() error
}

type CLI interface {
	ShowText(text string)
	AskConfirmation(question string) (bool, error)
}

type UseCase struct {
	configAdapter CFG
	cliAdapter    CLI
}

func NewUseCase(configAdapter CFG, cliAdapter CLI) (*UseCase, error) {
	if configAdapter == nil {
		return nil, errs.NewValidationRequiredError("configAdapter")
	}

	if cliAdapter == nil {
		return nil, errs.NewValidationRequiredError("cliAdapter")
	}

	return &UseCase{configAdapter, cliAdapter}, nil
}

func (u *UseCase) Execute(skipConfirm bool) error {
	var (
		ok  = skipConfirm
		err error
	)

	if !ok {
		ok, err = u.cliAdapter.AskConfirmation(
			"Are you sure you want to delete the existing configuration?")
		if err != nil {
			return err
		}
	}

	if !ok {
		u.cliAdapter.ShowText("Cancelled")

		return nil
	}

	if err := u.configAdapter.ClearCurrentContext(); err != nil {
		return err
	}

	return nil
}
