package clear

import (
	"github.com/POSIdev-community/aictl/internal/core/port"
	"github.com/POSIdev-community/aictl/pkg/errs"
)

type UseCase struct {
	configAdapter port.Config
	cliAdapter    port.Cli
}

func NewUseCase(configAdapter port.Config, cliAdapter port.Cli) (*UseCase, error) {
	if configAdapter == nil {
		return nil, errs.NewValidationRequiredError("configAdapter")
	}

	if cliAdapter == nil {
		return nil, errs.NewValidationRequiredError("cliAdapter")
	}

	return &UseCase{configAdapter, cliAdapter}, nil
}

func (u *UseCase) Execute() error {
	ok, err := u.cliAdapter.AskConfirmation(
		"Are you sure you want to delete the existing configuration?")
	if err != nil {
		return err
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
