package policies

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetProjectPolicies(ctx context.Context, projectId uuid.UUID) (io.ReadCloser, error)
}

type CLI interface {
	ReturnText(ctx context.Context, text string)
}

type UseCase struct {
	aiAdapter  AI
	cliAdapter CLI
	cfg        *config.Config
}

func NewUseCase(aiAdapter AI, cliAdapter CLI, cfg *config.Config) (*UseCase, error) {
	if aiAdapter == nil {
		return nil, validation.NewRequiredError("aiAdapter")
	}

	if cliAdapter == nil {
		return nil, validation.NewRequiredError("cliAdapter")
	}

	return &UseCase{
		aiAdapter:  aiAdapter,
		cliAdapter: cliAdapter,
		cfg:        cfg,
	}, nil
}

func (u *UseCase) Execute(ctx context.Context) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	r, err := u.aiAdapter.GetProjectPolicies(ctx, u.cfg.ProjectId())
	if err != nil {
		return fmt.Errorf("get project policies: %w", err)
	}
	defer func() {
		_ = r.Close()
	}()

	body, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("read project policies: %w", err)
	}

	formatted, err := formatProjectPolicies(body)
	if err != nil {
		return fmt.Errorf("format project policies: %w", err)
	}

	u.cliAdapter.ReturnText(ctx, formatted)

	return nil
}

// formatProjectPolicies extracts the policies payload from SecurityPoliciesModel
// (securityPolicies JSON string) and returns it pretty-printed when possible.
// If the payload is not strict JSON (e.g. contains comments), it is returned as-is.
func formatProjectPolicies(body []byte) (string, error) {
	body = bytes.TrimSpace(body)
	if len(body) == 0 {
		return "[]\n", nil
	}

	policies, err := extractPoliciesJSON(body)
	if err != nil {
		return "", err
	}

	var v any
	if err := json.Unmarshal(policies, &v); err != nil {
		out := string(policies)
		if !strings.HasSuffix(out, "\n") {
			out += "\n"
		}

		return out, nil
	}

	out, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return "", fmt.Errorf("indent policies json: %w", err)
	}

	return string(out) + "\n", nil
}

func extractPoliciesJSON(body []byte) ([]byte, error) {
	if body[0] == '[' {
		return body, nil
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(body, &fields); err != nil {
		return nil, fmt.Errorf("parse security policies model: %w", err)
	}

	sp, ok := fields["securityPolicies"]
	if !ok || len(bytes.TrimSpace(sp)) == 0 || string(bytes.TrimSpace(sp)) == "null" {
		return []byte("[]"), nil
	}

	sp = bytes.TrimSpace(sp)
	switch sp[0] {
	case '[':
		return sp, nil
	case '"':
		var s string
		if err := json.Unmarshal(sp, &s); err != nil {
			return nil, fmt.Errorf("securityPolicies string: %w", err)
		}

		s = strings.TrimSpace(s)
		if s == "" {
			return []byte("[]"), nil
		}

		return []byte(s), nil
	default:
		return nil, fmt.Errorf("securityPolicies must be a JSON array or a JSON string")
	}
}
