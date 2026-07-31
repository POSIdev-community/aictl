package get

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeScanningUC struct {
	called int
	filter *uuid.UUID
}

func (f *fakeScanningUC) Execute(_ context.Context, projectFilter *uuid.UUID) error {
	f.called++
	f.filter = projectFilter
	return nil
}

func TestGetScanningCmd(t *testing.T) {
	t.Cleanup(resetGetPackageFlags)
	resetGetPackageFlags()
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	cfg := mustGetCfg(t)
	uc := &fakeScanningUC{}
	ucs := defaultGetUCs()
	ucs.scanning = uc
	root := buildGetRoot(t, cfg, ucs)
	require.NoError(t, cmdtest.Execute(t, root.Command, "scanning", "-p", projectID.String()))
	require.Equal(t, projectID, *uc.filter)

	t.Run("without_project_filter", func(t *testing.T) {
		resetGetPackageFlags()
		uc2 := &fakeScanningUC{}
		ucs2 := defaultGetUCs()
		ucs2.scanning = uc2
		root2 := buildGetRoot(t, mustGetCfg(t), ucs2)
		require.NoError(t, cmdtest.Execute(t, root2.Command, "scanning"))
		require.Nil(t, uc2.filter)
	})
}
