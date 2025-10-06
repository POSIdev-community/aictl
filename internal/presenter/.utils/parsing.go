package _utils

import (
	"bufio"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
	"os"
)

func ParseUUIDs(args []string) ([]uuid.UUID, error) {
	uuids := make([]uuid.UUID, len(args))
	for i, arg := range args {
		id, err := uuid.Parse(arg)
		if err != nil {
			return nil, errs.NewValidationFieldError(arg, "invalid uuid")
		}

		uuids[i] = id
	}

	return uuids, nil
}

func ReadArgsFromStdin(args []string) []string {
	if len(args) != 1 || args[0] != "-" {
		return args
	}

	scanner := bufio.NewScanner(os.Stdin)
	var inputLines []string

	for scanner.Scan() {
		inputLines = append(inputLines, scanner.Text())
	}

	return inputLines
}
