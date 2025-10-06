package set

import (
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
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

func (u *UseCase) Execute(cfg *config.Config) error {
	if err := u.configAdapter.StoreContext(cfg); err != nil {
		return err
	}

	return nil
}
