package unset

import (
	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
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
		return nil, validation.NewRequiredError("configAdapter")
	}

	return &UseCase{configAdapter, cfg}, nil
}

func (u *UseCase) Execute(uriUnset, tokenUnset, tlsUnset, cacertUnset, projectIdUnset, branchIdUnset bool) error {
	uri := u.cfg.Uri()
	token := u.cfg.Token()
	tls := u.cfg.TLSSkip()
	caCert := u.cfg.CACertPath()
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

	if cacertUnset {
		caCert = ""
	}

	if projectIdUnset {
		projectId = uuid.Nil
	}

	if branchIdUnset {
		branchId = uuid.Nil
	}

	newCfg := config.NewConfig(uri, token, tls, projectId, branchId)
	if err := newCfg.ApplyCACertPath(caCert); err != nil {
		return err
	}

	if err := u.configAdapter.StoreContext(newCfg); err != nil {
		return err
	}

	return nil
}
