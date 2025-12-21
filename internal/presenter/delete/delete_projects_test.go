package delete

import (
	"context"
	"testing"

	"github.com/POSIdev-community/aictl/internal/presenter/delete/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewDeleteProjectsCommand_UUIDParsing(t *testing.T) {
	valid1 := "123e4567-e89b-12d3-a456-426614174000"
	valid2 := "123e4567-e89b-12d3-a456-426614174001"
	uuid1, _ := uuid.Parse(valid1)
	uuid2, _ := uuid.Parse(valid2)

	type args struct {
		projectIds []uuid.UUID
	}

	tests := []struct {
		name      string
		cliArgs   []string
		want      args
		expectErr bool
	}{
		{
			name:      "single valid UUID",
			cliArgs:   []string{valid1},
			want:      args{projectIds: []uuid.UUID{uuid1}},
			expectErr: false,
		},
		{
			name:      "multiple valid UUIDs",
			cliArgs:   []string{valid1, valid2},
			want:      args{projectIds: []uuid.UUID{uuid1, uuid2}},
			expectErr: false,
		},
		{
			name:      "invalid UUID",
			cliArgs:   []string{"not-a-uuid"},
			want:      args{},
			expectErr: true,
		},
		{
			name:      "mixed valid and invalid",
			cliArgs:   []string{valid1, "invalid"},
			want:      args{},
			expectErr: true,
		},
		{
			name:      "empty args",
			cliArgs:   []string{},
			want:      args{},
			expectErr: true, // MinimumNArgs(1)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := new(mocks.UseCaseDeleteProjects)
			if !tt.expectErr {
				mockUC.On("Execute", mock.Anything, tt.want.projectIds).Return(nil)
			}

			cmd := NewDeleteProjectsCommand(mockUC)
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
