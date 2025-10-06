package await

import (
	"context"
	"fmt"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanstage"
	"github.com/POSIdev-community/aictl/internal/core/port"
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

type UseCase struct {
	aiAdapter  port.Ai
	cliAdapter port.Cli
}

func NewUseCase(aiAdapter port.Ai, cliAdapter port.Cli) (*UseCase, error) {
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

			place := 0
			for i, id := range queue {
				if id == scanId {
					place = i
				}
			}

			u.cliAdapter.ShowText(fmt.Sprintf("%s: %d in queue", Enqueued, place))
		} else {
			u.cliAdapter.ShowText(fmt.Sprintf("%s: %d%%", stage.Stage, stage.Value))
		}

		time.Sleep(3 * time.Second)
	}

	if err != nil || !ScanComplete(stage) {
		return fmt.Errorf("scan stage %s in project %s", stage.Stage, cfg.ProjectId())
	}

	u.cliAdapter.ShowText(fmt.Sprintf("Scan '%s'", stage.Stage))

	return nil
}

func ScanComplete(stage scanstage.ScanStage) bool {
	return stage.Stage == Done || stage.Stage == Failed || stage.Stage == Aborted
}
