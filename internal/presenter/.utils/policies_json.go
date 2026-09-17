package _utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// NormalizePoliciesJSON prepares payload for PUT security policies.
//
// The AI API expects SecurityPoliciesModel:
//
//	{"checkSecurityPoliciesAccordance":bool,"securityPolicies":"<json string>"}
//
// where securityPolicies is a JSON *string* containing the policies array.
// A bare policies array (aisa --policy-settings-file) is wrapped automatically.
// Comments inside the policies payload are preserved (embedded as the string value).
func NormalizePoliciesJSON(raw []byte) ([]byte, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return nil, fmt.Errorf("empty")
	}

	switch trimmed[0] {
	case '[':
		// Keep original bytes (including // and /* */ comments) inside securityPolicies.
		return encodePoliciesModel(false, trimmed)
	case '{':
		return normalizePoliciesObject(trimmed)
	default:
		return nil, fmt.Errorf("policies data must be a JSON array or object")
	}
}

func normalizePoliciesObject(raw []byte) ([]byte, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		// Likely a bare policy object with comments — wrap as a one-element array textually.
		wrapped := make([]byte, 0, len(raw)+2)
		wrapped = append(wrapped, '[')
		wrapped = append(wrapped, raw...)
		wrapped = append(wrapped, ']')

		return encodePoliciesModel(false, wrapped)
	}

	sp, ok := fields["securityPolicies"]
	if !ok {
		wrapped := make([]byte, 0, len(raw)+2)
		wrapped = append(wrapped, '[')
		wrapped = append(wrapped, raw...)
		wrapped = append(wrapped, ']')

		return encodePoliciesModel(false, wrapped)
	}

	check := false
	if c, exists := fields["checkSecurityPoliciesAccordance"]; exists {
		if err := json.Unmarshal(c, &check); err != nil {
			return nil, fmt.Errorf("checkSecurityPoliciesAccordance: %w", err)
		}
	}

	sp = bytes.TrimSpace(sp)
	switch {
	case len(sp) == 0 || string(sp) == "null":
		return encodePoliciesModel(check, []byte("[]"))
	case sp[0] == '[':
		return encodePoliciesModel(check, sp)
	case sp[0] == '"':
		var s string
		if err := json.Unmarshal(sp, &s); err != nil {
			return nil, err
		}

		s = string(bytes.TrimSpace([]byte(s)))
		if s == "" {
			s = "[]"
		}

		return encodePoliciesModel(check, []byte(s))
	default:
		return nil, fmt.Errorf("securityPolicies must be a JSON array or a JSON string")
	}
}

func encodePoliciesModel(checkAccordance bool, policiesJSON []byte) ([]byte, error) {
	return json.Marshal(struct {
		CheckSecurityPoliciesAccordance bool   `json:"checkSecurityPoliciesAccordance"`
		SecurityPolicies                string `json:"securityPolicies"`
	}{
		CheckSecurityPoliciesAccordance: checkAccordance,
		SecurityPolicies:                string(policiesJSON),
	})
}
