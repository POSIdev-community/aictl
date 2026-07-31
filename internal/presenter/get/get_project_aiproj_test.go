package get

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeProjectAiprojUC struct {
	called int
	out    string
}

func (f *fakeProjectAiprojUC) Execute(_ context.Context, outputPath string) error {
	f.called++
	f.out = outputPath
	return nil
}

func TestGetProjectAiprojCmd(t *testing.T) {
	t.Cleanup(resetGetPackageFlags)
	resetGetPackageFlags()
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	cfg := mustGetCfg(t)
	out := filepath.Join(t.TempDir(), "out.aiproj")

	t.Run("output_flag", func(t *testing.T) {
		resetGetPackageFlags()
		uc := &fakeProjectAiprojUC{}
		ucs := defaultGetUCs()
		ucs.projectAiproj = uc
		root := buildGetRoot(t, cfg, ucs)
		require.NoError(t, cmdtest.Execute(t, root.Command, "project", "aiproj", "-p", projectID.String(), "-o", out))
		require.Equal(t, out, uc.out)
		require.Equal(t, projectID, cfg.ProjectId())
	})

	t.Run("force_required_when_exists", func(t *testing.T) {
		resetGetPackageFlags()
		exist := filepath.Join(t.TempDir(), "exists.aiproj")
		require.NoError(t, os.WriteFile(exist, []byte("x"), 0o644))
		uc := &fakeProjectAiprojUC{}
		ucs := defaultGetUCs()
		ucs.projectAiproj = uc
		root := buildGetRoot(t, mustGetCfg(t), ucs)
		require.Error(t, cmdtest.Execute(t, root.Command, "project", "aiproj", "-p", projectID.String(), "-o", exist))
		require.Equal(t, 0, uc.called)
	})

	t.Run("force_overwrites", func(t *testing.T) {
		resetGetPackageFlags()
		exist := filepath.Join(t.TempDir(), "exists.aiproj")
		require.NoError(t, os.WriteFile(exist, []byte("x"), 0o644))
		uc := &fakeProjectAiprojUC{}
		ucs := defaultGetUCs()
		ucs.projectAiproj = uc
		root := buildGetRoot(t, mustGetCfg(t), ucs)
		require.NoError(t, cmdtest.Execute(t, root.Command, "project", "aiproj", "-p", projectID.String(), "-o", exist, "-f"))
		require.Equal(t, exist, uc.out)
	})
}
