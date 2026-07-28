package v5_x

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/pkg/clientai/v5_x"
	"github.com/google/uuid"
)

func (a *ClientAI5x) setBlackBoxSettings(ctx context.Context, projectId uuid.UUID, scanSettings *settings.ScanSettings) error {
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

func toBlackBoxSettingsModel(payload common.BlackBoxPayload) v5_x.BlackBoxSettingsModel {
	model := v5_x.BlackBoxSettingsModel{
		Site:                  payload.Site,
		Level:                 toBlackBoxScanLevel(payload.Level),
		ScanScope:             toScanScope(payload.ScanScope),
		SslCheck:              payload.SslCheck,
		RunAutocheckAfterScan: payload.RunAutocheckAfterScan,
		IsActive:              payload.IsActive,
	}

	if payload.AdditionalHttpHeaders != nil {
		headers := make([]v5_x.HttpHeaderModel, len(*payload.AdditionalHttpHeaders))
		for i, header := range *payload.AdditionalHttpHeaders {
			headers[i] = v5_x.HttpHeaderModel{
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

func toBlackBoxAddresses(source *[]common.BlackBoxAddress) *[]v5_x.BlackBoxAddressModel {
	if source == nil {
		return nil
	}

	result := make([]v5_x.BlackBoxAddressModel, len(*source))
	for i, address := range *source {
		result[i] = v5_x.BlackBoxAddressModel{
			Address: address.Address,
			Format:  toBlackBoxFormat(address.Format),
		}
	}

	return &result
}

func toBlackBoxAuthentication(source *common.BlackBoxAuth) *v5_x.BlackBoxAuthenticationFullModel {
	if source == nil {
		return nil
	}

	result := &v5_x.BlackBoxAuthenticationFullModel{}
	if source.Cookie != nil {
		result.Cookie = &v5_x.BlackBoxRawCookieAuthenticationModel{
			Cookie:             source.Cookie.Cookie,
			ValidationAddress:  source.Cookie.ValidationAddress,
			ValidationTemplate: source.Cookie.ValidationTemplate,
		}
	}
	if source.Form != nil {
		result.Form = &v5_x.BlackBoxFormAuthenticationModel{
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
		result.Http = &v5_x.BlackBoxHttpAuthenticationModel{
			Login:             source.Http.Login,
			Password:          source.Http.Password,
			ValidationAddress: source.Http.ValidationAddress,
		}
	}

	return result
}

func toBlackBoxProxySettings(source *common.BlackBoxProxy) *v5_x.BlackBoxProxySettingsModel {
	if source == nil {
		return nil
	}

	var port *int32
	if source.Port != nil {
		value := int32(*source.Port)
		port = &value
	}

	return &v5_x.BlackBoxProxySettingsModel{
		IsActive: source.Enabled,
		Host:     source.Host,
		Login:    source.Login,
		Password: source.Password,
		Port:     port,
		Type:     toProxyType(source.Type),
	}
}

func toBlackBoxScanLevel(value *string) *v5_x.BlackBoxScanLevel {
	if value == nil || *value == "" {
		return nil
	}

	level := v5_x.BlackBoxScanLevel(*value)

	return &level
}

func toScanScope(value *string) *v5_x.ScanScope {
	if value == nil || *value == "" {
		return nil
	}

	scope := v5_x.ScanScope(*value)

	return &scope
}

func toBlackBoxFormat(value *string) *v5_x.BlackBoxFormat {
	if value == nil || *value == "" {
		return nil
	}

	format := v5_x.BlackBoxFormat(*value)

	return &format
}

func toBlackBoxFormDetection(value *string) *v5_x.BlackBoxFormDetection {
	if value == nil || *value == "" {
		return nil
	}

	detection := v5_x.BlackBoxFormDetection(*value)

	return &detection
}

func toProxyType(value *string) *v5_x.ProxyType {
	if value == nil || *value == "" {
		return nil
	}

	proxyType := v5_x.ProxyType(*value)

	return &proxyType
}
