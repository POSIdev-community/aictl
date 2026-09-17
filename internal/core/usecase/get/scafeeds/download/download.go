package download

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/POSIdev-community/aictl/internal/core/domain/scafeeds"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	utils "github.com/POSIdev-community/aictl/internal/core/usecase/.utils"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	DownloadScaFeeds(ctx context.Context, version string) (io.ReadCloser, string, error)
}

type CLI interface {
	ShowTextf(ctx context.Context, format string, a ...any)
}

type UseCase struct {
	aiAdapter  AI
	cliAdapter CLI
}

func NewUseCase(aiAdapter AI, cliAdapter CLI) (*UseCase, error) {
	if aiAdapter == nil {
		return nil, validation.NewRequiredError("aiAdapter")
	}
	if cliAdapter == nil {
		return nil, validation.NewRequiredError("cliAdapter")
	}

	return &UseCase{aiAdapter, cliAdapter}, nil
}

func (u *UseCase) Execute(ctx context.Context, version, outPath string) error {
	if err := u.aiAdapter.InitializeWithRetry(ctx); err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	body, fileName, err := u.aiAdapter.DownloadScaFeeds(ctx, version)
	if err != nil {
		return fmt.Errorf("download sca feeds: %w", err)
	}
	defer func() { _ = body.Close() }()

	if fileName == "" {
		fileName = scafeeds.DefaultArchiveName(version)
	}

	dest, err := resolveDownloadPath(outPath, fileName)
	if err != nil {
		return err
	}

	if err := utils.CopyFileToPath(body, dest); err != nil {
		return fmt.Errorf("write sca feeds to %s: %w", dest, err)
	}

	u.cliAdapter.ShowTextf(ctx, "saved sca feeds to '%s'", dest)

	return nil
}

func resolveDownloadPath(outPath, fileName string) (string, error) {
	if outPath == "" {
		dest := fileName
		if fshelper.PathExists(dest) && fshelper.IsFile(dest) {
			return "", validation.NewError("'output' path exists")
		}

		return dest, nil
	}

	if !fshelper.PathExists(outPath) {
		return outPath, nil
	}

	if fshelper.IsDirectory(outPath) {
		dest := filepath.Join(outPath, fileName)
		if fshelper.PathExists(dest) && fshelper.IsFile(dest) {
			return "", validation.NewError("'output' path exists")
		}

		return dest, nil
	}

	if fshelper.IsFile(outPath) {
		return "", validation.NewError("'output' path exists")
	}

	return "", validation.NewError("invalid output path")
}
