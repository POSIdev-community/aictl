package context

import (
	"context"
	"testing"

	"github.com/POSIdev-community/aictl/internal/presenter/context/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewConfigShowCommand_FormatFlags(t *testing.T) {
	type args struct {
		json bool
		yaml bool
	}

	tests := []struct {
		name      string
		cliArgs   []string
		want      args
		expectErr bool
	}{
		{
			name:      "no flags",
			cliArgs:   []string{},
			want:      args{},
			expectErr: false,
		},
		{
			name:      "json",
			cliArgs:   []string{"--json"},
			want:      args{json: true},
			expectErr: false,
		},
		{
			name:      "yaml",
			cliArgs:   []string{"--yaml"},
			want:      args{yaml: true},
			expectErr: false,
		},
		{
			name:      "conflict: json and yaml",
			cliArgs:   []string{"--json", "--yaml"},
			want:      args{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := new(mocks.UseCaseConfigShow)
			if !tt.expectErr {
				mockUC.On("Execute", mock.Anything, tt.want.json, tt.want.yaml).Return(nil)
			}

			cmd := NewConfigShowCommand(mockUC)
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
