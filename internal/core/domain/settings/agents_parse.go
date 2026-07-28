package settings

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

func ParsePreferredAgentsCSV(s string) ([]uuid.UUID, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return []uuid.UUID{}, nil
	}

	parts := strings.Split(s, ",")
	agents := make([]uuid.UUID, 0, len(parts))
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return nil, validation.NewFieldError("agents", fmt.Sprintf("empty agent id at position %d", i+1))
		}

		id, err := uuid.Parse(part)
		if err != nil {
			return nil, validation.NewFieldError("agents", fmt.Sprintf("invalid uuid %q", part))
		}

		agents = append(agents, id)
	}

	return agents, nil
}
