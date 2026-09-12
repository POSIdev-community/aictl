package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/internal/core/domain/version"
)

func RequireProjectScanSettings(serverVersion version.Version) error {
	minVersion, err := version.NewVersion("6.0.0")
	if err != nil {
		return fmt.Errorf("parse min version: %w", err)
	}

	if serverVersion.Less(minVersion) {
		return validation.NewError("priority and preferred agents settings are not supported on server version 5.x")
	}

	return nil
}

func CopyFileToPath(srcFile io.ReadCloser, fullDestPath string) error {
	destFile, err := os.Create(fullDestPath)
	if err != nil {
		return fmt.Errorf("create target file: %w", err)
	}

	_, copyErr := io.Copy(destFile, srcFile)
	closeErr := destFile.Close()
	if copyErr != nil {
		_ = os.Remove(fullDestPath)

		return fmt.Errorf("copy file: %w", copyErr)
	}
	if closeErr != nil {
		_ = os.Remove(fullDestPath)

		return fmt.Errorf("close target file: %w", closeErr)
	}

	return nil
}
