package create

import (
	"context"
	"testing"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	_utils "github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/POSIdev-community/aictl/internal/presenter/create/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewCreateBranchCmd_ArgsAndFlags(t *testing.T) {
	validUUID := "123e4567-e89b-12d3-a456-426614174000"
	pid, _ := uuid.Parse(validUUID)

	type args struct {
		branchName string
		scanTarget string
		safe       bool
	}

	tests := []struct {
		name      string
		cliArgs   []string
		flags     []string
		want      args
		expectErr bool
	}{
		{
			name:      "minimal with project-id",
			cliArgs:   []string{"main"},
			flags:     []string{"--project-id", validUUID},
			want:      args{branchName: "main"},
			expectErr: false,
		},
		{
			name:      "with scan-target",
			cliArgs:   []string{"dev"},
			flags:     []string{"--project-id", validUUID, "--scan-target", "/tmp"},
			want:      args{branchName: "dev", scanTarget: "/tmp"},
			expectErr: false,
		},
		{
			name:      "with safe flag",
			cliArgs:   []string{"fix"},
			flags:     []string{"--project-id", validUUID, "--safe"},
			want:      args{branchName: "fix", safe: true},
			expectErr: false,
		},
		{
			name:      "all flags",
			cliArgs:   []string{"release"},
			flags:     []string{"--project-id", validUUID, "--scan-target", ".", "--safe"},
			want:      args{branchName: "release", scanTarget: ".", safe: true},
			expectErr: false,
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
		{
			name:      "invalid project-id",
			cliArgs:   []string{"x"},
			flags:     []string{"--project-id", "not-a-uuid"},
			expectErr: true,
		},
		{
			name:      "missing project-id → validation error",
			cliArgs:   []string{"x"},
			flags:     []string{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewConfig(config.Uri{}, "", false, uuid.Nil, uuid.Nil)

			mockUC := new(mocks.UseCaseCreateBranch)
			if !tt.expectErr {
				mockUC.On("Execute",
					mock.Anything,
					tt.want.branchName,
					tt.want.scanTarget,
					tt.want.safe,
				).Return(nil)
			}

			branchCmd := NewCreateBranchCmd(cfg, mockUC).Command

			branchCmd.Flags().BoolVar(&safeFlag, "safe", false, "")
			_utils.AddConnectionPersistentFlags(branchCmd)

			branchCmd.SetArgs(append(tt.flags, tt.cliArgs...))
			branchCmd.SetContext(context.Background())

			err := branchCmd.Execute()
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				mockUC.AssertExpectations(t)
				require.Equal(t, pid, cfg.ProjectId())
			}
		})
	}
}
