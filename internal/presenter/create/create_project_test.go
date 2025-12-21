package create

import (
	"context"
	"testing"

	"github.com/POSIdev-community/aictl/internal/presenter/create/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewCreateProjectCmd_ArgsAndFlags(t *testing.T) {
	type args struct {
		projectName string
		safe        bool
	}

	tests := []struct {
		name      string
		cliArgs   []string
		flags     []string
		want      args
		expectErr bool
	}{
		{
			name:      "minimal: project name only",
			cliArgs:   []string{"my-project"},
			flags:     nil,
			want:      args{projectName: "my-project"},
			expectErr: false,
		},
		{
			name:      "with safe flag",
			cliArgs:   []string{"prod"},
			flags:     []string{"--safe"},
			want:      args{projectName: "prod", safe: true},
			expectErr: false,
		},
		{
			name:      "no args → ExactArgs(1) violation",
			cliArgs:   []string{},
			flags:     []string{},
			expectErr: true,
		},
		{
			name:      "too many args",
			cliArgs:   []string{"a", "b"},
			flags:     nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := new(mocks.UseCaseCreateProject)
			if !tt.expectErr {
				mockUC.On("Execute",
					mock.Anything,
					tt.want.projectName,
					tt.want.safe,
				).Return(nil)
			}

			projectCmd := NewCreateProjectCmd(mockUC).Command

			// Привязываем safe как persistent-флаг (как в NewCreateCmd)
			projectCmd.Flags().BoolVar(&safeFlag, "safe", false, "if resource exists, return its id without error")

			projectCmd.SetArgs(append(tt.flags, tt.cliArgs...))
			projectCmd.SetContext(context.Background())

			err := projectCmd.Execute()
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				mockUC.AssertExpectations(t)
			}
		})
	}
}
