package ai

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type MultipartField struct {
	Key   string
	Value string
}

func prepareMultipartBody(archivePath string, fields ...MultipartField) (*bytes.Buffer, string, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return nil, "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	filename := filepath.Base(archivePath)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, "", err
	}

	if _, err = io.Copy(part, file); err != nil {
		return nil, "", err
	}

	for _, field := range fields {
		ff, err := writer.CreateFormField(field.Key)
		if err != nil {
			return nil, "", err
		}
		if _, err = ff.Write([]byte(field.Value)); err != nil {
			return nil, "", err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	return body, writer.FormDataContentType(), nil
}

func prepareArchive(sourcePath string) (string, error) {
	// Проверяем существование пути
	info, err := os.Stat(sourcePath)
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}

	// Если это ZIP архив - возвращаем путь как есть
	if !info.IsDir() && strings.HasSuffix(strings.ToLower(sourcePath), ".zip") {
		return sourcePath, nil
	}

	// Создаем временный файл для архива
	tmpFile, err := os.CreateTemp("", "archive_*.zip")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	archivePath := tmpFile.Name()

	// Создаем ZIP архив
	zipWriter := zip.NewWriter(tmpFile)
	defer zipWriter.Close()

	// Функция для добавления файла в архив
	addFileToZip := func(filePath string, info os.FileInfo) error {
		// Открываем файл
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Создаем заголовок файла в архиве
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Устанавливаем метод сжатия
		header.Method = zip.Deflate

		// Получаем относительный путь для архива
		relPath, err := filepath.Rel(sourcePath, filePath)
		if err != nil {
			// Если не получается получить относительный путь, используем полный путь
			relPath = filePath
		}

		// Заменяем разделители пути на Unix-style для совместимости
		header.Name = filepath.ToSlash(relPath)

		// Создаем запись в архиве
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// Копируем содержимое файла в архив
		_, err = io.Copy(writer, file)
		return err
	}

	// Функция для добавления директории в архив (пустой)
	addDirToZip := func(dirPath string, info os.FileInfo) error {
		// Получаем относительный путь для архива
		relPath, err := filepath.Rel(sourcePath, dirPath)
		if err != nil {
			relPath = dirPath
		}

		// Создаем запись директории (добавляем trailing slash)
		header := &zip.FileHeader{
			Name:     filepath.ToSlash(relPath) + "/",
			Method:   zip.Deflate,
			Modified: info.ModTime(),
		}

		_, err = zipWriter.CreateHeader(header)
		return err
	}

	if info.IsDir() {
		// Обрабатываем директорию рекурсивно
		err = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Пропускаем корневую директорию
			if path == sourcePath {
				return nil
			}

			if info.IsDir() {
				// Добавляем запись для директории
				return addDirToZip(path, info)
			} else {
				// Добавляем файл
				return addFileToZip(path, info)
			}
		})

		if err != nil {
			os.Remove(archivePath) // Удаляем временный файл в случае ошибки
			return "", fmt.Errorf("failed to walk directory: %w", err)
		}
	} else {
		filename := filepath.Base(sourcePath)

		header := &zip.FileHeader{
			Name:     filepath.ToSlash(filename),
			Method:   zip.Deflate,
			Modified: info.ModTime(),
		}
		header.SetMode(info.Mode())

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			os.Remove(archivePath)
			return "", err
		}

		file, err := os.Open(sourcePath)
		if err != nil {
			os.Remove(archivePath)
			return "", err
		}
		defer file.Close()

		if _, err = io.Copy(writer, file); err != nil {
			os.Remove(archivePath)
			return "", err
		}
	}

	// Закрываем writer чтобы записать все данные
	if err := zipWriter.Close(); err != nil {
		os.Remove(archivePath)
		return "", fmt.Errorf("failed to close zip writer: %w", err)
	}

	return archivePath, nil
}

func createStubScanTarget() (string, error) {
	tempDir, err := os.MkdirTemp("", "source_*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	exampleFile := filepath.Join(tempDir, "aictl.temp")
	if err := os.WriteFile(exampleFile, []byte("this is temporal file for creation not empty branch"), 0644); err != nil {
		return "", fmt.Errorf("failed to create example file: %w", err)
	}

	return tempDir, nil
}

func getOrDefault[T any](value *T, defaultValue T) T {
	if value == nil {
		return defaultValue
	}

	return *value
}

func reference[T any](value T) *T {
	return &value
}
