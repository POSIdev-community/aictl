package set

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeSetSettingsUC struct {
	called    int
	rawAiproj []byte
}

func (f *fakeSetSettingsUC) Execute(_ context.Context, rawAiproj []byte) error {
	f.called++
	f.rawAiproj = append([]byte(nil), rawAiproj...)
	return nil
}

func TestSetProjectSettingsCmd(t *testing.T) {
	t.Cleanup(resetSetProjectID)
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	t.Run("from_arg", func(t *testing.T) {
		resetSetProjectID()
		uc := &fakeSetSettingsUC{}
		root := buildSetRoot(t, uc, noopSetPoliciesUC{}, &fakeSetExclusionsUC{})
		payload := `{"Version":"1.0"}`
		require.NoError(t, cmdtest.Execute(t, root.Command, "project", "settings", payload, "-p", projectID.String()))
		require.JSONEq(t, payload, string(uc.rawAiproj))
	})

	t.Run("from_file", func(t *testing.T) {
		resetSetProjectID()
		path := filepath.Join(t.TempDir(), "aiproj.json")
		require.NoError(t, os.WriteFile(path, []byte(`{"Version":"1.1"}`), 0o644))
		uc := &fakeSetSettingsUC{}
		root := buildSetRoot(t, uc, noopSetPoliciesUC{}, &fakeSetExclusionsUC{})
		require.NoError(t, cmdtest.Execute(t, root.Command, "project", "settings", "-f", path, "-p", projectID.String()))
		require.JSONEq(t, `{"Version":"1.1"}`, string(uc.rawAiproj))
	})

	t.Run("stdin_arg_dash", func(t *testing.T) {
		resetSetProjectID()
		uc := &fakeSetSettingsUC{}
		root := buildSetRoot(t, uc, noopSetPoliciesUC{}, &fakeSetExclusionsUC{})
		cmdtest.WithStdin(t, `{"Version":"2.0"}`+"\n", func() {
			require.NoError(t, cmdtest.Execute(t, root.Command, "project", "settings", "-", "-p", projectID.String()))
		})
		require.JSONEq(t, `{"Version":"2.0"}`, string(uc.rawAiproj))
	})
}
