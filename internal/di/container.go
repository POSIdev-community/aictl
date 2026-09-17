package di

import (
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter"
)

func InitializeCmd(cfg *config.Config) (*presenter.CmdRoot, error) {
	a, err := newAdapters(cfg)
	if err != nil {
		return nil, err
	}

	cmdContext, err := buildContextCmd(a)
	if err != nil {
		return nil, err
	}

	cmdCreate, err := buildCreateCmd(a)
	if err != nil {
		return nil, err
	}

	cmdDelete, err := buildDeleteCmd(a)
	if err != nil {
		return nil, err
	}

	cmdGet, err := buildGetCmd(a)
	if err != nil {
		return nil, err
	}

	cmdScan, err := buildScanCmd(a)
	if err != nil {
		return nil, err
	}

	cmdSet, err := buildSetCmd(a)
	if err != nil {
		return nil, err
	}

	cmdUpdate, err := buildUpdateCmd(a)
	if err != nil {
		return nil, err
	}

	cmdCheck, err := buildCheckCmd(a)
	if err != nil {
		return nil, err
	}

	return presenter.NewRootCmd(cmdContext, cmdCreate, cmdDelete, cmdGet, cmdScan, cmdSet, cmdUpdate, cmdCheck), nil
}
