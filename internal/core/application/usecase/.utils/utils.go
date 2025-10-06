package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CopyFileToPath(srcFile io.ReadCloser, destPath string, filename string) error {
	fullDestPath := filepath.Join(destPath, filename)

	destFile, err := os.Create(fullDestPath)
	if err != nil {
		return fmt.Errorf("не удалось создать файл назначения: %v", err)
	}

	defer func(destFile *os.File) {
		err := destFile.Close()
		if err != nil {
			// TODO log it
		}
	}(destFile)

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("ошибка копирования файла: %v", err)
	}

	return nil
}
