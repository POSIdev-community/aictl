package v6_0

import (
	"net/http"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	"github.com/POSIdev-community/aictl/internal/core/apperror"
	"github.com/POSIdev-community/aictl/pkg/clientai/v6_0"
)

func CheckResponse(rsp *http.Response, resourceName string) error {
	if rsp == nil {
		return apperror.NewEmptyResponseError(resourceName)
	}

	if rsp.StatusCode < 400 {
		return nil
	}

	body := common.ReadErrorBody(rsp.Body)

	return checkResponseCommon(rsp.StatusCode, body, nil, resourceName)
}

func CheckResponseByModel(statusCode int, body string, model *v6_0.ApiErrorModel) error {
	if model != nil && model.ErrorCode != nil {
		details := map[string]*string{}
		if model.Details != nil {
			details = *model.Details
		}

		return apperror.CheckApiErrorModel(statusCode, string(*model.ErrorCode), details)
	}

	return checkResponseCommon(statusCode, body, nil, "")
}

func checkResponseCommon(statusCode int, body string, model *v6_0.ApiErrorModel, resourceName string) error {
	if statusCode < 400 {
		return nil
	}

	if statusCode >= 500 {
		return apperror.NewServerError(statusCode, body)
	}

	if model != nil {
		var errorCode = string(*model.ErrorCode)

		return apperror.CheckApiErrorModel(statusCode, errorCode, *model.Details)
	}

	return apperror.CheckResponseErrors(statusCode, body, resourceName)
}
