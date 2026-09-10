package _utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
)

// ReadAiprojInput loads aiproj JSON from -f path, positional arg (or "-" stdin), or piped stdin.
// When filePath is set, args are ignored (same as set project settings).
func ReadAiprojInput(filePath string, args []string) ([]byte, error) {
	return ReadAiprojInputFrom(filePath, args, os.Stdin)
}

func ReadAiprojInputFrom(filePath string, args []string, stdin io.Reader) ([]byte, error) {
	var raw []byte

	switch {
	case filePath != "":
		if !fshelper.PathExists(filePath) {
			return nil, validation.NewMessageError(fmt.Sprintf("file %s does not exist", filePath))
		}

		if !fshelper.IsFile(filePath) {
			return nil, validation.NewMessageError(fmt.Sprintf("path %s does not a file", filePath))
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, validation.NewMessageError(fmt.Sprintf("read aiproj file: %v", err))
		}

		raw = content
	default:
		args = ReadArgsFromReader(args, stdin)
		if len(args) == 0 {
			piped, ok := readStdinIfPiped(stdin)
			if !ok || len(piped) == 0 {
				return nil, validation.NewMessageError("aiproj data required")
			}

			raw = piped
		} else {
			raw = []byte(args[0])
		}
	}

	if len(raw) == 0 {
		return nil, validation.NewMessageError("aiproj data required")
	}

	if !json.Valid(raw) {
		return nil, validation.NewMessageError("invalid aiproj data: not valid json")
	}

	return raw, nil
}

func readStdinIfPiped(stdin io.Reader) ([]byte, bool) {
	f, ok := stdin.(*os.File)
	if ok {
		fi, err := f.Stat()
		if err != nil || (fi.Mode()&os.ModeCharDevice) != 0 {
			return nil, false
		}
	}

	data, err := io.ReadAll(stdin)
	if err != nil {
		return nil, false
	}

	return data, true
}
