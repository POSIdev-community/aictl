package utils

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	domainlicense "github.com/POSIdev-community/aictl/internal/core/domain/license"
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
)

// LicenseAI is the port used to apply license constraints before scan start.
type LicenseAI interface {
	GetLicense(ctx context.Context) (*domainlicense.License, error)
	GetProjectSettings(ctx context.Context, projectId uuid.UUID) (settings.ScanSettings, error)
	SetProjectSettings(ctx context.Context, projectId uuid.UUID, settings *settings.ScanSettings) error
}

// LicenseCLI reports warnings about disabled modules.
type LicenseCLI interface {
	ShowTextf(ctx context.Context, format string, a ...any)
}

// ApplyLicenseOptions controls license application before scan start.
type ApplyLicenseOptions struct {
	SkipLanguageCheck bool
}

// ApplyLicenseToProjectSettings checks languages, disables unlicensed modules, and persists when safe.
func ApplyLicenseToProjectSettings(
	ctx context.Context,
	ai LicenseAI,
	cli LicenseCLI,
	projectId uuid.UUID,
	opts ApplyLicenseOptions,
) error {
	lic, err := ai.GetLicense(ctx)
	if err != nil {
		return fmt.Errorf("get license: %w", err)
	}

	current, err := ai.GetProjectSettings(ctx, projectId)
	if err != nil {
		return fmt.Errorf("get project settings: %w", err)
	}

	updated, removed, err := domainlicense.PrepareSettings(lic, current, opts.SkipLanguageCheck)
	if err != nil {
		return err
	}

	if len(removed) == 0 {
		return nil
	}

	if err := ai.SetProjectSettings(ctx, projectId, &updated); err != nil {
		return fmt.Errorf("set project settings: %w", err)
	}

	for _, module := range removed {
		cli.ShowTextf(ctx, "license: disabled unlicensed scan module '%s'", module)
	}

	return nil
}
