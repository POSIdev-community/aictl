package scans

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/regexfilter"
	"github.com/POSIdev-community/aictl/internal/core/domain/scan"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetScans(ctx context.Context, branchId uuid.UUID) ([]scan.Scan, error)
	GetLastScan(ctx context.Context, branchId uuid.UUID) (*scan.Scan, error)
}

type CLI interface {
	ShowScans(ctx context.Context, scans []scan.Scan)
	ShowScansQuite(ctx context.Context, scans []scan.Scan)
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

func (u *UseCase) Execute(ctx context.Context, filter regexfilter.RegexFilter, quite bool, latest bool) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	var scans []scan.Scan
	if latest {
		lastScan, err := u.aiAdapter.GetLastScan(ctx, u.cfg.BranchId())
		if err != nil {
			return fmt.Errorf("get last scan: %w", err)
		}

		scans = []scan.Scan{*lastScan}
	} else {
		scans, err = u.aiAdapter.GetScans(ctx, u.cfg.BranchId())
		if err != nil {
			return fmt.Errorf("get scans: %w", err)
		}
	}

	filteredScans := make([]scan.Scan, 0, len(scans))
	if filter.Empty() {
		filteredScans = scans
	} else {
		for _, s := range scans {
			if filter.Execute(s.Id.String()) || filter.Execute(s.ScanLabel) {
				filteredScans = append(filteredScans, s)
			}
		}
	}

	if quite {
		u.cliAdapter.ShowScansQuite(ctx, filteredScans)
	} else {
		u.cliAdapter.ShowScans(ctx, filteredScans)
	}

	return nil
}
