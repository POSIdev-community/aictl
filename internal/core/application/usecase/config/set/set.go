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
	cfg           *config.Config
}

func NewUseCase(configAdapter CFG, cfg *config.Config) (*UseCase, error) {
	if configAdapter == nil {
		return nil, errs.NewValidationRequiredError("configAdapter")
	}

	return &UseCase{configAdapter, cfg}, nil
}

func (u *UseCase) Execute() error {
	if err := u.configAdapter.StoreContext(u.cfg); err != nil {
		return err
	}

	return nil
}
