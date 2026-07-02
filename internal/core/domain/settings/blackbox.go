package settings

import (
	v110 "github.com/POSIdev-community/aiproj/model/v1_10"
	v19 "github.com/POSIdev-community/aiproj/model/v1_9"
)

func (s *ScanSettings) applyBlackBoxFromV19(bb *v19.AIProjBlackBoxSettings) {
	if bb == nil {
		return
	}

	s.BlackBoxEnabled = true

	if bb.Site != nil {
		s.BlackBoxSettings.Site = *bb.Site
	}
	if bb.Level != nil {
		s.BlackBoxSettings.Level = string(*bb.Level)
	}
	if bb.ScanScope != nil {
		s.BlackBoxSettings.ScanScope = string(*bb.ScanScope)
	}
	if bb.SslCheck != nil {
		s.BlackBoxSettings.SslCheck = *bb.SslCheck
	}
	if bb.RunAutocheckAfterScan != nil {
		s.BlackBoxSettings.RunAutocheckAfterScan = *bb.RunAutocheckAfterScan
	}

	s.BlackBoxSettings.AdditionalHttpHeaders = mapHTTPHeadersV19(bb.AdditionalHttpHeaders)
	s.BlackBoxSettings.WhiteListedAddresses = mapWhiteListedAddressesV19(bb.WhiteListedAddresses)
	s.BlackBoxSettings.BlackListedAddresses = mapBlackListedAddressesV19(bb.BlackListedAddresses)
	s.BlackBoxSettings.Authentication = mapAuthenticationV19(bb.Authentication)
	s.BlackBoxSettings.ProxySettings = mapProxySettingsV19(bb.ProxySettings)
}

func (s *ScanSettings) applyBlackBoxFromV110(bb *v110.AIProjBlackBoxSettings) {
	if bb == nil {
		return
	}

	s.BlackBoxEnabled = true

	if bb.Site != nil {
		s.BlackBoxSettings.Site = *bb.Site
	}
	if bb.Level != nil {
		s.BlackBoxSettings.Level = string(*bb.Level)
	}
	if bb.ScanScope != nil {
		s.BlackBoxSettings.ScanScope = string(*bb.ScanScope)
	}
	if bb.SslCheck != nil {
		s.BlackBoxSettings.SslCheck = *bb.SslCheck
	}
	if bb.RunAutocheckAfterScan != nil {
		s.BlackBoxSettings.RunAutocheckAfterScan = *bb.RunAutocheckAfterScan
	}

	s.BlackBoxSettings.AdditionalHttpHeaders = mapHTTPHeadersV110(bb.AdditionalHttpHeaders)
	s.BlackBoxSettings.WhiteListedAddresses = mapWhiteListedAddressesV110(bb.WhiteListedAddresses)
	s.BlackBoxSettings.BlackListedAddresses = mapBlackListedAddressesV110(bb.BlackListedAddresses)
	s.BlackBoxSettings.Authentication = mapAuthenticationV110(bb.Authentication)
	s.BlackBoxSettings.ProxySettings = mapProxySettingsV110(bb.ProxySettings)
}

func mapHTTPHeadersV19(headers *v19.AIProjBlackBoxSettingsAdditionalHttpHeaders) []HTTPHeader {
	if headers == nil || len(*headers) == 0 {
		return nil
	}

	result := make([]HTTPHeader, len(*headers))
	for i, header := range *headers {
		if header.Key != nil {
			result[i].Key = *header.Key
		}
		if header.Value != nil {
			result[i].Value = *header.Value
		}
	}

	return result
}

func mapHTTPHeadersV110(headers *v110.AIProjBlackBoxSettingsAdditionalHttpHeaders) []HTTPHeader {
	if headers == nil || len(*headers) == 0 {
		return nil
	}

	result := make([]HTTPHeader, len(*headers))
	for i, header := range *headers {
		if header.Key != nil {
			result[i].Key = *header.Key
		}
		if header.Value != nil {
			result[i].Value = *header.Value
		}
	}

	return result
}

func mapWhiteListedAddressesV19(addresses *v19.AIProjBlackBoxSettingsWhiteListedAddresses) []AddressEntry {
	if addresses == nil || len(*addresses) == 0 {
		return nil
	}

	result := make([]AddressEntry, len(*addresses))
	for i, address := range *addresses {
		entry := AddressEntry{}
		if address.Address != nil {
			entry.Address = *address.Address
		}
		if address.Format != nil {
			entry.Format = string(*address.Format)
		}
		result[i] = entry
	}

	return result
}

func mapBlackListedAddressesV19(addresses *v19.AIProjBlackBoxSettingsBlackListedAddresses) []AddressEntry {
	if addresses == nil || len(*addresses) == 0 {
		return nil
	}

	result := make([]AddressEntry, len(*addresses))
	for i, address := range *addresses {
		entry := AddressEntry{}
		if address.Address != nil {
			entry.Address = *address.Address
		}
		if address.Format != nil {
			entry.Format = string(*address.Format)
		}
		result[i] = entry
	}

	return result
}

func mapWhiteListedAddressesV110(addresses *v110.AIProjBlackBoxSettingsWhiteListedAddresses) []AddressEntry {
	if addresses == nil || len(*addresses) == 0 {
		return nil
	}

	result := make([]AddressEntry, len(*addresses))
	for i, address := range *addresses {
		entry := AddressEntry{}
		if address.Address != nil {
			entry.Address = *address.Address
		}
		if address.Format != nil {
			entry.Format = string(*address.Format)
		}
		result[i] = entry
	}

	return result
}

func mapBlackListedAddressesV110(addresses *v110.AIProjBlackBoxSettingsBlackListedAddresses) []AddressEntry {
	if addresses == nil || len(*addresses) == 0 {
		return nil
	}

	result := make([]AddressEntry, len(*addresses))
	for i, address := range *addresses {
		entry := AddressEntry{}
		if address.Address != nil {
			entry.Address = *address.Address
		}
		if address.Format != nil {
			entry.Format = string(*address.Format)
		}
		result[i] = entry
	}

	return result
}

func mapAuthenticationV19(auth *v19.AIProjBlackBoxSettingsAuthentication) *BlackBoxAuthentication {
	if auth == nil {
		return nil
	}

	result := &BlackBoxAuthentication{}
	if auth.Type != nil {
		result.Type = string(*auth.Type)
	}

	if auth.Cookie != nil {
		result.Cookie = &BlackBoxCookieAuth{
			Cookie:             auth.Cookie.Cookie,
			ValidationAddress:  auth.Cookie.ValidationAddress,
			ValidationTemplate: auth.Cookie.ValidationTemplate,
		}
	}

	if auth.Form != nil {
		form := &BlackBoxFormAuth{
			FormDetection:      string(auth.Form.FormDetection),
			FormAddress:        auth.Form.FormAddress,
			Login:              auth.Form.Login,
			Password:           auth.Form.Password,
			ValidationTemplate: auth.Form.ValidationTemplate,
		}
		if auth.Form.FormXPath != nil {
			form.FormXPath = *auth.Form.FormXPath
		}
		if auth.Form.LoginKey != nil {
			form.LoginKey = *auth.Form.LoginKey
		}
		if auth.Form.PasswordKey != nil {
			form.PasswordKey = *auth.Form.PasswordKey
		}
		result.Form = form
	}

	if auth.Http != nil {
		result.Http = &BlackBoxHTTPAuth{
			Login:             auth.Http.Login,
			Password:          auth.Http.Password,
			ValidationAddress: auth.Http.ValidationAddress,
		}
	}

	return result
}

func mapAuthenticationV110(auth *v110.AIProjBlackBoxSettingsAuthentication) *BlackBoxAuthentication {
	if auth == nil {
		return nil
	}

	result := &BlackBoxAuthentication{}
	if auth.Type != nil {
		result.Type = string(*auth.Type)
	}

	if auth.Cookie != nil {
		result.Cookie = &BlackBoxCookieAuth{
			Cookie:             auth.Cookie.Cookie,
			ValidationAddress:  auth.Cookie.ValidationAddress,
			ValidationTemplate: auth.Cookie.ValidationTemplate,
		}
	}

	if auth.Form != nil {
		form := &BlackBoxFormAuth{
			FormDetection:      string(auth.Form.FormDetection),
			FormAddress:        auth.Form.FormAddress,
			Login:              auth.Form.Login,
			Password:           auth.Form.Password,
			ValidationTemplate: auth.Form.ValidationTemplate,
		}
		if auth.Form.FormXPath != nil {
			form.FormXPath = *auth.Form.FormXPath
		}
		if auth.Form.LoginKey != nil {
			form.LoginKey = *auth.Form.LoginKey
		}
		if auth.Form.PasswordKey != nil {
			form.PasswordKey = *auth.Form.PasswordKey
		}
		result.Form = form
	}

	if auth.Http != nil {
		result.Http = &BlackBoxHTTPAuth{
			Login:             auth.Http.Login,
			Password:          auth.Http.Password,
			ValidationAddress: auth.Http.ValidationAddress,
		}
	}

	return result
}

func mapProxySettingsV19(proxy *v19.AIProjBlackBoxSettingsProxySettings) *BlackBoxProxySettings {
	if proxy == nil {
		return nil
	}

	result := &BlackBoxProxySettings{}
	if proxy.Enabled != nil {
		result.Enabled = *proxy.Enabled
	}
	if proxy.Host != nil {
		result.Host = *proxy.Host
	}
	if proxy.Login != nil {
		result.Login = *proxy.Login
	}
	if proxy.Password != nil {
		result.Password = *proxy.Password
	}
	if proxy.Port != nil {
		result.Port = *proxy.Port
	}
	if proxy.Type != nil {
		result.Type = string(*proxy.Type)
	}

	return result
}

func mapProxySettingsV110(proxy *v110.AIProjBlackBoxSettingsProxySettings) *BlackBoxProxySettings {
	if proxy == nil {
		return nil
	}

	result := &BlackBoxProxySettings{}
	if proxy.Enabled != nil {
		result.Enabled = *proxy.Enabled
	}
	if proxy.Host != nil {
		result.Host = *proxy.Host
	}
	if proxy.Login != nil {
		result.Login = *proxy.Login
	}
	if proxy.Password != nil {
		result.Password = *proxy.Password
	}
	if proxy.Port != nil {
		result.Port = *proxy.Port
	}
	if proxy.Type != nil {
		result.Type = string(*proxy.Type)
	}

	return result
}
