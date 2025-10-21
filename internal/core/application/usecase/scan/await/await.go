package await

import (
	"context"
	"fmt"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanstage"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
	"time"
)

const (
	Enqueued = "Enqueued"
	Aborted  = "Aborted"
	Done     = "Done"
	Failed   = "Failed"
)

type AI interface {
	GetScanStage(ctx context.Context, projectId uuid.UUID, scanId uuid.UUID) (scanstage.ScanStage, error)
	GetScanQueue(ctx context.Context) ([]uuid.UUID, error)
}

type CLI interface {
	ShowText(text string)
	ShowTextF(format string, a ...any)
}

type UseCase struct {
	aiAdapter  AI
	cliAdapter CLI
}

func NewUseCase(aiAdapter AI, cliAdapter CLI) (*UseCase, error) {
	if aiAdapter == nil {
		return nil, errs.NewValidationRequiredError("aiAdapter")
	}

	if cliAdapter == nil {
		return nil, errs.NewValidationRequiredError("cliAdapter")
	}

	return &UseCase{aiAdapter, cliAdapter}, nil
}

func (u *UseCase) Execute(ctx context.Context, cfg *config.Config, scanId uuid.UUID) error {
	failCount := 0
	stage := scanstage.ScanStage{}
	var err error
	for failCount < 3 {
		stage, err = u.aiAdapter.GetScanStage(ctx, cfg.ProjectId(), scanId)
		if err != nil {
			failCount++
			time.Sleep(3 * time.Second)
			u.cliAdapter.ShowText("...")
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
				u.cliAdapter.ShowText("...")
				continue
			}

			place := 1
			for i, id := range queue {
				if id == scanId {
					place = i + 1
				}
			}

			u.cliAdapter.ShowTextF("%s: %d/%d", Enqueued, place, len(queue))
		} else {
			u.cliAdapter.ShowTextF("%s: %d%%", stage.Stage, stage.Value)
		}

		time.Sleep(3 * time.Second)
	}

	if err != nil || !ScanComplete(stage) {
		return fmt.Errorf("scan stage %s in project %s", stage.Stage, cfg.ProjectId())
	}

	u.cliAdapter.ShowTextF("Scan '%s'", stage.Stage)

	return nil
}

func ScanComplete(stage scanstage.ScanStage) bool {
	return stage.Stage == Done || stage.Stage == Failed || stage.Stage == Aborted
}
