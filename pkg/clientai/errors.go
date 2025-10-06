package clientai

import (
	"github.com/POSIdev-community/aictl/pkg/errs"
	"io"
	"net/http"
)

func CheckResponse(rsp *http.Response, resourceName string) error {
	if rsp == nil {
		return errs.NewNilResponseError()
	}

	if rsp.StatusCode < 400 {
		return nil
	}

	var bytes []byte
	var body string
	_, err := rsp.Body.Read(bytes)
	if err != nil {
		if err != io.EOF {
			body = ""
		}
	}

	return checkResponseCommon(rsp.StatusCode, body, nil, resourceName)
}

func CheckResponseByModel(statusCode int, body string, model *ApiErrorModel) error {
	if model != nil && model.ErrorCode != nil && model.Details != nil {
		var errorCode = string(*model.ErrorCode)

		return errs.CheckApiErrorModel(statusCode, errorCode, *model.Details)
	}

	return checkResponseCommon(statusCode, body, nil, "")
}

func checkResponseCommon(statusCode int, body string, model *ApiErrorModel, resourceName string) error {
	if statusCode < 400 {
		return nil
	}

	if statusCode >= 500 {
		return errs.NewServerError(statusCode, body)
	}

	if model != nil {
		var errorCode = string(*model.ErrorCode)

		return errs.CheckApiErrorModel(statusCode, errorCode, *model.Details)
	}

	return errs.CheckResponseErrors(statusCode, body, resourceName)
}
