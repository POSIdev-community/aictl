package scanning

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	domainqueue "github.com/POSIdev-community/aictl/internal/core/domain/queue"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetActiveScans(ctx context.Context) ([]domainqueue.Entry, error)
}

type CLI interface {
	ShowQueue(ctx context.Context, entries []domainqueue.Entry)
}

type UseCase struct {
	aiAdapter  AI
	cliAdapter CLI
}

func NewUseCase(aiAdapter AI, cliAdapter CLI) (*UseCase, error) {
	if aiAdapter == nil {
		return nil, validation.NewRequiredError("aiAdapter")
	}

	if cliAdapter == nil {
		return nil, validation.NewRequiredError("cliAdapter")
	}

	return &UseCase{aiAdapter, cliAdapter}, nil
}

// Execute lists active scans. projectFilter is opt-in via -p only; ctx project is ignored.
func (u *UseCase) Execute(ctx context.Context, projectFilter *uuid.UUID) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	entries, err := u.aiAdapter.GetActiveScans(ctx)
	if err != nil {
		return fmt.Errorf("get active scans: %w", err)
	}

	if projectFilter != nil {
		filtered := make([]domainqueue.Entry, 0, len(entries))
		for _, e := range entries {
			if e.ProjectId == *projectFilter {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}

	u.cliAdapter.ShowQueue(ctx, entries)

	return nil
}
