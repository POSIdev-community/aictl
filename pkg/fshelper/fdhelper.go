package fshelper

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

func IsArchive(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".zip":
		return true
	default:
		return false
	}
}
