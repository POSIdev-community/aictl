package context

import (
	"testing"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/context/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigSetCommand_FlagHandling(t *testing.T) {
	validUUID := "123e4567-e89b-12d3-a456-426614174000"
	emptyUUID := uuid.Nil

	tests := []struct {
		name      string
		args      []string
		assertCfg func(t *testing.T, cfg *config.Config)
		expectErr bool
	}{
		{
			name: "uri flag",
			args: []string{"--uri", "https://example.com"},
			assertCfg: func(t *testing.T, cfg *config.Config) {
				assert.Equal(t, "https://example.com", cfg.UriString())
			},
		},
		{
			name: "token flag",
			args: []string{"--token", "secret123"},
			assertCfg: func(t *testing.T, cfg *config.Config) {
				assert.Equal(t, "secret123", cfg.Token())
			},
		},
		{
			name: "tls-skip flag",
			args: []string{"--tls-skip"},
			assertCfg: func(t *testing.T, cfg *config.Config) {
				assert.True(t, cfg.TLSSkip())
			},
		},
		{
			name: "no-tls-skip flag",
			args: []string{"--no-tls-skip"},
			assertCfg: func(t *testing.T, cfg *config.Config) {
				assert.False(t, cfg.TLSSkip())
			},
		},
		{
			name: "project-id flag (valid)",
			args: []string{"--project-id", validUUID},
			assertCfg: func(t *testing.T, cfg *config.Config) {
				expected, _ := uuid.Parse(validUUID)
				assert.Equal(t, expected, cfg.ProjectId())
			},
		},
		{
			name: "branch-id flag (valid)",
			args: []string{"--branch-id", validUUID},
			assertCfg: func(t *testing.T, cfg *config.Config) {
				expected, _ := uuid.Parse(validUUID)
				assert.Equal(t, expected, cfg.BranchId())
			},
		},
		{
			name: "multiple flags",
			args: []string{"-u", "http://a", "-t", "t", "--tls-skip", "-p", validUUID, "-b", validUUID},
			assertCfg: func(t *testing.T, cfg *config.Config) {
				assert.Equal(t, "http://a", cfg.UriString())
				assert.Equal(t, "t", cfg.Token())
				assert.True(t, cfg.TLSSkip())
				pid, _ := uuid.Parse(validUUID)
				bid, _ := uuid.Parse(validUUID)
				assert.Equal(t, pid, cfg.ProjectId())
				assert.Equal(t, bid, cfg.BranchId())
			},
		},
		{
			name:      "conflicting tls flags",
			args:      []string{"--tls-skip", "--no-tls-skip"},
			expectErr: true,
		},
		{
			name:      "no flags → validation error",
			args:      []string{},
			expectErr: true,
		},
		{
			name:      "project-id flag (invalid UUID)",
			args:      []string{"--project-id", "not-a-uuid"},
			expectErr: true,
		},
		{
			name:      "branch-id flag (empty)",
			args:      []string{"--branch-id", ""},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewConfig(
				config.Uri{}, "", false, emptyUUID, emptyUUID,
			)

			mockUC := new(mocks.UseCaseConfigSet)
			if !tt.expectErr {
				mockUC.On("Execute").Return(nil)
			}

			cmd := NewConfigSetCommand(cfg, mockUC)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				tt.assertCfg(t, cfg)
				mockUC.AssertExpectations(t)
			}
		})
	}
}
