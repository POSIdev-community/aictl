package show

import (
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/errs"
)

type CFG interface {
	StringJson(config *config.Config) (string, error)
	StringYaml(config *config.Config) (string, error)
	String(config *config.Config) (string, error)
}

type CLI interface {
	ReturnText(text string)
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

	u.cliAdapter.ReturnText(str)

	return nil
}
