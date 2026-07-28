package await

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/apperror"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanstage"
)

type run struct {
	*UseCase
	scanId           uuid.UUID
	failOnScanFailed bool
}

func (r *run) waitUntilDone(ctx context.Context, updates <-chan scanstage.ScanStage) error {
	timer := time.NewTimer(r.pollInterval)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case stage, ok := <-updates:
			if !ok {
				updates = nil
				done, err := r.onPollTick(ctx, timer, true)
				if err != nil {
					return err
				}
				if done {
					return nil
				}

				continue
			}
			drainAndReset(timer, r.pollInterval)

			r.showStage(ctx, stage)
			if stage.IsCompleted() {
				return r.finish(ctx, stage)
			}
		case <-timer.C:
			done, err := r.onPollTick(ctx, timer, updates == nil)
			if err != nil {
				return err
			}
			if done {
				return nil
			}
		}
	}
}

func (r *run) onPollTick(ctx context.Context, timer *time.Timer, showDots bool) (bool, error) {
	done, err := r.checkStage(ctx)
	if err != nil {
		if isTerminalAwaitErr(err) {
			return false, err
		}
		r.cliAdapter.ShowTextf(ctx, "error getting scan stage: %v", err.Error())
		if showDots {
			r.cliAdapter.ShowText(ctx, "...")
		}
		timer.Reset(r.pollInterval)

		return false, nil
	}
	if done {
		return true, nil
	}
	timer.Reset(r.pollInterval)

	return false, nil
}

// checkStage fetches current scan stage, prints it, and finishes if complete.
// Returns (true, nil) when the scan is done successfully, or (true, FailError) when
// fail-on-scan-failed triggers.
func (r *run) checkStage(ctx context.Context) (bool, error) {
	stage, err := r.aiAdapter.GetScanStage(ctx, r.cfg.ProjectId(), r.scanId)
	if err != nil {
		return false, err
	}

	r.showStage(ctx, stage)
	if stage.IsCompleted() {
		return true, r.finish(ctx, stage)
	}

	return false, nil
}

func (r *run) showStage(ctx context.Context, stage scanstage.ScanStage) {
	if stage.Stage == "" {
		return
	}

	if stage.IsQueued() {
		item, queueErr := r.aiAdapter.GetScanItem(ctx, r.scanId)
		if queueErr != nil {
			if !errors.Is(queueErr, context.Canceled) {
				r.cliAdapter.ShowTextf(ctx, "error getting scan queue: %v", queueErr.Error())
			}
			r.cliAdapter.ShowTextf(ctx, "%s", strings.ToLower(stage.Stage))

			return
		}

		if item.OutOf > 0 {
			r.cliAdapter.ShowTextf(ctx, "%s: %d/%d", strings.ToLower(stage.Stage), item.Place, item.OutOf)
		} else {
			r.cliAdapter.ShowTextf(ctx, "%s", strings.ToLower(stage.Stage))
		}

		return
	}

	r.cliAdapter.ShowTextf(ctx, "%s: %d%%", strings.ToLower(stage.Stage), stage.Value)
}

func (r *run) finish(ctx context.Context, stage scanstage.ScanStage) error {
	r.cliAdapter.ShowTextf(ctx, "Scan '%s'", stage.Stage)
	r.cliAdapter.ReturnText(ctx, stage.Stage)

	if r.failOnScanFailed && stage.IsFailed() {
		return apperror.NewScanFailedError(stage.Stage)
	}

	return nil
}

func drainAndReset(timer *time.Timer, d time.Duration) {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
	timer.Reset(d)
}
