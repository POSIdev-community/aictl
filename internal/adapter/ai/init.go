package ai

import (
	"context"
	"errors"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	"github.com/POSIdev-community/aictl/internal/adapter/ai/v5_x"
	"github.com/POSIdev-community/aictl/internal/adapter/ai/v6_0"
	"github.com/POSIdev-community/aictl/internal/adapter/ai/v6_1"
	"github.com/POSIdev-community/aictl/internal/adapter/ai/v6_x"
)

var clientInitializers = []common.Initializer{
	v6_x.Initializer,
	v6_1.Initializer,
	v6_0.Initializer,
	v5_x.Initializer,
}

func (a *Adapter) Initialize(ctx context.Context) error {
	state := common.InitState{}
	var candidateErrs []error

	for _, init := range clientInitializers {
		client, nextState, matched, err := init.TryInitialize(ctx, a.baseClient, a.cfg, state)
		state = nextState
		if matched {
			a.serverVersion = state.Version
			a.activeClient = client

			return a.activeClient.CheckLicense(ctx)
		}
		if err != nil {
			candidateErrs = append(candidateErrs, err)
		}
		a.baseClient.Reset()
	}

	if len(candidateErrs) > 0 {
		return fmt.Errorf("initialize ai client: no compatible client found: %w", errors.Join(candidateErrs...))
	}

	return fmt.Errorf("initialize ai client: no compatible client found")
}
