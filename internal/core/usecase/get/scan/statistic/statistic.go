package statistic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/statistic"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetScanStatistic(ctx context.Context, projectId, scanResultId uuid.UUID) (*statistic.Statistic, error)
}

type CLI interface {
	ShowReader(io.Reader) error
	ShowTextf(ctx context.Context, format string, args ...any)
	ShowScanStatistic(context.Context, *statistic.Statistic)
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

func (u *UseCase) Execute(ctx context.Context, scanId uuid.UUID, fullDestPath string, jsonOutput bool) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	u.cliAdapter.ShowTextf(ctx, "getting scan statistic, scan-id '%v'", scanId.String())

	stat, err := u.aiAdapter.GetScanStatistic(ctx, u.cfg.ProjectId(), scanId)
	if err != nil {
		return fmt.Errorf("get scan sbom: %w", err)
	}

	if fullDestPath == "" && !jsonOutput {
		u.cliAdapter.ShowScanStatistic(ctx, stat)

		return nil
	}

	jsonData, err := json.MarshalIndent(stat, "", "    ")
	if err != nil {
		return fmt.Errorf("marshal scan statistic: %w", err)
	}
	jsonData = append(jsonData, '\n')

	if fullDestPath != "" {
		if err = os.WriteFile(fullDestPath, jsonData, 0644); err != nil {
			return fmt.Errorf("write statistic to path %s: %w", fullDestPath, err)
		}

		return nil
	}

	reader := bytes.NewReader(jsonData)
	readCloser := io.NopCloser(reader)
	defer func() {
		_ = readCloser.Close()
	}()

	if err := u.cliAdapter.ShowReader(reader); err != nil {
		return fmt.Errorf("print statistic: %w", err)
	}

	return nil
}
