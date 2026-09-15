package v6_x

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/apperror"
)

func TestCheckResponse_IncludesBodyOnBadRequest(t *testing.T) {
	t.Parallel()

	body := `{"Message":"Invalid policy field"}`
	rsp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(strings.NewReader(body)),
	}

	err := CheckResponse(rsp, "policies")
	require.Error(t, err)

	var badReq *apperror.BadRequestError
	require.True(t, errors.As(err, &badReq))
	require.Equal(t, "Bad Request error", badReq.Error())
	require.Equal(t, body, badReq.Body())
}

func TestCheckResponse_Nil(t *testing.T) {
	t.Parallel()

	err := CheckResponse(nil, "policies")
	var empty *apperror.EmptyResponseError
	require.ErrorAs(t, err, &empty)
}
