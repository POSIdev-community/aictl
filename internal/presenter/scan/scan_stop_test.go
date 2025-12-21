package scan

import (
	"context"
	"testing"

	"github.com/POSIdev-community/aictl/internal/presenter/scan/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewScanStopCmd_ArgsAndFlags(t *testing.T) {
	sid, _ := uuid.Parse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")

	type args struct {
		scanId uuid.UUID
	}

	tests := []struct {
		name      string
		cliArgs   []string
		want      args
		expectErr bool
	}{
		{
			name:      "valid scan id",
			cliArgs:   []string{sid.String()},
			want:      args{scanId: sid},
			expectErr: false,
		},
		{
			name:      "invalid scan id (not uuid)",
			cliArgs:   []string{"not-a-uuid"},
			expectErr: true,
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
			mockUC := new(mocks.UseCaseScanStop)
			if !tt.expectErr {
				mockUC.On("Execute", mock.Anything, tt.want.scanId).Return(nil)
			}

			stopCmd := NewScanStopCmd(mockUC).Command

			stopCmd.SetArgs(tt.cliArgs)
			stopCmd.SetContext(context.Background())

			err := stopCmd.Execute()
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				mockUC.AssertExpectations(t)
			}
		})
	}
}
