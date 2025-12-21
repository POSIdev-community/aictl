package context

import (
	"context"
	"testing"

	"github.com/POSIdev-community/aictl/internal/presenter/context/mocks"
	"github.com/stretchr/testify/require"
)

func TestNewConfigUnsetCommand_UnsetFlags(t *testing.T) {
	type args struct {
		uri       bool
		token     bool
		tlsSkip   bool
		projectId bool
		branchId  bool
	}

	tests := []struct {
		name      string
		cliArgs   []string
		want      args
		expectErr bool
	}{
		{
			name:      "uri",
			cliArgs:   []string{"--uri"},
			want:      args{uri: true},
			expectErr: false,
		},
		{
			name:      "token",
			cliArgs:   []string{"--token"},
			want:      args{token: true},
			expectErr: false,
		},
		{
			name:      "tls-skip",
			cliArgs:   []string{"--tls-skip"},
			want:      args{tlsSkip: true},
			expectErr: false,
		},
		{
			name:      "project-id",
			cliArgs:   []string{"--project-id"},
			want:      args{projectId: true},
			expectErr: false,
		},
		{
			name:      "branch-id",
			cliArgs:   []string{"--branch-id"},
			want:      args{branchId: true},
			expectErr: false,
		},
		{
			name:      "short flags (-u -t -p -b)",
			cliArgs:   []string{"-u", "-t", "-p", "-b"},
			want:      args{uri: true, token: true, projectId: true, branchId: true},
			expectErr: false,
		},
		{
			name:      "all flags",
			cliArgs:   []string{"--uri", "--token", "--tls-skip", "--project-id", "--branch-id"},
			want:      args{uri: true, token: true, tlsSkip: true, projectId: true, branchId: true},
			expectErr: false,
		},
		{
			name:      "no flags",
			cliArgs:   []string{},
			want:      args{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := new(mocks.UseCaseConfigUnset)
			if !tt.expectErr {
				mockUC.On("Execute",
					tt.want.uri,
					tt.want.token,
					tt.want.tlsSkip,
					tt.want.projectId,
					tt.want.branchId,
				).Return(nil)
			}

			cmd := NewConfigUnsetCommand(mockUC)
			cmd.SetArgs(tt.cliArgs)
			cmd.SetContext(context.Background())

			err := cmd.Execute()
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				mockUC.AssertExpectations(t)
			}
		})
	}
}
