package ai

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/POSIdev-community/aictl/pkg/logger"
)

type MultipartField struct {
	Key   string
	Value string
}

func prepareMultipartBody(
	ctx context.Context,
	archivePath string,
	fields ...MultipartField) (io.ReadCloser, string, error) {

	log := logger.FromContext(ctx)
	progress := createProgressReporter(log)

	file, err := os.Open(archivePath)
	if err != nil {
		return nil, "", fmt.Errorf("error opening file: %w", err)
	}

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	contentType := writer.FormDataContentType()

	go func() {

		defer func() {
			writer.Close()
			pw.Close()
		}()

		done := make(chan struct{})
		defer close(done)
		go func() {
			select {
			case <-ctx.Done():
				pw.CloseWithError(ctx.Err())
			case <-done:
			}
		}()

		filename := filepath.Base(archivePath)
		part, err := writer.CreateFormFile("file", filename)
		if err != nil {
			pw.CloseWithError(fmt.Errorf("failed to create form file: %w", err))
			return
		}

		fi, err := file.Stat()
		if err != nil {
			pw.CloseWithError(fmt.Errorf("stat file: %w", err))
			return
		}
		totalSize := fi.Size()
		if totalSize == 0 {
			pw.CloseWithError(errors.New("archive is empty"))
			return
		}

		buf := make([]byte, 1*1024*1024)
		uploaded := int64(0)

		bytesPerPercent := totalSize / 100
		if bytesPerPercent == 0 {
			bytesPerPercent = 1 // на случай очень маленьких файлов
		}

		progress(0)
		for {
			n, readErr := file.Read(buf)
			if n > 0 {
				// Пишем чанк в multipart
				if _, writeErr := part.Write(buf[:n]); writeErr != nil {
					pw.CloseWithError(fmt.Errorf("write to multipart part: %w", writeErr))
					return
				}
				uploaded += int64(n)

				currentPercent := int(uploaded / bytesPerPercent)
				if currentPercent > 100 {
					currentPercent = 100
				}

				progress(currentPercent)
			}

			if readErr == io.EOF {
				break
			}
			if readErr != nil {
				pw.CloseWithError(fmt.Errorf("read archive: %w", readErr))
				return
			}
		}

		if progress != nil {
			progress(100)
		}

		// Поля (после файла — стандарт для multipart)
		for _, field := range fields {
			if err := writer.WriteField(field.Key, field.Value); err != nil {
				pw.CloseWithError(fmt.Errorf("failed to write field %q: %w", field.Key, err))
				return
			}
		}
	}()

	return &multipartReadCloser{Reader: pr, pw: pw, file: file}, contentType, nil
}

type multipartReadCloser struct {
	io.Reader
	pw   *io.PipeWriter
	file *os.File
}

func (mrc *multipartReadCloser) Close() error {
	_ = mrc.pw.CloseWithError(errors.New("body closed by caller"))
	return mrc.file.Close()
}

func prepareArchive(sourcePath string) (archivePath string, err error) {
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

	archivePath = tmpFile.Name()

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
		err = filepath.WalkDir(sourcePath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// Пропускаем корневую директорию
			if path == sourcePath {
				return nil
			}

			if d.IsDir() {
				info, err := d.Info()
				if err != nil {
					return err
				}
				return addDirToZip(path, info)
			}

			if d.Type().IsRegular() {
				info, err := d.Info()
				if err != nil {
					return err
				}
				return addFileToZip(path, info)
			}

			return nil
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

func createProgressReporter(log *logger.Logger) func(int) {
	var lastPrintedPercent = -1

	return func(sentPercent int) {
		percent := sentPercent / 10 * 10

		if percent > lastPrintedPercent {
			lastPrintedPercent = percent

			log.StdErrf("updating sources: %d%%", percent)
		}
	}
}
