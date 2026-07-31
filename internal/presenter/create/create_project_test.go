package create

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
	"github.com/POSIdev-community/aictl/pkg/gitignore"
)

type fakeCreateProjectUC struct {
	called int
	name   string
	safe   bool
}

func (f *fakeCreateProjectUC) Execute(_ context.Context, projectName string, safe bool) error {
	f.called++
	f.name, f.safe = projectName, safe
	return nil
}

type noopCreateBranchUC struct{}

func (noopCreateBranchUC) Execute(context.Context, *config.Config, string, string, bool, gitignore.Exclusions, string) error {
	return nil
}

func TestCreateProjectCmd(t *testing.T) {
	t.Cleanup(resetSafeFlag)
	cfg := cmdtest.MustCfg(t)

	t.Run("with_safe", func(t *testing.T) {
		resetSafeFlag()
		uc := &fakeCreateProjectUC{}
		root := NewCreateCmd(cfg, NewCreateBranchCmd(cfg, noopCreateBranchUC{}), NewCreateProjectCmd(uc))
		require.NoError(t, cmdtest.Execute(t, root.Command, "project", "my-app", "--safe"))
		require.Equal(t, 1, uc.called)
		require.Equal(t, "my-app", uc.name)
		require.True(t, uc.safe)
	})

	t.Run("stdin_dash", func(t *testing.T) {
		resetSafeFlag()
		uc := &fakeCreateProjectUC{}
		root := NewCreateCmd(cfg, NewCreateBranchCmd(cfg, noopCreateBranchUC{}), NewCreateProjectCmd(uc))
		cmdtest.WithStdin(t, "from-stdin\n", func() {
			require.NoError(t, cmdtest.Execute(t, root.Command, "project", "-"))
		})
		require.Equal(t, "from-stdin", uc.name)
		require.False(t, uc.safe)
	})
}
