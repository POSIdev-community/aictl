package settings_test

import (
	"encoding/json"
	"fmt"
	"testing"

	v110 "github.com/POSIdev-community/aiproj/model/v1_10"
	v19 "github.com/POSIdev-community/aiproj/model/v1_9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
)

func fromV110(t *testing.T, s *settings.ScanSettings, raw string) {
	t.Helper()
	var p v110.AIProj
	require.NoError(t, json.Unmarshal([]byte(raw), &p))
	require.NoError(t, s.UpdateFromV110(&p))
}

func TestUpdateFromAiprojGoDownloadDependencies(t *testing.T) {
	t.Parallel()

	raw := `{"Version":"%s","ProjectName":"demo","ProgrammingLanguages":["Go"],"ScanModules":["StaticCodeAnalysis"],
		"GoSettings":{"DownloadDependencies":true}}`

	t.Run("1.10", func(t *testing.T) {
		t.Parallel()
		var s settings.ScanSettings
		fromV110(t, &s, fmt.Sprintf(raw, "1.10"))
		assert.True(t, s.GoSettings.DownloadDependencies)
	})

	t.Run("1.9", func(t *testing.T) {
		t.Parallel()
		var p v19.AIProj
		require.NoError(t, json.Unmarshal([]byte(fmt.Sprintf(raw, "1.9")), &p))
		var s settings.ScanSettings
		require.NoError(t, s.UpdateFromV19(&p))
		assert.True(t, s.GoSettings.DownloadDependencies)
	})
}

func TestUpdateFromAiprojBlackBoxEnabledOnlyByScanModules(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		initial     bool
		raw         string
		wantEnabled bool
		wantSite    string
	}{
		{
			name: "block without BlackBox module does not enable",
			raw: `{"Version":"1.10","ProjectName":"demo","ProgrammingLanguages":["Python"],"ScanModules":["StaticCodeAnalysis"],
				"BlackBoxSettings":{"Site":"http://10.0.0.1:8081/"}}`,
			wantEnabled: false,
			wantSite:    "http://10.0.0.1:8081/",
		},
		{
			name: "empty block without BlackBox module does not enable",
			raw: `{"Version":"1.10","ProjectName":"demo","ProgrammingLanguages":["Python"],"ScanModules":["StaticCodeAnalysis"],
				"BlackBoxSettings":{}}`,
			wantEnabled: false,
		},
		{
			name:    "previously enabled is turned off when module is not listed",
			initial: true,
			raw: `{"Version":"1.10","ProjectName":"demo","ProgrammingLanguages":["Python"],
				"ScanModules":["StaticCodeAnalysis"]}`,
			wantEnabled: false,
		},
		{
			name: "BlackBox module enables",
			raw: `{"Version":"1.10","ProjectName":"demo","ProgrammingLanguages":["Python"],"ScanModules":["StaticCodeAnalysis","BlackBox"],
				"BlackBoxSettings":{"Site":"http://10.0.0.1:8081/"}}`,
			wantEnabled: true,
			wantSite:    "http://10.0.0.1:8081/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := settings.ScanSettings{BlackBoxEnabled: tt.initial}
			fromV110(t, &s, tt.raw)
			assert.Equal(t, tt.wantEnabled, s.BlackBoxEnabled)
			assert.True(t, s.BlackBoxManaged, "ScanModules from aiproj define the black box state")
			assert.Equal(t, tt.wantSite, s.BlackBoxSettings.Site)
		})
	}
}
