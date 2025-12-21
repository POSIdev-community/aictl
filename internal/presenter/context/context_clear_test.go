package context

import (
	"testing"

	"github.com/POSIdev-community/aictl/internal/presenter/context/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewConfigClearCommand_YesFlag(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{"long flag", []string{"--yes"}, true},
		{"short flag", []string{"-y"}, true},
		{"no flag", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := new(mocks.UseCaseConfigClear)
			mockUC.On("Execute", mock.Anything, tt.expected).Return(nil)

			cmd := NewConfigClearCommand(mockUC)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			require.NoError(t, err)

			mockUC.AssertExpectations(t)
		})
	}
}
