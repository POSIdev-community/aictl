package show

import (
	"fmt"
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

func (u *UseCase) Execute(config *config.Config, json bool, yaml bool) error {
	if json && yaml {
		return fmt.Errorf("cannot use both json and yaml format")
	}

	var (
		str string
		err error
	)

	switch {
	case json:
		str, err = u.configAdapter.StringJson(config)
	case yaml:
		str, err = u.configAdapter.StringYaml(config)
	default:
		str, err = u.configAdapter.String(config)
	}

	if err != nil {
		return err
	}

	u.cliAdapter.ShowText(str)

	return nil
}
