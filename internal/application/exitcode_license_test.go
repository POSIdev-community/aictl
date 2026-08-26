package application

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/apperror"
)

func TestMapExitCode_LicenseCheckFailure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		err      error
		wantCode int
	}{
		{
			name: "invalid_license",
			err: fmt.Errorf("initialize with retry: %w", fmt.Errorf("initialize ai adapter: %w",
				fmt.Errorf("license is invalid"))),
			wantCode: ExitCodeUnknown,
		},
		{
			name: "license_http_500",
			err: fmt.Errorf("initialize with retry: %w", fmt.Errorf("initialize ai adapter: %w",
				fmt.Errorf("ai check license: %w", apperror.NewServerError(http.StatusInternalServerError, "")))),
			wantCode: ExitCodeAPI,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			code, msg := mapExitCode(tt.err)
			require.NotEqual(t, ExitCodeSuccess, code)
			require.Equal(t, tt.wantCode, code)
			require.NotEmpty(t, msg)
		})
	}
}
