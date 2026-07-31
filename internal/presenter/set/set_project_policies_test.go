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

type fakeSetPoliciesUC struct {
	called  int
	rawJSON []byte
}

func (f *fakeSetPoliciesUC) Execute(_ context.Context, rawJSON []byte) error {
	f.called++
	f.rawJSON = append([]byte(nil), rawJSON...)
	return nil
}

func TestSetProjectPoliciesCmd(t *testing.T) {
	t.Cleanup(resetSetProjectID)
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	t.Run("invalid_json", func(t *testing.T) {
		resetSetProjectID()
		uc := &fakeSetPoliciesUC{}
		root := buildSetRoot(t, noopSetSettingsUC{}, uc, &fakeSetExclusionsUC{})
		require.Error(t, cmdtest.Execute(t, root.Command, "project", "policies", "not-json", "-p", projectID.String()))
		require.Equal(t, 0, uc.called)
	})

	t.Run("from_arg", func(t *testing.T) {
		resetSetProjectID()
		uc := &fakeSetPoliciesUC{}
		root := buildSetRoot(t, noopSetSettingsUC{}, uc, &fakeSetExclusionsUC{})
		require.NoError(t, cmdtest.Execute(t, root.Command, "project", "policies", `{"rules":[]}`, "-p", projectID.String()))
		require.JSONEq(t, `{"rules":[]}`, string(uc.rawJSON))
	})

	t.Run("from_file", func(t *testing.T) {
		resetSetProjectID()
		path := filepath.Join(t.TempDir(), "policies.json")
		require.NoError(t, os.WriteFile(path, []byte(`{"ok":true}`), 0o644))
		uc := &fakeSetPoliciesUC{}
		root := buildSetRoot(t, noopSetSettingsUC{}, uc, &fakeSetExclusionsUC{})
		require.NoError(t, cmdtest.Execute(t, root.Command, "project", "policies", "-f", path, "-p", projectID.String()))
		require.JSONEq(t, `{"ok":true}`, string(uc.rawJSON))
	})

	t.Run("stdin_arg_dash", func(t *testing.T) {
		resetSetProjectID()
		uc := &fakeSetPoliciesUC{}
		root := buildSetRoot(t, noopSetSettingsUC{}, uc, &fakeSetExclusionsUC{})
		cmdtest.WithStdin(t, `{"a":1}`+"\n", func() {
			require.NoError(t, cmdtest.Execute(t, root.Command, "project", "policies", "-", "-p", projectID.String()))
		})
		require.JSONEq(t, `{"a":1}`, string(uc.rawJSON))
	})

	t.Run("stdin_file_dash", func(t *testing.T) {
		resetSetProjectID()
		uc := &fakeSetPoliciesUC{}
		root := buildSetRoot(t, noopSetSettingsUC{}, uc, &fakeSetExclusionsUC{})
		cmdtest.WithStdin(t, `{"b":2}`+"\n", func() {
			require.NoError(t, cmdtest.Execute(t, root.Command, "project", "policies", "-f", "-", "-p", projectID.String()))
		})
		require.JSONEq(t, `{"b":2}`, string(uc.rawJSON))
	})
}
