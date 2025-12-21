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

func TestNewScanStartBranchCmd_ArgsAndFlags(t *testing.T) {
	validBranchID := "123e4567-e89b-12d3-a456-426614174000"
	bid, _ := uuid.Parse(validBranchID)

	tests := []struct {
		name       string
		cliArgs    []string
		wantBranch uuid.UUID
		expectErr  bool
	}{
		{
			name:       "valid branch id",
			cliArgs:    []string{validBranchID},
			wantBranch: bid,
			expectErr:  false,
		},
		{
			name:       "empty branch id → validation error",
			cliArgs:    []string{""},
			wantBranch: uuid.Nil,
			expectErr:  true,
		},
		{
			name:       "invalid branch id (not uuid)",
			cliArgs:    []string{"not-a-uuid"},
			wantBranch: uuid.Nil,
			expectErr:  true,
		},
		{
			name:      "no args → ExactArgs(1) violation",
			cliArgs:   []string{},
			expectErr: true,
		},
		{
			name:      "too many args",
			cliArgs:   []string{"a", "b"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewConfig(config.Uri{}, "", false, uuid.Nil, uuid.Nil)

			mockUC := new(mocks.UseCaseScanStartBranch)
			if !tt.expectErr {
				mockUC.On("Execute", mock.Anything).Return(nil)
			}

			startCmd := NewScanStartBranchCmd(cfg, mockUC).Command

			startCmd.SetArgs(tt.cliArgs)
			startCmd.SetContext(context.Background())

			err := startCmd.Execute()
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				mockUC.AssertExpectations(t)
				require.Equal(t, tt.wantBranch, cfg.BranchId())
			}
		})
	}
}
