package integration

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	"github.com/POSIdev-community/aictl/internal/adapter/ai/v5_x"
	"github.com/POSIdev-community/aictl/internal/adapter/ai/v6_0"
	"github.com/POSIdev-community/aictl/internal/adapter/ai/v6_1"
	"github.com/POSIdev-community/aictl/internal/adapter/ai/v6_x"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/statistic"
)

type scanStatisticClient interface {
	Initialize(ctx context.Context, cfg *config.Config) error
	GetScanStatistic(ctx context.Context, projectId, scanResultId uuid.UUID) (*statistic.Statistic, error)
	GetScanIssues(ctx context.Context, projectId, scanResultId uuid.UUID) ([]statistic.Issue, error)
}

func TestGetScanStatistic_WithTriage(t *testing.T) {
	t.Parallel()

	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	scanID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	issuesJSON := `[
		{"approvalState":"Discard","level":"High"},
		{"approvalState":"None","level":"High"},
		{"approvalState":"Approval","level":"Medium"},
		{"approvalState":"Discard","level":"Low"},
		{"approvalState":"None","level":"Potential"}
	]`

	cases := []struct {
		name      string
		statBody  string
		newClient func(*common.BaseClient) scanStatisticClient
	}{
		{
			name:      "v5_x",
			statBody:  `{"filesTotal":10,"filesScanned":8,"urlsTotal":1,"urlsScanned":1,"scanDuration":"PT00H01M00S","policyState":"None","high":99,"total":99}`,
			newClient: func(base *common.BaseClient) scanStatisticClient { return v5_x.NewAiClient(base) },
		},
		{
			name:      "v6_0",
			statBody:  `{"filesTotal":10,"filesScanned":8,"urlsTotal":1,"urlsScanned":1,"scanDuration":"PT00H01M00S","policyState":"None","high":99,"total":99}`,
			newClient: func(base *common.BaseClient) scanStatisticClient { return v6_0.NewAiClient(base) },
		},
		{
			name:      "v6_1",
			statBody:  `{"filesTotal":10,"filesScanned":8,"urlsTotal":1,"urlsScanned":1,"scanDuration":"PT00H01M00S","policyState":"None","high":99,"total":99}`,
			newClient: func(base *common.BaseClient) scanStatisticClient { return v6_1.NewAiClient(base) },
		},
		{
			name:      "v6_x",
			statBody:  `{"filesTotal":10,"filesScanned":8,"urlsTotal":1,"urlsScanned":1,"scanDuration":"PT00H01M00S","securityPolicyStatus":{"policyState":"None"},"high":99,"total":99}`,
			newClient: func(base *common.BaseClient) scanStatisticClient { return v6_x.NewAiClient(base) },
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := newStatisticAuthServer(t, projectID, scanID, tc.statBody, issuesJSON)
			client := tc.newClient(common.NewBaseClient())
			cfg := mustConfig(t, srv.URL, "api-token-for-tests")

			require.NoError(t, client.Initialize(t.Context(), cfg))

			meta, err := client.GetScanStatistic(t.Context(), projectID, scanID)
			require.NoError(t, err)
			issues, err := client.GetScanIssues(t.Context(), projectID, scanID)
			require.NoError(t, err)

			assertTriageCounts(t, meta, issues)
		})
	}
}

func assertTriageCounts(t *testing.T, meta *statistic.Statistic, issues []statistic.Issue) {
	t.Helper()

	require.Equal(t, int32(10), meta.FilesTotal)
	require.Equal(t, int32(8), meta.FilesScanned)
	require.Equal(t, "None", meta.PolicyState)
	require.Equal(t, int32(0), meta.Total) // severity filled by use case, not adapter

	all := statistic.CountSeverities(issues, false)
	require.Equal(t, int32(2), all.High)
	require.Equal(t, int32(1), all.Medium)
	require.Equal(t, int32(1), all.Low)
	require.Equal(t, int32(1), all.Potential)
	require.Equal(t, int32(5), all.Total())

	triaged := statistic.CountSeverities(issues, true)
	require.Equal(t, int32(1), triaged.High)
	require.Equal(t, int32(1), triaged.Medium)
	require.Equal(t, int32(0), triaged.Low)
	require.Equal(t, int32(1), triaged.Potential)
	require.Equal(t, int32(3), triaged.Total())
}

func newStatisticAuthServer(t *testing.T, projectID, scanID uuid.UUID, statBody, issuesBody string) *httptest.Server {
	t.Helper()

	statPath := fmt.Sprintf("/api/projects/%s/scanResults/%s/statistic", projectID, scanID)
	issuesPath := fmt.Sprintf("/api/projects/%s/scanResults/%s/issues", projectID, scanID)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == SigninPath:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"accessToken":"access-1","refreshToken":"refresh-1"}`))
		case r.Method == http.MethodGet && r.URL.Path == statPath:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(statBody))
		case r.Method == http.MethodGet && r.URL.Path == issuesPath:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(issuesBody))
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)

	return srv
}
