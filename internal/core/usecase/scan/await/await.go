package await

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/apperror"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/queue"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanstage"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

const (
	DefaultPollInterval    = 10 * time.Second
	DefaultNotFoundRetries = 5
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetScanStage(ctx context.Context, projectId uuid.UUID, scanId uuid.UUID) (scanstage.ScanStage, error)
	GetScanItem(ctx context.Context, id uuid.UUID) (queue.Item, error)
	WatchScanStage(ctx context.Context, scanId uuid.UUID) (<-chan scanstage.ScanStage, <-chan error, error)
}

type CLI interface {
	ShowText(ctx context.Context, text string)
	ShowTextf(ctx context.Context, format string, a ...any)
	ReturnText(ctx context.Context, text string)
}

type UseCase struct {
	aiAdapter          AI
	cliAdapter         CLI
	cfg                *config.Config
	pollInterval       time.Duration
	notFoundMaxRetries int
}

func NewUseCase(aiAdapter AI, cliAdapter CLI, cfg *config.Config, pollInterval time.Duration) (*UseCase, error) {
	if aiAdapter == nil {
		return nil, validation.NewRequiredError("aiAdapter")
	}

	if cliAdapter == nil {
		return nil, validation.NewRequiredError("cliAdapter")
	}

	if pollInterval <= 0 {
		return nil, validation.NewInvalidError("pollInterval")
	}

	return &UseCase{
		aiAdapter:          aiAdapter,
		cliAdapter:         cliAdapter,
		cfg:                cfg,
		pollInterval:       pollInterval,
		notFoundMaxRetries: DefaultNotFoundRetries,
	}, nil
}

func (u *UseCase) Execute(ctx context.Context, scanId uuid.UUID, failOnScanFailed bool) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	r := &run{UseCase: u, scanId: scanId, failOnScanFailed: failOnScanFailed}
	u.cliAdapter.ShowTextf(ctx, "awaiting scan, id '%v'", scanId.String())

	done, err := r.checkStage(ctx)
	if err != nil {
		if fatal, fatalErr := r.classifyPollError(err); fatal {
			return fatalErr
		}
		u.cliAdapter.ShowTextf(ctx, "error getting scan stage: %v", err.Error())
	} else if done {
		return nil
	}

	updates, watchErrs, err := u.aiAdapter.WatchScanStage(ctx, scanId)
	if err != nil {
		updates = nil
		watchErrs = nil
	}

	return r.waitUntilDone(ctx, updates, watchErrs)
}

func isImmediateTerminalAwaitErr(err error) bool {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	var (
		failErr  *apperror.FailError
		authErr  *apperror.AuthenticationError
		authzErr *apperror.AuthorizationError
	)

	return errors.As(err, &failErr) || errors.As(err, &authErr) || errors.As(err, &authzErr)
}

func isNotFoundErr(err error) bool {
	var (
		notFound     *apperror.NotFoundError
		notFoundByID *apperror.NotFoundByIdError
	)

	return errors.As(err, &notFound) || errors.As(err, &notFoundByID)
}
