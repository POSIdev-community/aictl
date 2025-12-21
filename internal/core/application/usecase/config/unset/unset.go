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
	cfg           *config.Config
}

func NewUseCase(configAdapter CFG, cfg *config.Config) (*UseCase, error) {
	if configAdapter == nil {
		return nil, errs.NewValidationRequiredError("configAdapter")
	}

	return &UseCase{configAdapter, cfg}, nil
}

func (u *UseCase) Execute(uriUnset, tokenUnset, tlsUnset, projectIdUnset, branchIdUnset bool) error {
	uri := u.cfg.Uri()
	token := u.cfg.Token()
	tls := u.cfg.TLSSkip()
	projectId := u.cfg.ProjectId()
	branchId := u.cfg.BranchId()

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
