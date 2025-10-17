package utils

import (
	"fmt"
	"io"
	"os"
)

func CopyFileToPath(srcFile io.ReadCloser, fullDestPath string) error {
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
