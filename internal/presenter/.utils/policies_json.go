package _utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// NormalizedPolicies is the result of parsing user input for set project policies.
// Check == nil means the caller should preserve the server's current value.
type NormalizedPolicies struct {
	PoliciesJSON []byte
	Check        *bool
}

// NormalizePoliciesJSON prepares policies payload for set project policies.
//
// Accepts the API SecurityPoliciesModel object, or a policies JSON array
// (as in aisa --policy-settings-file). Comments inside the policies payload
// are preserved. When checkSecurityPoliciesAccordance is omitted (array input
// or object without the field), Check is nil.
func NormalizePoliciesJSON(raw []byte) (NormalizedPolicies, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return NormalizedPolicies{}, fmt.Errorf("empty")
	}

	switch trimmed[0] {
	case '[':
		return NormalizedPolicies{PoliciesJSON: trimmed}, nil
	case '{':
		return normalizePoliciesObject(trimmed)
	default:
		return NormalizedPolicies{}, fmt.Errorf("policies data must be a JSON array or object")
	}
}

func normalizePoliciesObject(raw []byte) (NormalizedPolicies, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		// Likely a bare policy object with comments — wrap as a one-element array textually.
		wrapped := make([]byte, 0, len(raw)+2)
		wrapped = append(wrapped, '[')
		wrapped = append(wrapped, raw...)
		wrapped = append(wrapped, ']')

		return NormalizedPolicies{PoliciesJSON: wrapped}, nil
	}

	sp, ok := fields["securityPolicies"]
	if !ok {
		wrapped := make([]byte, 0, len(raw)+2)
		wrapped = append(wrapped, '[')
		wrapped = append(wrapped, raw...)
		wrapped = append(wrapped, ']')

		return NormalizedPolicies{PoliciesJSON: wrapped}, nil
	}

	var check *bool
	if c, exists := fields["checkSecurityPoliciesAccordance"]; exists {
		var v bool
		if err := json.Unmarshal(c, &v); err != nil {
			return NormalizedPolicies{}, fmt.Errorf("checkSecurityPoliciesAccordance: %w", err)
		}
		check = &v
	}

	sp = bytes.TrimSpace(sp)
	switch {
	case len(sp) == 0 || string(sp) == "null":
		return NormalizedPolicies{PoliciesJSON: []byte("[]"), Check: check}, nil
	case sp[0] == '[':
		return NormalizedPolicies{PoliciesJSON: sp, Check: check}, nil
	case sp[0] == '"':
		var s string
		if err := json.Unmarshal(sp, &s); err != nil {
			return NormalizedPolicies{}, err
		}

		s = string(bytes.TrimSpace([]byte(s)))
		if s == "" {
			s = "[]"
		}

		return NormalizedPolicies{PoliciesJSON: []byte(s), Check: check}, nil
	default:
		return NormalizedPolicies{}, fmt.Errorf("securityPolicies must be a JSON array or a JSON string")
	}
}
