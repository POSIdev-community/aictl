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

func TestNewScanAwaitCmd_ArgsAndFlags(t *testing.T) {
	validUUID := "123e4567-e89b-12d3-a456-426614174000"
	pid, _ := uuid.Parse(validUUID)
	sid, _ := uuid.Parse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")

	type args struct {
		scanId    uuid.UUID
		projectId uuid.UUID
	}

	tests := []struct {
		name      string
		cliArgs   []string
		flags     []string
		want      args
		expectErr bool
	}{
		{
			name:      "minimal with project-id and valid scan id",
			cliArgs:   []string{sid.String()},
			flags:     []string{"--project-id", validUUID},
			want:      args{scanId: sid, projectId: pid},
			expectErr: false,
		},
		{
			name:      "missing project-id → validation error",
			cliArgs:   []string{sid.String()},
			flags:     []string{},
			want:      args{scanId: sid},
			expectErr: true,
		},
		{
			name:      "invalid scan id (not uuid)",
			cliArgs:   []string{"not-a-uuid"},
			flags:     []string{"--project-id", validUUID},
			expectErr: true,
		},
		{
			name:      "invalid project-id",
			cliArgs:   []string{sid.String()},
			flags:     []string{"--project-id", "invalid"},
			expectErr: true,
		},
		{
			name:      "no args → ExactArgs(1) violation",
			cliArgs:   []string{},
			flags:     []string{"--project-id", validUUID},
			expectErr: true,
		},
		{
			name:      "too many args",
			cliArgs:   []string{"a", "b"},
			flags:     []string{"--project-id", validUUID},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewConfig(config.Uri{}, "", false, uuid.Nil, uuid.Nil)

			mockUC := new(mocks.UseCaseScanAwait)
			if !tt.expectErr {
				mockUC.On("Execute", mock.Anything, tt.want.scanId).Return(nil)
			}

			awaitCmd := NewScanAwaitCmd(cfg, mockUC).Command

			awaitCmd.SetArgs(append(tt.flags, tt.cliArgs...))
			awaitCmd.SetContext(context.Background())

			err := awaitCmd.Execute()
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				mockUC.AssertExpectations(t)
				require.Equal(t, tt.want.projectId, cfg.ProjectId())
			}
		})
	}
}
