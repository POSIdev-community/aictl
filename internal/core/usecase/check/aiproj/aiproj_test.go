package aiproj

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/apperror"
)

type fakeCLI struct {
	stdout []string
	stderr []string
}

func (f *fakeCLI) ShowTextf(_ context.Context, format string, a ...any) {
	f.stderr = append(f.stderr, fmt.Sprintf(format, a...))
}

func (f *fakeCLI) ReturnText(_ context.Context, text string) {
	f.stdout = append(f.stdout, text)
}

func TestUseCase_OK_JSON(t *testing.T) {
	t.Parallel()

	cli := &fakeCLI{}
	uc, err := NewUseCase(cli)
	require.NoError(t, err)

	raw := []byte(`{
		"Version": "1.11",
		"ProjectName": "demo",
		"ProgrammingLanguages": ["Go"],
		"ScanModules": ["StaticCodeAnalysis"]
	}`)

	require.NoError(t, uc.Execute(context.Background(), raw, "", true))
	require.Len(t, cli.stdout, 1)
	require.Empty(t, cli.stderr)

	var payload map[string]any
	require.NoError(t, json.Unmarshal([]byte(cli.stdout[0]), &payload))
	require.Equal(t, true, payload["ok"])
	require.Equal(t, "1.11", payload["version"])
	require.Equal(t, "demo", payload["projectName"])
	require.Equal(t, []any{"Go"}, payload["languages"])
	require.Equal(t, []any{}, payload["errors"])
}

func TestUseCase_OK_VerboseHuman(t *testing.T) {
	t.Parallel()

	cli := &fakeCLI{}
	uc, err := NewUseCase(cli)
	require.NoError(t, err)

	raw := []byte(`{
		"Version": "1.11",
		"ProjectName": "demo",
		"ProgrammingLanguages": ["Go"],
		"ScanModules": ["StaticCodeAnalysis"]
	}`)

	require.NoError(t, uc.Execute(context.Background(), raw, "", false))
	require.Empty(t, cli.stdout)
	require.Equal(t, []string{"aiproj ok (version 1.11)"}, cli.stderr)
}

func TestUseCase_Fail_JSON(t *testing.T) {
	t.Parallel()

	cli := &fakeCLI{}
	uc, err := NewUseCase(cli)
	require.NoError(t, err)

	raw := []byte(`{
		"Version": "1.11",
		"ProjectName": 123,
		"ProgrammingLanguages": ["Go"],
		"ScanModules": ["StaticCodeAnalysis"]
	}`)

	err = uc.Execute(context.Background(), raw, "1.11", true)
	require.Error(t, err)

	var failErr *apperror.FailError
	require.ErrorAs(t, err, &failErr)
	require.Equal(t, "aiproj schema invalid", failErr.Error())
	require.Len(t, cli.stdout, 1)

	var payload map[string]any
	require.NoError(t, json.Unmarshal([]byte(cli.stdout[0]), &payload))
	require.Equal(t, false, payload["ok"])
	require.Equal(t, "1.11", payload["version"])
	require.Nil(t, payload["projectName"])
	require.Equal(t, []any{"Go"}, payload["languages"])
	require.NotEmpty(t, payload["errors"])
}

func TestUseCase_Fail_Human(t *testing.T) {
	t.Parallel()

	cli := &fakeCLI{}
	uc, err := NewUseCase(cli)
	require.NoError(t, err)

	raw := []byte(`{
		"Version": "1.10",
		"ProjectName": "demo",
		"ProgrammingLanguages": ["Go"],
		"ScanModules": ["StaticCodeAnalysis"]
	}`)

	err = uc.Execute(context.Background(), raw, "1.11", false)
	require.Error(t, err)

	var failErr *apperror.FailError
	require.ErrorAs(t, err, &failErr)
	require.Equal(t, `aiproj Version is "1.10", expected "1.11"`, failErr.Error())
	require.Empty(t, cli.stdout)
}
