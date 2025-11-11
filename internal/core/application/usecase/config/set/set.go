package set

import (
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/errs"
)

type CFG interface {
	StoreContext(cfg *config.Config) error
}

type UseCase struct {
	configAdapter CFG
}

func NewUseCase(configAdapter CFG) (*UseCase, error) {
	if configAdapter == nil {
		return nil, errs.NewValidationRequiredError("configAdapter")
	}

	return &UseCase{configAdapter}, nil
}

func (u *UseCase) Execute(cfg *config.Config) error {
	if err := u.configAdapter.StoreContext(cfg); err != nil {
		return err
	}

	return nil
}
