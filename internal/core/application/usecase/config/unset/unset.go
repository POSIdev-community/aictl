package unset

import (
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/port"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
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
