package utils

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/securitypolicies"
)

type ProjectPoliciesAI interface {
	GetProjectPolicies(ctx context.Context, projectId uuid.UUID) (io.ReadCloser, error)
}

func LoadProjectPolicies(ctx context.Context, ai ProjectPoliciesAI, projectId uuid.UUID) (securitypolicies.Model, error) {
	r, err := ai.GetProjectPolicies(ctx, projectId)
	if err != nil {
		return securitypolicies.Model{}, fmt.Errorf("get project policies: %w", err)
	}
	defer func() {
		_ = r.Close()
	}()

	body, err := io.ReadAll(r)
	if err != nil {
		return securitypolicies.Model{}, fmt.Errorf("read project policies: %w", err)
	}

	model, err := securitypolicies.Parse(body)
	if err != nil {
		return securitypolicies.Model{}, fmt.Errorf("parse project policies: %w", err)
	}

	return model, nil
}

type ProjectPoliciesWriter interface {
	ProjectPoliciesAI
	SetProjectPolicies(ctx context.Context, projectId uuid.UUID, rawJSON []byte) error
}

func SetProjectPolicyCheck(ctx context.Context, ai ProjectPoliciesWriter, projectId uuid.UUID, enabled bool) error {
	current, err := LoadProjectPolicies(ctx, ai, projectId)
	if err != nil {
		return err
	}

	raw, err := securitypolicies.Encode(securitypolicies.Model{
		CheckAccordance: enabled,
		Policies:        current.Policies,
	})
	if err != nil {
		return fmt.Errorf("encode security policies: %w", err)
	}

	if err := ai.SetProjectPolicies(ctx, projectId, raw); err != nil {
		return fmt.Errorf("set project policies: %w", err)
	}

	return nil
}
