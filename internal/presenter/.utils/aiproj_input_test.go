package _utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

func TestReadAiprojInputFrom_File(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "aiproj.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"Version":"1.11"}`), 0o644))

	raw, err := ReadAiprojInputFrom(path, []string{`{"Version":"ignored"}`}, strings.NewReader(""))
	require.NoError(t, err)
	require.JSONEq(t, `{"Version":"1.11"}`, string(raw))
}

func TestReadAiprojInputFrom_Arg(t *testing.T) {
	t.Parallel()

	raw, err := ReadAiprojInputFrom("", []string{`{"Version":"1.10"}`}, strings.NewReader(""))
	require.NoError(t, err)
	require.JSONEq(t, `{"Version":"1.10"}`, string(raw))
}

func TestReadAiprojInputFrom_Empty(t *testing.T) {
	t.Parallel()

	_, err := ReadAiprojInputFrom("", nil, strings.NewReader(""))
	require.Error(t, err)

	var msgErr *validation.MessageError
	require.ErrorAs(t, err, &msgErr)
	require.Equal(t, "aiproj data required", msgErr.Message)
}

func TestReadAiprojInputFrom_InvalidJSON(t *testing.T) {
	t.Parallel()

	_, err := ReadAiprojInputFrom("", []string{`{`}, strings.NewReader(""))
	require.Error(t, err)

	var msgErr *validation.MessageError
	require.ErrorAs(t, err, &msgErr)
	require.Equal(t, "invalid aiproj data: not valid json", msgErr.Message)
}

func TestReadAiprojInputFrom_MissingFile(t *testing.T) {
	t.Parallel()

	_, err := ReadAiprojInputFrom(filepath.Join(t.TempDir(), "missing.json"), nil, strings.NewReader(""))
	require.Error(t, err)

	var msgErr *validation.MessageError
	require.ErrorAs(t, err, &msgErr)
	require.Contains(t, msgErr.Message, "does not exist")
}

func TestReadAiprojInputFrom_DashStdin(t *testing.T) {
	t.Parallel()

	raw, err := ReadAiprojInputFrom("", []string{"-"}, strings.NewReader(`{"Version":"1.9"}`+"\n"))
	require.NoError(t, err)
	require.JSONEq(t, `{"Version":"1.9"}`, string(raw))
}
