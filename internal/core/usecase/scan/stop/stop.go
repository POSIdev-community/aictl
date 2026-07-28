package stop

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	StopScan(ctx context.Context, scanResultId uuid.UUID) error
}

type UseCase struct {
	aiAdapter AI
}

func NewUseCase(aiAdapter AI) (*UseCase, error) {
	if aiAdapter == nil {
		return nil, validation.NewRequiredError("aiAdapter")
	}

	return &UseCase{aiAdapter}, nil
}

func (u *UseCase) Execute(ctx context.Context, scanResultId uuid.UUID) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	err = u.aiAdapter.StopScan(ctx, scanResultId)
	if err != nil {
		return err
	}

	return nil
}
