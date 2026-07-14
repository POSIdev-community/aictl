package await

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/queue"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanstage"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/google/uuid"
)

const (
	Enqueued = "Enqueued"
	Aborted  = "Aborted"
	Done     = "Done"
	Failed   = "Failed"

	pollInterval = 10 * time.Second
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetScanStage(ctx context.Context, projectId uuid.UUID, scanId uuid.UUID) (scanstage.ScanStage, error)
	GetScanItem(ctx context.Context, id uuid.UUID) (queue.Item, error)
	WatchScanStage(ctx context.Context, scanId uuid.UUID) (<-chan scanstage.ScanStage, error)
}

type CLI interface {
	ShowText(ctx context.Context, text string)
	ShowTextf(ctx context.Context, format string, a ...any)
	ReturnText(ctx context.Context, text string)
}

type UseCase struct {
	aiAdapter  AI
	cliAdapter CLI
	cfg        *config.Config
}

func NewUseCase(aiAdapter AI, cliAdapter CLI, cfg *config.Config) (*UseCase, error) {
	if aiAdapter == nil {
		return nil, validation.NewRequiredError("aiAdapter")
	}

	if cliAdapter == nil {
		return nil, validation.NewRequiredError("cliAdapter")
	}

	return &UseCase{aiAdapter, cliAdapter, cfg}, nil
}

func (u *UseCase) Execute(ctx context.Context, scanId uuid.UUID) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	u.cliAdapter.ShowTextf(ctx, "awating scan, id '%v'", scanId.String())

	stage, err := u.aiAdapter.GetScanStage(ctx, u.cfg.ProjectId(), scanId)
	if err != nil {
		return fmt.Errorf("get scan stage: %w", err)
	}
	u.showStage(ctx, scanId, stage)
	if scanComplete(stage) {
		u.finish(ctx, stage)

		return nil
	}

	updates, err := u.aiAdapter.WatchScanStage(ctx, scanId)
	if err != nil {
		return u.pollUntilDone(ctx, scanId)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case stage, ok := <-updates:
			if !ok {
				return u.pollUntilDone(ctx, scanId)
			}
			u.showStage(ctx, scanId, stage)
			if scanComplete(stage) {
				u.finish(ctx, stage)

				return nil
			}
		}
	}
}

func (u *UseCase) pollUntilDone(ctx context.Context, scanId uuid.UUID) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		stage, err := u.aiAdapter.GetScanStage(ctx, u.cfg.ProjectId(), scanId)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return err
			}
			u.cliAdapter.ShowTextf(ctx, "error getting scan stage: %v", err.Error())
			time.Sleep(pollInterval)
			u.cliAdapter.ShowText(ctx, "...")

			continue
		}

		u.showStage(ctx, scanId, stage)
		if scanComplete(stage) {
			u.finish(ctx, stage)

			return nil
		}

		time.Sleep(pollInterval)
	}
}

func (u *UseCase) showStage(ctx context.Context, scanId uuid.UUID, stage scanstage.ScanStage) {
	if stage.Stage == "" {
		return
	}

	if stage.Stage == Enqueued {
		item, queueErr := u.aiAdapter.GetScanItem(ctx, scanId)
		if queueErr != nil {
			if !errors.Is(queueErr, context.Canceled) {
				u.cliAdapter.ShowTextf(ctx, "error getting scan queue: %v", queueErr.Error())
			}
			u.cliAdapter.ShowTextf(ctx, "%s", strings.ToLower(stage.Stage))

			return
		}

		if item.OutOf > 0 {
			u.cliAdapter.ShowTextf(ctx, "%s: %d/%d", strings.ToLower(stage.Stage), item.Place, item.OutOf)
		} else {
			u.cliAdapter.ShowTextf(ctx, "%s", strings.ToLower(stage.Stage))
		}

		return
	}

	u.cliAdapter.ShowTextf(ctx, "%s: %d%%", strings.ToLower(stage.Stage), stage.Value)
}

func (u *UseCase) finish(ctx context.Context, stage scanstage.ScanStage) {
	u.cliAdapter.ShowTextf(ctx, "Scan '%s'", stage.Stage)
	u.cliAdapter.ReturnText(ctx, stage.Stage)
}

func scanComplete(stage scanstage.ScanStage) bool {
	return stage.Stage == Done || stage.Stage == Failed || stage.Stage == Aborted
}
