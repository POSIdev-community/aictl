package exclusions

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
)

func TestUseCase_Execute(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	projectID := uuid.New()
	cfg := config.NewConfig(config.Uri{}, "", true, projectID, uuid.New())

	ai := NewMockAI(t)
	ai.On("InitializeWithRetry", ctx).Return(nil).Once()
	ai.On("GetProjectExclusions", ctx, projectID).Return("*.log\nnode_modules/", nil).Once()

	cli := NewMockCLI(t)
	cli.On("ReturnText", ctx, "*.log\nnode_modules/").Return().Once()

	uc, err := NewUseCase(ai, cli, cfg)
	require.NoError(t, err)
	require.NoError(t, uc.Execute(ctx))
}
