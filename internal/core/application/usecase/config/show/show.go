package show

import (
	"context"
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
	ReturnText(ctx context.Context, text string)
}

type UseCase struct {
	configAdapter CFG
	cliAdapter    CLI
	cfg           *config.Config
}

func NewUseCase(configAdapter CFG, cliAdapter CLI, cfg *config.Config) (*UseCase, error) {
	if configAdapter == nil {
		return nil, errs.NewValidationRequiredError("configAdapter")
	}

	if cliAdapter == nil {
		return nil, errs.NewValidationRequiredError("cliAdapter")
	}

	return &UseCase{configAdapter, cliAdapter, cfg}, nil
}

func (u *UseCase) Execute(ctx context.Context, json bool, yaml bool) error {
	if json && yaml {
		return fmt.Errorf("cannot use both json and yaml format")
	}

	var (
		str string
		err error
	)

	switch {
	case json:
		str, err = u.configAdapter.StringJson(u.cfg)
	case yaml:
		str, err = u.configAdapter.StringYaml(u.cfg)
	default:
		str, err = u.configAdapter.String(u.cfg)
	}

	if err != nil {
		return err
	}

	u.cliAdapter.ReturnText(ctx, str)

	return nil
}
