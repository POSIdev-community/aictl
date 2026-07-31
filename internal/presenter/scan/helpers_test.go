package scan

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/scantype"
	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

func resetScanStartFlags() {
	scanLabel = ""
	fullScan = false
}

type noopAwaitUC struct{}

func (noopAwaitUC) Execute(context.Context, uuid.UUID, bool) error { return nil }

type noopCheckPoliciesUC struct{}

func (noopCheckPoliciesUC) Execute(context.Context, uuid.UUID, bool) error { return nil }

type noopStartBranchUC struct{}

func (noopStartBranchUC) Execute(context.Context, string, scantype.Type) error { return nil }

type noopStartProjectUC struct{}

func (noopStartProjectUC) Execute(context.Context, string, scantype.Type) error { return nil }

type noopStopUC struct{}

func (noopStopUC) Execute(context.Context, uuid.UUID) error { return nil }

func buildScanRoot(
	t *testing.T,
	cfg *config.Config,
	await UseCaseScanAwait,
	check UseCaseScanCheckPolicies,
	startBranch UseCaseScanStartBranch,
	startProject UseCaseScanStartProject,
	stop UseCaseScanStop,
) *CmdScan {
	t.Helper()
	require.NotNil(t, cfg)
	preScan := NewPersistentPreRunEScanCmd(cfg)
	preStart := NewPersistentPreRunEScanStartCmd(preScan)
	start := NewScanStartCmd(preStart, NewScanStartBranchCmd(cfg, startBranch), NewScanStartProjectCmd(cfg, startProject))
	return NewScanCmd(preScan, NewScanAwaitCmd(cfg, await), NewScanCheckPoliciesCmd(cfg, check), start, NewScanStopCmd(stop))
}

func mustScanCfg(t *testing.T) *config.Config {
	t.Helper()
	return cmdtest.MustCfg(t)
}
