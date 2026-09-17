package create

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeCreateSbomProjectUC struct {
	called int
	name   string
	path   string
	safe   bool
}

func (f *fakeCreateSbomProjectUC) Execute(_ context.Context, projectName, sbomPath string, safe bool) error {
	f.called++
	f.name, f.path, f.safe = projectName, sbomPath, safe
	return nil
}

func TestCreateSbomProjectCmd(t *testing.T) {
	t.Cleanup(resetSafeFlag)

	sbomFile := filepath.Join(t.TempDir(), "sbom.json")
	require.NoError(t, os.WriteFile(sbomFile, []byte(`{}`), 0o644))

	t.Run("with_name", func(t *testing.T) {
		resetSafeFlag()
		cfg := cmdtest.MustCfg(t)
		uc := &fakeCreateSbomProjectUC{}
		root := NewCreateCmd(cfg, NewCreateBranchCmd(cfg, noopCreateBranchUC{}), NewCreateProjectCmd(noopCreateProjectUC{}), NewCreateSbomProjectCmd(uc))
		require.NoError(t, cmdtest.Execute(t, root.Command, "sbom-project", "my-sbom", "--file", sbomFile, "--safe"))
		require.Equal(t, 1, uc.called)
		require.Equal(t, "my-sbom", uc.name)
		require.Equal(t, sbomFile, uc.path)
		require.True(t, uc.safe)
	})

	t.Run("stdin_dash", func(t *testing.T) {
		resetSafeFlag()
		cfg := cmdtest.MustCfg(t)
		uc := &fakeCreateSbomProjectUC{}
		root := NewCreateCmd(cfg, NewCreateBranchCmd(cfg, noopCreateBranchUC{}), NewCreateProjectCmd(noopCreateProjectUC{}), NewCreateSbomProjectCmd(uc))
		cmdtest.WithStdin(t, "from-stdin\n", func() {
			require.NoError(t, cmdtest.Execute(t, root.Command, "sbom-project", "-", "--file", sbomFile))
		})
		require.Equal(t, "from-stdin", uc.name)
		require.False(t, uc.safe)
	})

	t.Run("missing_name", func(t *testing.T) {
		resetSafeFlag()
		cfg := cmdtest.MustCfg(t)
		uc := &fakeCreateSbomProjectUC{}
		root := NewCreateCmd(cfg, NewCreateBranchCmd(cfg, noopCreateBranchUC{}), NewCreateProjectCmd(noopCreateProjectUC{}), NewCreateSbomProjectCmd(uc))
		err := cmdtest.Execute(t, root.Command, "sbom-project", "--file", sbomFile)
		require.Error(t, err)
		var requiredErr *validation.RequiredError
		require.True(t, errors.As(err, &requiredErr))
		require.Equal(t, "project-name", requiredErr.Field)
		require.Equal(t, 0, uc.called)
	})
}
