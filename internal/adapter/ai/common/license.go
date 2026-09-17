package common

import (
	domainlicense "github.com/POSIdev-community/aictl/internal/core/domain/license"
)

// LicenseFromLanguages builds a License with legacy LicensedModules (nil).
func LicenseFromLanguages(languages []string) *domainlicense.License {
	return &domainlicense.License{Languages: languages}
}

// MapLanguageGroups converts API language group enums to strings.
func MapLanguageGroups[T ~string](langs *[]T) []string {
	if langs == nil {
		return nil
	}
	out := make([]string, len(*langs))
	for i := range *langs {
		out[i] = string((*langs)[i])
	}

	return out
}

// LicensedModuleInput is a version-agnostic licensed module row from the API.
type LicensedModuleInput struct {
	ID              string
	Enabled         bool
	ScanModuleTypes []string
}

// LicenseFromAPI builds a domain License. modulesPresent distinguishes nil (legacy) from empty slice (new format).
func LicenseFromAPI(languages []string, modules []LicensedModuleInput, modulesPresent bool) *domainlicense.License {
	lic := &domainlicense.License{Languages: languages}
	if !modulesPresent {
		return lic
	}

	out := make([]domainlicense.LicensedModule, 0, len(modules))
	for _, m := range modules {
		out = append(out, domainlicense.LicensedModule{
			ID:              m.ID,
			Enabled:         m.Enabled,
			ScanModuleTypes: append([]string(nil), m.ScanModuleTypes...),
		})
	}
	lic.LicensedModules = out

	return lic
}
