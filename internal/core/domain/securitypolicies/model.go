package securitypolicies

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// Model is the AI API SecurityPoliciesModel.
type Model struct {
	CheckAccordance bool
	// Policies is the raw JSON text stored in securityPolicies (usually an array).
	Policies string
}

// Parse reads a SecurityPoliciesModel JSON body (or a bare policies array).
func Parse(raw []byte) (Model, error) {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 {
		return Model{Policies: "[]"}, nil
	}

	if raw[0] == '[' {
		return Model{Policies: string(raw)}, nil
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return Model{}, fmt.Errorf("parse security policies model: %w", err)
	}

	check := false
	if c, ok := fields["checkSecurityPoliciesAccordance"]; ok {
		if err := json.Unmarshal(c, &check); err != nil {
			return Model{}, fmt.Errorf("checkSecurityPoliciesAccordance: %w", err)
		}
	}

	policies, err := extractPoliciesJSON(fields)
	if err != nil {
		return Model{}, err
	}

	return Model{CheckAccordance: check, Policies: string(policies)}, nil
}

// Encode marshals Model into API SecurityPoliciesModel JSON.
func Encode(m Model) ([]byte, error) {
	policies := strings.TrimSpace(m.Policies)
	if policies == "" {
		policies = "[]"
	}

	return json.Marshal(struct {
		CheckSecurityPoliciesAccordance bool   `json:"checkSecurityPoliciesAccordance"`
		SecurityPolicies                string `json:"securityPolicies"`
	}{
		CheckSecurityPoliciesAccordance: m.CheckAccordance,
		SecurityPolicies:                policies,
	})
}

func extractPoliciesJSON(fields map[string]json.RawMessage) ([]byte, error) {
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
