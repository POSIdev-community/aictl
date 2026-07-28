package ai

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	"github.com/POSIdev-community/aictl/internal/adapter/ai/v5_x"
	"github.com/POSIdev-community/aictl/internal/adapter/ai/v6_0"
	"github.com/POSIdev-community/aictl/internal/adapter/ai/v6_x"
)

var clientInitializers = []common.Initializer{
	v6_x.Initializer,
	v6_0.Initializer,
	v5_x.Initializer,
}

func (a *Adapter) Initialize(ctx context.Context) error {
	state := common.InitState{}

	for _, init := range clientInitializers {
		client, nextState, matched, err := init.TryInitialize(ctx, a.baseClient, a.cfg, state)
		state = nextState
		if err != nil {
			return err
		}
		if !matched {
			a.baseClient.Reset()
			continue
		}

		a.serverVersion = state.Version
		a.activeClient = client

		//return a.activeClient.CheckLicense(ctx)

		return nil
	}

	return fmt.Errorf("initialize ai client: no compatible client found")
}
