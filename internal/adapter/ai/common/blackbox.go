package common

import (
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
)

type BlackBoxAddress struct {
	Address *string
	Format  *string
}

type BlackBoxHTTPHeader struct {
	Key   *string
	Value *string
}

type BlackBoxAuthForm struct {
	FormDetection      *string
	FormAddress        *string
	FormXPath          *string
	Login              *string
	LoginKey           *string
	Password           *string
	PasswordKey        *string
	ValidationTemplate *string
}

type BlackBoxAuthCookie struct {
	Cookie             *string
	ValidationAddress  *string
	ValidationTemplate *string
}

type BlackBoxAuthHTTP struct {
	Login             *string
	Password          *string
	ValidationAddress *string
}

type BlackBoxAuth struct {
	Type   *string
	Form   *BlackBoxAuthForm
	Cookie *BlackBoxAuthCookie
	Http   *BlackBoxAuthHTTP
}

type BlackBoxProxy struct {
	Enabled  *bool
	Host     *string
	Login    *string
	Password *string
	Port     *int
	Type     *string
}

type BlackBoxPayload struct {
	Site                  *string
	Level                 *string
	ScanScope             *string
	SslCheck              *bool
	RunAutocheckAfterScan *bool
	AdditionalHttpHeaders *[]BlackBoxHTTPHeader
	WhiteListedAddresses  *[]BlackBoxAddress
	BlackListedAddresses  *[]BlackBoxAddress
	Authentication        *BlackBoxAuth
	ProxySettings         *BlackBoxProxy
	IsActive              *bool
}

func BuildBlackBoxPayload(source settings.BlackBoxSettings, enabled bool) BlackBoxPayload {
	payload := BlackBoxPayload{
		Site:                  Reference(source.Site),
		Level:                 Reference(source.Level),
		ScanScope:             Reference(source.ScanScope),
		SslCheck:              Reference(source.SslCheck),
		RunAutocheckAfterScan: Reference(source.RunAutocheckAfterScan),
		IsActive:              Reference(enabled),
	}

	if len(source.AdditionalHttpHeaders) > 0 {
		headers := make([]BlackBoxHTTPHeader, len(source.AdditionalHttpHeaders))
		for i, header := range source.AdditionalHttpHeaders {
			headers[i] = BlackBoxHTTPHeader{
				Key:   Reference(header.Key),
				Value: Reference(header.Value),
			}
		}
		payload.AdditionalHttpHeaders = &headers
	}

	payload.WhiteListedAddresses = mapAddressEntries(source.WhiteListedAddresses)
	payload.BlackListedAddresses = mapAddressEntries(source.BlackListedAddresses)
	payload.Authentication = mapBlackBoxAuth(source.Authentication)
	payload.ProxySettings = mapBlackBoxProxy(source.ProxySettings)

	return payload
}

func mapAddressEntries(entries []settings.AddressEntry) *[]BlackBoxAddress {
	if len(entries) == 0 {
		return nil
	}

	result := make([]BlackBoxAddress, len(entries))
	for i, entry := range entries {
		result[i] = BlackBoxAddress{
			Address: Reference(entry.Address),
			Format:  Reference(entry.Format),
		}
	}

	return &result
}

func mapBlackBoxAuth(auth *settings.BlackBoxAuthentication) *BlackBoxAuth {
	if auth == nil {
		return nil
	}

	result := &BlackBoxAuth{
		Type: Reference(authType(auth)),
	}

	if auth.Cookie != nil {
		result.Cookie = &BlackBoxAuthCookie{
			Cookie:             Reference(auth.Cookie.Cookie),
			ValidationAddress:  Reference(auth.Cookie.ValidationAddress),
			ValidationTemplate: Reference(auth.Cookie.ValidationTemplate),
		}
	}

	if auth.Form != nil {
		result.Form = &BlackBoxAuthForm{
			FormDetection:      Reference(auth.Form.FormDetection),
			FormAddress:        Reference(auth.Form.FormAddress),
			FormXPath:          Reference(auth.Form.FormXPath),
			Login:              Reference(auth.Form.Login),
			LoginKey:           Reference(auth.Form.LoginKey),
			Password:           Reference(auth.Form.Password),
			PasswordKey:        Reference(auth.Form.PasswordKey),
			ValidationTemplate: Reference(auth.Form.ValidationTemplate),
		}
	}

	if auth.Http != nil {
		result.Http = &BlackBoxAuthHTTP{
			Login:             Reference(auth.Http.Login),
			Password:          Reference(auth.Http.Password),
			ValidationAddress: Reference(auth.Http.ValidationAddress),
		}
	}

	return result
}

func mapBlackBoxProxy(proxy *settings.BlackBoxProxySettings) *BlackBoxProxy {
	if proxy == nil {
		return nil
	}

	return &BlackBoxProxy{
		Enabled:  Reference(proxy.Enabled),
		Host:     Reference(proxy.Host),
		Login:    Reference(proxy.Login),
		Password: Reference(proxy.Password),
		Port:     Reference(proxy.Port),
		Type:     Reference(proxy.Type),
	}
}

// authType returns the explicit authentication type or, when it is omitted in aiproj,
// infers it from the single credentials block provided. Without a type the server
// silently falls back to None and black box scans without logging in.
func authType(auth *settings.BlackBoxAuthentication) string {
	if auth.Type != "" {
		return auth.Type
	}

	switch {
	case auth.Http != nil && auth.Form == nil && auth.Cookie == nil:
		return "Http"
	case auth.Form != nil && auth.Http == nil && auth.Cookie == nil:
		return "Form"
	case auth.Cookie != nil && auth.Http == nil && auth.Form == nil:
		return "RawCookie"
	default:
		return ""
	}
}
