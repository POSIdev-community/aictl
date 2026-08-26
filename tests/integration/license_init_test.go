package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/adapter/ai"
)

func TestAdapterInitialize_LicenseFailure(t *testing.T) {
	t.Parallel()

	const apiToken = "api-token-for-tests"

	tests := []struct {
		name        string
		version     string
		license     []LicenseStep
		wantSubstr  string
		wantSuccess bool
	}{
		{
			name:       "license_http_500",
			version:    "6.1.0",
			license:    []LicenseStep{{Status: http.StatusInternalServerError}},
			wantSubstr: "check license",
		},
		{
			name:       "license_http_403",
			version:    "6.0.5",
			license:    []LicenseStep{{Status: http.StatusForbidden}},
			wantSubstr: "check license",
		},
		{
			name:       "license_invalid",
			version:    "5.4.0",
			license:    []LicenseStep{LicenseInvalid()},
			wantSubstr: "license is invalid",
		},
		{
			name:        "license_valid",
			version:     "6.1.0",
			license:     []LicenseStep{LicenseOK()},
			wantSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := NewInitServer(t, InitServerConfig{
				APIToken: apiToken,
				Version:  tt.version,
				Signin:   []Step{OKSignin("access-1", "refresh-1")},
				License:  tt.license,
			})

			adapter, err := ai.NewAdapter(mustConfig(t, srv.URL, apiToken))
			require.NoError(t, err)

			err = adapter.Initialize(t.Context())
			if tt.wantSuccess {
				require.NoError(t, err)
				require.Equal(t, 1, srv.LicenseCalls)

				return
			}

			require.Error(t, err)
			require.Contains(t, err.Error(), tt.wantSubstr)
			require.Equal(t, 1, srv.LicenseCalls)
		})
	}
}

func TestAdapterInitializeWithRetry_LicenseFailure(t *testing.T) {
	t.Parallel()

	const apiToken = "api-token-for-tests"

	srv := NewInitServer(t, InitServerConfig{
		APIToken: apiToken,
		Version:  "6.1.0",
		Signin:   []Step{OKSignin("access-1", "refresh-1")},
		License:  []LicenseStep{LicenseInvalid()},
	})

	adapter, err := ai.NewAdapter(mustConfig(t, srv.URL, apiToken))
	require.NoError(t, err)

	err = adapter.InitializeWithRetry(t.Context())
	require.Error(t, err)
	require.Contains(t, err.Error(), "initialize ai adapter")
	require.Contains(t, err.Error(), "license is invalid")
}
