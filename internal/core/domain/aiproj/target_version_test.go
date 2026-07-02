package aiproj

import (
	"testing"

	"github.com/POSIdev-community/aictl/internal/core/domain/version"
	"github.com/stretchr/testify/require"
)

func TestTargetForServer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		server string
		want   string
	}{
		{name: "5.4", server: "5.4.0", want: "1.9"},
		{name: "5.9", server: "5.9.9", want: "1.9"},
		{name: "6.0", server: "6.0.0", want: "1.10"},
		{name: "6.1", server: "6.1.0", want: "1.11"},
		{name: "6.5", server: "6.5.0", want: "1.11"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			serverVersion, err := version.NewVersion(tt.server)
			require.NoError(t, err)

			got := getTargetVersion(serverVersion)

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParseAndMigrateForServer(t *testing.T) {
	t.Parallel()

	input := []byte(`{
		"Version": "1.8",
		"ProjectName": "demo",
		"ProgrammingLanguages": ["Go"],
		"ScanModules": ["StaticCodeAnalysis"],
		"GoSettings": {"CustomParameters": "+v"}
	}`)

	serverVersion, err := version.NewVersion("6.1.0")
	require.NoError(t, err)

	parsed, err := ParseAndMigrateForServer(input, serverVersion)
	require.NoError(t, err)
	require.Equal(t, "1.11", parsed.AiprojVersion())
	require.NotNil(t, parsed.V111)
	require.Equal(t, "demo", parsed.V111.ProjectName)
	require.NotNil(t, parsed.V111.GoSettings)
	require.NotNil(t, parsed.V111.GoSettings.CustomParameters)
	require.Equal(t, "+v", *parsed.V111.GoSettings.CustomParameters)
}
