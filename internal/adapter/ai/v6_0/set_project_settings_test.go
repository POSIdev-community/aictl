package v6_0

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
	clientai "github.com/POSIdev-community/aictl/pkg/clientai/v6_0"
)

// fakeSettingsServer records what SetProjectSettings actually sends to AIE.
type fakeSettingsServer struct {
	mu             sync.Mutex
	blackBoxActive bool
	settingsBody   map[string]any
	blackBoxBody   map[string]any
	blackBoxPuts   int
}

func (f *fakeSettingsServer) client(t *testing.T, projectID uuid.UUID) *ClientAI60 {
	base := "/api/projects/" + projectID.String()

	return newSettingsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()

		switch {
		case r.Method == http.MethodPut && r.URL.Path == base+"/settings":
			require.NoError(t, json.NewDecoder(r.Body).Decode(&f.settingsBody))
		case r.Method == http.MethodPut && r.URL.Path == base+"/preferredAgentsSettings":
		case r.Method == http.MethodGet && r.URL.Path == base+"/blackBoxSettings":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"isActive": f.blackBoxActive})

			return
		case r.Method == http.MethodPut && r.URL.Path == base+"/blackBoxSettings":
			f.blackBoxPuts++
			require.NoError(t, json.NewDecoder(r.Body).Decode(&f.blackBoxBody))
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)

			return
		}
		w.WriteHeader(http.StatusOK)
	}))
}

var settingsTestProjectID = uuid.MustParse("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")

func TestSetProjectSettingsSendsSkipGitIgnoreAndGoDownloadDependencies(t *testing.T) {
	t.Parallel()

	f := &fakeSettingsServer{}
	s := settings.ScanSettings{
		SkipGitIgnoreFiles: true,
		GoSettings:         settings.GoSettings{DownloadDependencies: true},
	}

	require.NoError(t, f.client(t, settingsTestProjectID).SetProjectSettings(t.Context(), settingsTestProjectID, &s))

	assert.Equal(t, true, f.settingsBody["skipGitIgnoreFiles"])
	goSettings, ok := f.settingsBody["goSettings"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, true, goSettings["downloadDependencies"])
	assert.Zero(t, f.blackBoxPuts, "black box is neither enabled nor active: no PUT expected")
}

func TestSetProjectSettingsSendsBlackBoxAuthenticationType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		auth     *settings.BlackBoxAuthentication
		wantType string
	}{
		{
			name: "explicit Http",
			auth: &settings.BlackBoxAuthentication{
				Type: "Http",
				Http: &settings.BlackBoxHTTPAuth{Login: "admin", Password: "secret", ValidationAddress: "http://10.0.0.1:8081/profile"},
			},
			wantType: "Http",
		},
		{
			name:     "type inferred from Form block",
			auth:     &settings.BlackBoxAuthentication{Form: &settings.BlackBoxFormAuth{Login: "admin"}},
			wantType: "Form",
		},
		{
			name:     "type inferred from Cookie block",
			auth:     &settings.BlackBoxAuthentication{Cookie: &settings.BlackBoxCookieAuth{Cookie: "sid=1"}},
			wantType: "RawCookie",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &fakeSettingsServer{}
			s := settings.ScanSettings{
				BlackBoxEnabled:  true,
				BlackBoxSettings: settings.BlackBoxSettings{Site: "http://10.0.0.1:8081/", Authentication: tt.auth},
			}

			require.NoError(t, f.client(t, settingsTestProjectID).SetProjectSettings(t.Context(), settingsTestProjectID, &s))

			require.Equal(t, 1, f.blackBoxPuts)
			assert.Equal(t, true, f.blackBoxBody["isActive"])
			auth, ok := f.blackBoxBody["authentication"].(map[string]any)
			require.True(t, ok)
			assert.Equal(t, tt.wantType, auth["type"])
		})
	}
}

func TestSetProjectSettingsDisablesActiveBlackBox(t *testing.T) {
	t.Parallel()

	f := &fakeSettingsServer{blackBoxActive: true}
	s := settings.ScanSettings{BlackBoxEnabled: false}

	require.NoError(t, f.client(t, settingsTestProjectID).SetProjectSettings(t.Context(), settingsTestProjectID, &s))

	require.Equal(t, 1, f.blackBoxPuts)
	assert.Equal(t, false, f.blackBoxBody["isActive"])
}

func newSettingsTestClient(t *testing.T, h http.Handler) *ClientAI60 {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)

	base := common.NewBaseClient()
	base.HttpClient = srv.Client()
	c := NewAiClient(base)
	api, err := clientai.NewClientWithResponses(srv.URL, clientai.WithHTTPClient(srv.Client()))
	require.NoError(t, err)
	c.ClientWithResponses = api

	return c
}
