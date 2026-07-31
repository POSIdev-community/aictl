package get

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeBranchUC struct {
	called   int
	branchID uuid.UUID
}

func (f *fakeBranchUC) Execute(_ context.Context, id uuid.UUID) error {
	f.called++
	f.branchID = id
	return nil
}

func TestGetBranchCmd(t *testing.T) {
	t.Cleanup(resetGetPackageFlags)
	resetGetPackageFlags()
	branchID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	uc := &fakeBranchUC{}
	ucs := defaultGetUCs()
	ucs.branch = uc
	root := buildGetRoot(t, mustGetCfg(t), ucs)
	require.NoError(t, cmdtest.Execute(t, root.Command, "branch", branchID.String()))
	require.Equal(t, branchID, uc.branchID)

	t.Run("stdin_dash", func(t *testing.T) {
		resetGetPackageFlags()
		uc2 := &fakeBranchUC{}
		ucs2 := defaultGetUCs()
		ucs2.branch = uc2
		root2 := buildGetRoot(t, mustGetCfg(t), ucs2)
		cmdtest.WithStdin(t, branchID.String()+"\n", func() {
			require.NoError(t, cmdtest.Execute(t, root2.Command, "branch", "-"))
		})
		require.Equal(t, branchID, uc2.branchID)
	})
}
