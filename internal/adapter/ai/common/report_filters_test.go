package common_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/POSIdev-community/aictl/internal/core/domain/version"
)

func TestValidateReportFiltersForVersion(t *testing.T) {
	t.Parallel()

	v5, err := version.NewVersion("5.4.0")
	require.NoError(t, err)
	v60, err := version.NewVersion("6.0.0")
	require.NoError(t, err)
	v61, err := version.NewVersion("6.1.0")
	require.NoError(t, err)
	v63, err := version.NewVersion("6.3.0")
	require.NoError(t, err)

	cases := []struct {
		name    string
		ver     version.Version
		filters report.Filters
		wantErr string
	}{
		{
			name:    "v5x rejects SecretDetection",
			ver:     v5,
			filters: report.Filters{Apply: true, ScanModules: []string{"SecretDetection"}},
			wantErr: ">= 6.0.0",
		},
		{
			name:    "v5x rejects MaliciousCodeDetection",
			ver:     v5,
			filters: report.Filters{Apply: true, ScanModules: []string{"MaliciousCodeDetection"}},
			wantErr: ">= 6.0.0",
		},
		{
			name:    "v5x rejects OneC",
			ver:     v5,
			filters: report.Filters{Apply: true, Languages: []string{"OneC"}},
			wantErr: ">= 6.0.0",
		},
		{
			name:    "v5x rejects Dart",
			ver:     v5,
			filters: report.Filters{Apply: true, Languages: []string{"Dart"}},
			wantErr: ">= 6.1.0",
		},
		{
			name:    "v60 accepts SecretDetection rejects Dart",
			ver:     v60,
			filters: report.Filters{Apply: true, ScanModules: []string{"SecretDetection"}, Languages: []string{"Dart"}},
			wantErr: ">= 6.1.0",
		},
		{
			name:    "v60 accepts OneC",
			ver:     v60,
			filters: report.Filters{Apply: true, Languages: []string{"OneC"}, ScanModules: []string{"SecretDetection"}},
		},
		{
			name:    "v61 accepts Dart and SecretDetection",
			ver:     v61,
			filters: report.Filters{Apply: true, Languages: []string{"Dart"}, ScanModules: []string{"SecretDetection"}},
		},
		{
			name:    "v63 accepts all",
			ver:     v63,
			filters: report.Filters{Apply: true, Languages: []string{"Dart", "OneC"}, ScanModules: []string{"SecretDetection", "MaliciousCodeDetection"}},
		},
		{
			name:    "empty filters skip checks",
			ver:     v5,
			filters: report.EmptyFilters(),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := common.ValidateReportFiltersForVersion(tc.filters, tc.ver)
			if tc.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
