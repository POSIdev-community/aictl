package check

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeCheckUC struct {
	called        int
	raw           []byte
	schemaVersion string
	jsonOut       bool
	err           error
}

func (f *fakeCheckUC) Execute(_ context.Context, raw []byte, schemaVersion string, jsonOut bool) error {
	f.called++
	f.raw = append([]byte(nil), raw...)
	f.schemaVersion = schemaVersion
	f.jsonOut = jsonOut

	return f.err
}

func buildCheckRoot(t *testing.T, uc UseCaseCheckAiproj) *CmdCheck {
	t.Helper()
	cmd := NewCheckCmd(NewPersistentPreRunECheckCmd(), NewCheckAiprojCmd(uc))
	cmdtest.SilenceCmd(cmd.Command)

	return cmd
}

func TestCheckAiprojCmd_FromFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "aiproj.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"Version":"1.11"}`), 0o644))

	uc := &fakeCheckUC{}
	root := buildCheckRoot(t, uc)
	require.NoError(t, cmdtest.Execute(t, root.Command, "aiproj", "-f", path, "--schema-version", "1.11", "--json"))
	require.Equal(t, 1, uc.called)
	require.JSONEq(t, `{"Version":"1.11"}`, string(uc.raw))
	require.Equal(t, "1.11", uc.schemaVersion)
	require.True(t, uc.jsonOut)
}

func TestCheckAiprojCmd_MissingInput(t *testing.T) {
	uc := &fakeCheckUC{}
	root := buildCheckRoot(t, uc)
	err := cmdtest.Execute(t, root.Command, "aiproj")
	require.Error(t, err)

	var msgErr *validation.MessageError
	require.ErrorAs(t, err, &msgErr)
	require.Equal(t, "aiproj data required", msgErr.Message)
	require.Equal(t, 0, uc.called)
}

func TestCheckAiprojCmd_UseCaseError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "aiproj.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"Version":"1.11"}`), 0o644))

	uc := &fakeCheckUC{err: validation.NewMessageError("aiproj schema invalid")}
	root := buildCheckRoot(t, uc)
	err := cmdtest.Execute(t, root.Command, "aiproj", "-f", path)
	require.Error(t, err)
	require.ErrorContains(t, err, "aiproj schema invalid")
}
