package scan

import (
	"context"
	"testing"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/scan/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewScanStartProjectCmd_ArgsAndFlags(t *testing.T) {
	validProjectID := "123e4567-e89b-12d3-a456-426614174000"
	pid, _ := uuid.Parse(validProjectID)

	tests := []struct {
		name        string
		cliArgs     []string
		wantProject uuid.UUID
		expectErr   bool
	}{
		{
			name:        "valid project id",
			cliArgs:     []string{validProjectID},
			wantProject: pid,
			expectErr:   false,
		},
		{
			name:      "empty project id → validation error",
			cliArgs:   []string{""},
			expectErr: true,
		},
		{
			name:      "invalid project id (not uuid)",
			cliArgs:   []string{"not-a-uuid"},
			expectErr: true,
		},
		{
			name:      "too many args → MaximumNArgs(1) violation",
			cliArgs:   []string{"a", "b"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewConfig(config.Uri{}, "", false, uuid.Nil, uuid.Nil)

			mockUC := new(mocks.UseCaseScanStartProject)
			if !tt.expectErr {
				mockUC.On("Execute", mock.Anything).Return(nil)
			}

			startCmd := NewScanStartProjectCmd(cfg, mockUC).Command

			startCmd.SetArgs(tt.cliArgs)
			startCmd.SetContext(context.Background())

			err := startCmd.Execute()
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				mockUC.AssertExpectations(t)
				require.Equal(t, tt.wantProject, cfg.ProjectId())
			}
		})
	}
}
