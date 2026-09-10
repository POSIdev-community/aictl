//go:build e2e

package e2e

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadStandsExample(t *testing.T) {
	root, err := RepoRoot()
	require.NoError(t, err)

	path := filepath.Join(root, "tests", "e2e", "stands.example.yaml")
	stands, err := LoadStands(path)
	require.Error(t, err, "example config must use placeholder tokens")
	require.Nil(t, stands)
}

func TestStandVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		stand   string
		want    string
		wantErr bool
	}{
		{name: "54", stand: "5.4", want: "5.4"},
		{name: "60", stand: "6.0", want: "6.0"},
		{name: "61", stand: "6.1", want: "6.1"},
		{name: "62", stand: "6.2", want: "6.2"},
		{name: "63", stand: "6.3", want: "6.3"},
		{name: "invalid", stand: "prod", wantErr: true},
		{name: "major only", stand: "6", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := StandVersion(tt.stand)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
