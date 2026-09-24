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
	"github.com/POSIdev-community/aictl/internal/core/domain/securitypolicies"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	usecaseutils "github.com/POSIdev-community/aictl/internal/core/usecase/.utils"
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

	model, err := usecaseutils.LoadProjectPolicies(ctx, u.aiAdapter, u.cfg.ProjectId())
	if err != nil {
		return err
	}

	formatted, err := formatPoliciesArray(model.Policies)
	if err != nil {
		return fmt.Errorf("format project policies: %w", err)
	}

	u.cliAdapter.ReturnText(ctx, formatted)

	return nil
}

// formatPoliciesArray pretty-prints the policies JSON array when possible.
// If the payload is not strict JSON (e.g. contains comments), it is returned as-is.
func formatPoliciesArray(policies string) (string, error) {
	policies = strings.TrimSpace(policies)
	if policies == "" {
		return "[]\n", nil
	}

	var v any
	if err := json.Unmarshal([]byte(policies), &v); err != nil {
		out := policies
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

// formatProjectPolicies kept for tests that pass full SecurityPoliciesModel bodies.
func formatProjectPolicies(body []byte) (string, error) {
	body = bytes.TrimSpace(body)
	if len(body) == 0 {
		return "[]\n", nil
	}

	model, err := securitypolicies.Parse(body)
	if err != nil {
		return "", err
	}

	return formatPoliciesArray(model.Policies)
}
