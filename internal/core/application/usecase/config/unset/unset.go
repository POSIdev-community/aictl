package unset

import (
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
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

func (u *UseCase) Execute(cfg *config.Config, uriUnset, tokenUnset, tlsUnset, projectIdUnset, branchIdUnset bool) error {
	uri := cfg.Uri()
	token := cfg.Token()
	tls := cfg.TLSSkip()
	projectId := cfg.ProjectId()
	branchId := cfg.BranchId()

	if uriUnset {
		uri = config.Uri{}
	}

	if tokenUnset {
		token = ""
	}

	if tlsUnset {
		tls = false
	}

	if projectIdUnset {
		projectId = uuid.Nil
	}

	if branchIdUnset {
		branchId = uuid.Nil
	}

	if err := u.configAdapter.StoreContext(config.NewConfig(uri, token, tls, projectId, branchId)); err != nil {
		return err
	}

	return nil
}
