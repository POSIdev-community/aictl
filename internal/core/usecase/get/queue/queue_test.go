package queue

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	domainqueue "github.com/POSIdev-community/aictl/internal/core/domain/queue"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanstage"
)

func TestUseCase_Execute_All(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	entries := []domainqueue.Entry{
		{ProjectId: uuid.New(), BranchId: uuid.New(), ScanId: uuid.New(), Stage: scanstage.Enqueued},
	}

	ai := NewMockAI(t)
	ai.On("InitializeWithRetry", ctx).Return(nil).Once()
	ai.On("GetScanQueue", ctx).Return(entries, nil).Once()

	cli := NewMockCLI(t)
	cli.On("ShowQueue", ctx, entries).Return().Once()

	uc, err := NewUseCase(ai, cli)
	require.NoError(t, err)
	require.NoError(t, uc.Execute(ctx, nil))
}

func TestUseCase_Execute_FilterByProject(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	projectA := uuid.New()
	projectB := uuid.New()
	entries := []domainqueue.Entry{
		{ProjectId: projectA, BranchId: uuid.New(), ScanId: uuid.New(), Stage: scanstage.Enqueued},
		{ProjectId: projectB, BranchId: uuid.New(), ScanId: uuid.New(), Stage: scanstage.Enqueued},
	}

	ai := NewMockAI(t)
	ai.On("InitializeWithRetry", ctx).Return(nil).Once()
	ai.On("GetScanQueue", ctx).Return(entries, nil).Once()

	cli := NewMockCLI(t)
	cli.On("ShowQueue", ctx, []domainqueue.Entry{entries[0]}).Return().Once()

	uc, err := NewUseCase(ai, cli)
	require.NoError(t, err)
	require.NoError(t, uc.Execute(ctx, &projectA))
}
