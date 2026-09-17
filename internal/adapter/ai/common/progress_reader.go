package common

import (
	"io"
)

// ProgressReadCloser wraps r and reports read progress 0..100 via onProgress.
// totalSize <= 0 disables percent tracking (still reports 0 on first Read and 100 on EOF).
func ProgressReadCloser(r io.ReadCloser, totalSize int64, onProgress func(percent int)) io.ReadCloser {
	if onProgress == nil {
		onProgress = func(int) {}
	}

	return &progressReadCloser{
		r:          r,
		totalSize:  totalSize,
		onProgress: onProgress,
	}
}

type progressReadCloser struct {
	r          io.ReadCloser
	totalSize  int64
	read       int64
	started    bool
	onProgress func(percent int)
}

func (p *progressReadCloser) Read(buf []byte) (int, error) {
	if !p.started {
		p.started = true
		p.onProgress(0)
	}

	n, err := p.r.Read(buf)
	if n > 0 {
		p.read += int64(n)
		if p.totalSize > 0 {
			percent := int(p.read * 100 / p.totalSize)
			if percent > 100 {
				percent = 100
			}
			p.onProgress(percent)
		}
	}

	if err == io.EOF {
		p.onProgress(100)
	}

	return n, err
}

func (p *progressReadCloser) Close() error {
	return p.r.Close()
}
