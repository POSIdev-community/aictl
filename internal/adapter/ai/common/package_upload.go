package common

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
)

// PackageFileMeta holds metadata required by POST /api/packages/{type}.
type PackageFileMeta struct {
	FileName string
	FileSize int64
	Hash     string // MD5 hex lowercase
}

// ComputePackageFileMeta returns basename, size and MD5 (lowercase hex) of path.
func ComputePackageFileMeta(path string) (PackageFileMeta, error) {
	file, err := os.Open(path)
	if err != nil {
		return PackageFileMeta{}, fmt.Errorf("open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	fi, err := file.Stat()
	if err != nil {
		return PackageFileMeta{}, fmt.Errorf("stat file: %w", err)
	}
	if !fi.Mode().IsRegular() {
		return PackageFileMeta{}, errors.New("path is not a regular file")
	}
	if fi.Size() == 0 {
		return PackageFileMeta{}, errors.New("package file is empty")
	}

	hasher := md5.New() //nolint:gosec // MD5 required by AIE Packages API
	if _, err = io.Copy(hasher, file); err != nil {
		return PackageFileMeta{}, fmt.Errorf("compute md5: %w", err)
	}

	return PackageFileMeta{
		FileName: filepath.Base(path),
		FileSize: fi.Size(),
		Hash:     hex.EncodeToString(hasher.Sum(nil)),
	}, nil
}

// PreparePackageMultipartBody builds multipart body for package upload.
// Form file field name is "package" (Packages API), not "file".
// onProgress receives raw upload percent 0..100; pass nil to disable.
func PreparePackageMultipartBody(
	ctx context.Context,
	path string,
	version string,
	meta PackageFileMeta,
	onProgress func(percent int),
) (io.ReadCloser, string, error) {
	if onProgress == nil {
		onProgress = func(int) {}
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, "", fmt.Errorf("open file: %w", err)
	}

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	contentType := writer.FormDataContentType()

	go func() {
		defer func() {
			_ = writer.Close()
			_ = pw.Close()
		}()

		done := make(chan struct{})
		defer close(done)
		go func() {
			select {
			case <-ctx.Done():
				_ = pw.CloseWithError(ctx.Err())
			case <-done:
			}
		}()

		part, err := writer.CreateFormFile("package", meta.FileName)
		if err != nil {
			_ = pw.CloseWithError(fmt.Errorf("create form file: %w", err))
			return
		}

		buf := make([]byte, 1*1024*1024)
		uploaded := int64(0)
		totalSize := meta.FileSize
		bytesPerPercent := totalSize / 100
		if bytesPerPercent == 0 {
			bytesPerPercent = 1
		}

		onProgress(0)
		for {
			n, readErr := file.Read(buf)
			if n > 0 {
				if _, writeErr := part.Write(buf[:n]); writeErr != nil {
					_ = pw.CloseWithError(fmt.Errorf("write to multipart part: %w", writeErr))
					return
				}
				uploaded += int64(n)

				currentPercent := int(uploaded / bytesPerPercent)
				if currentPercent > 100 {
					currentPercent = 100
				}
				onProgress(currentPercent)
			}

			if readErr == io.EOF {
				break
			}
			if readErr != nil {
				_ = pw.CloseWithError(fmt.Errorf("read package: %w", readErr))
				return
			}
		}
		onProgress(100)

		fields := []MultipartField{
			{Key: "version", Value: version},
			{Key: "fileName", Value: meta.FileName},
			{Key: "fileSize", Value: strconv.FormatInt(meta.FileSize, 10)},
			{Key: "hash", Value: meta.Hash},
		}
		for _, field := range fields {
			if err := writer.WriteField(field.Key, field.Value); err != nil {
				_ = pw.CloseWithError(fmt.Errorf("write field %q: %w", field.Key, err))
				return
			}
		}
	}()

	return &multipartReadCloser{Reader: pr, pw: pw, file: file}, contentType, nil
}
