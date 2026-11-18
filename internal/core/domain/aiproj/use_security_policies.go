package aiproj

// UseSecurityPolicies returns aiproj UseSecurityPolicies when set, else nil.
func (r *Result) UseSecurityPolicies() *bool {
	if r == nil {
		return nil
	}

	switch r.Version {
	case "1.9":
		if r.V19 == nil {
			return nil
		}

		return r.V19.UseSecurityPolicies
	case "1.10":
		if r.V110 == nil {
			return nil
		}

		return r.V110.UseSecurityPolicies
	case "1.11":
		if r.V111 == nil {
			return nil
		}

		return r.V111.UseSecurityPolicies
	case "1.12":
		if r.V112 == nil {
			return nil
		}

		return r.V112.UseSecurityPolicies
	default:
		return nil
	}
}
