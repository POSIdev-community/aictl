package context

import (
	"context"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
)

type noopClearUC struct{}

func (noopClearUC) Execute(context.Context, bool) error { return nil }

type noopShowUC struct{}

func (noopShowUC) Execute(context.Context, bool, bool) error { return nil }

type noopSetUC struct{}

func (noopSetUC) Execute() error { return nil }

type noopUnsetUC struct{}

func (noopUnsetUC) Execute(bool, bool, bool, bool, bool) error { return nil }

func emptyCfg() *config.Config {
	return config.NewConfig(config.Uri{}, "", false, uuid.Nil, uuid.Nil)
}

func newCtxRoot(clear UseCaseConfigClear, set UseCaseConfigSet, show UseCaseConfigShow, unset UseCaseConfigUnset) *CmdContext {
	return NewContextCmd(
		NewConfigClearCommand(clear),
		NewConfigSetCommand(emptyCfg(), set),
		NewConfigShowCommand(show),
		NewConfigUnsetCommand(unset),
	)
}
