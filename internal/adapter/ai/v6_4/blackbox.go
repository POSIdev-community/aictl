package v6_4

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/pkg/clientai/v6_4"
)

func (a *ClientAI64) setBlackBoxSettings(ctx context.Context, projectId uuid.UUID, scanSettings *settings.ScanSettings) error {
	payload := common.BuildBlackBoxPayload(scanSettings.BlackBoxSettings, scanSettings.BlackBoxEnabled)
	body := toBlackBoxSettingsModel(payload)

	res, err := a.PutApiProjectsProjectIdBlackBoxSettingsWithResponse(ctx, projectId, body, a.AddJWTToHeader)
	if err != nil {
		return fmt.Errorf("put black box settings request: %w", err)
	}

	statusCode := res.StatusCode()
	responseBody := string(res.Body)
	if err = CheckResponseByModel(statusCode, responseBody, nil); err != nil {
		return fmt.Errorf("put black box settings: %w", err)
	}

	return nil
}

func toBlackBoxSettingsModel(payload common.BlackBoxPayload) v6_4.BlackBoxSettingsModel {
	model := v6_4.BlackBoxSettingsModel{
		Site:                  payload.Site,
		Level:                 toBlackBoxScanLevel(payload.Level),
		ScanScope:             toScanScope(payload.ScanScope),
		SslCheck:              payload.SslCheck,
		RunAutocheckAfterScan: payload.RunAutocheckAfterScan,
		IsActive:              payload.IsActive,
	}

	if payload.AdditionalHttpHeaders != nil {
		headers := make([]v6_4.HttpHeaderModel, len(*payload.AdditionalHttpHeaders))
		for i, header := range *payload.AdditionalHttpHeaders {
			headers[i] = v6_4.HttpHeaderModel{
				Key:   header.Key,
				Value: header.Value,
			}
		}
		model.AdditionalHttpHeaders = &headers
	}

	model.WhiteListedAddresses = toBlackBoxAddresses(payload.WhiteListedAddresses)
	model.BlackListedAddresses = toBlackBoxAddresses(payload.BlackListedAddresses)
	model.Authentication = toBlackBoxAuthentication(payload.Authentication)
	model.ProxySettings = toBlackBoxProxySettings(payload.ProxySettings)

	return model
}

func toBlackBoxAddresses(source *[]common.BlackBoxAddress) *[]v6_4.BlackBoxAddressModel {
	if source == nil {
		return nil
	}

	result := make([]v6_4.BlackBoxAddressModel, len(*source))
	for i, address := range *source {
		result[i] = v6_4.BlackBoxAddressModel{
			Address: address.Address,
			Format:  toBlackBoxFormat(address.Format),
		}
	}

	return &result
}

func toBlackBoxAuthentication(source *common.BlackBoxAuth) *v6_4.BlackBoxAuthenticationFullModel {
	if source == nil {
		return nil
	}

	result := &v6_4.BlackBoxAuthenticationFullModel{}
	if source.Cookie != nil {
		result.Cookie = &v6_4.BlackBoxRawCookieAuthenticationModel{
			Cookie:             source.Cookie.Cookie,
			ValidationAddress:  source.Cookie.ValidationAddress,
			ValidationTemplate: source.Cookie.ValidationTemplate,
		}
	}
	if source.Form != nil {
		result.Form = &v6_4.BlackBoxFormAuthenticationModel{
			FormDetection:      toBlackBoxFormDetection(source.Form.FormDetection),
			FormAddress:        source.Form.FormAddress,
			FormXPath:          source.Form.FormXPath,
			Login:              source.Form.Login,
			LoginKey:           source.Form.LoginKey,
			Password:           source.Form.Password,
			PasswordKey:        source.Form.PasswordKey,
			ValidationTemplate: source.Form.ValidationTemplate,
		}
	}
	if source.Http != nil {
		result.Http = &v6_4.BlackBoxHttpAuthenticationModel{
			Login:             source.Http.Login,
			Password:          source.Http.Password,
			ValidationAddress: source.Http.ValidationAddress,
		}
	}

	return result
}

func toBlackBoxProxySettings(source *common.BlackBoxProxy) *v6_4.BlackBoxProxySettingsModel {
	if source == nil {
		return nil
	}

	var port *int32
	if source.Port != nil {
		value := int32(*source.Port)
		port = &value
	}

	return &v6_4.BlackBoxProxySettingsModel{
		IsActive: source.Enabled,
		Host:     source.Host,
		Login:    source.Login,
		Password: source.Password,
		Port:     port,
		Type:     toProxyType(source.Type),
	}
}

func toBlackBoxScanLevel(value *string) *v6_4.BlackBoxScanLevel {
	if value == nil || *value == "" {
		return nil
	}

	level := v6_4.BlackBoxScanLevel(*value)

	return &level
}

func toScanScope(value *string) *v6_4.ScanScope {
	if value == nil || *value == "" {
		return nil
	}

	scope := v6_4.ScanScope(*value)

	return &scope
}

func toBlackBoxFormat(value *string) *v6_4.BlackBoxFormat {
	if value == nil || *value == "" {
		return nil
	}

	format := v6_4.BlackBoxFormat(*value)

	return &format
}

func toBlackBoxFormDetection(value *string) *v6_4.BlackBoxFormDetection {
	if value == nil || *value == "" {
		return nil
	}

	detection := v6_4.BlackBoxFormDetection(*value)

	return &detection
}

func toProxyType(value *string) *v6_4.ProxyType {
	if value == nil || *value == "" {
		return nil
	}

	proxyType := v6_4.ProxyType(*value)

	return &proxyType
}
