package await

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanstage"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
)

const (
	Enqueued = "Enqueued"
	Aborted  = "Aborted"
	Done     = "Done"
	Failed   = "Failed"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetScanStage(ctx context.Context, projectId uuid.UUID, scanId uuid.UUID) (scanstage.ScanStage, error)
	GetScanQueue(ctx context.Context) ([]uuid.UUID, error)
}

type CLI interface {
	ShowText(ctx context.Context, text string)
	ShowTextf(ctx context.Context, format string, a ...any)
}

type UseCase struct {
	aiAdapter  AI
	cliAdapter CLI
	cfg        *config.Config
}

func NewUseCase(aiAdapter AI, cliAdapter CLI, cfg *config.Config) (*UseCase, error) {
	if aiAdapter == nil {
		return nil, errs.NewValidationRequiredError("aiAdapter")
	}

	if cliAdapter == nil {
		return nil, errs.NewValidationRequiredError("cliAdapter")
	}

	return &UseCase{aiAdapter, cliAdapter, cfg}, nil
}

func (u *UseCase) Execute(ctx context.Context, scanId uuid.UUID) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize with jwt retry: %w", err)
	}

	u.cliAdapter.ShowTextf(ctx, "awating scan, id '%v'", scanId.String())

	failCount := 0
	stage := scanstage.ScanStage{}
	for failCount < 3 {
		stage, err = u.aiAdapter.GetScanStage(ctx, u.cfg.ProjectId(), scanId)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return err
			}

			failCount++
			time.Sleep(3 * time.Second)
			u.cliAdapter.ShowText(ctx, "...")
			continue
		}

		failCount = 0
		if ScanComplete(stage) {
			break
		}

		if stage.Stage == Enqueued {
			queue, err := u.aiAdapter.GetScanQueue(ctx)
			if err != nil {
				failCount++
				time.Sleep(3 * time.Second)
				u.cliAdapter.ShowText(ctx, "...")
				continue
			}

			place := 1
			for i, id := range queue {
				if id == scanId {
					place = i + 1
				}
			}

			u.cliAdapter.ShowTextf(ctx, "%s: %d/%d", strings.ToLower(stage.Stage), place, len(queue))
		} else {
			u.cliAdapter.ShowTextf(ctx, "%s: %d%%", strings.ToLower(stage.Stage), stage.Value)
		}

		time.Sleep(3 * time.Second)
	}

	if err != nil || !ScanComplete(stage) {
		return fmt.Errorf("scan stage %s in project %s", stage.Stage, u.cfg.ProjectId())
	}

	u.cliAdapter.ShowTextf(ctx, "Scan '%s'", stage.Stage)

	return nil
}

func ScanComplete(stage scanstage.ScanStage) bool {
	return stage.Stage == Done || stage.Stage == Failed || stage.Stage == Aborted
}
