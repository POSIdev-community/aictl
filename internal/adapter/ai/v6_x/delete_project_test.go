package v6_x

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/apperror"
	clientai "github.com/POSIdev-community/aictl/pkg/clientai/v6_x"
)

func TestDeleteProjectNotFound(t *testing.T) {
	t.Parallel()

	projectID := uuid.MustParse("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
	client := newTestClient6x(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/projects/"+projectID.String(), r.URL.Path)

		code := clientai.ApiErrorTypePROJECTNOTFOUND
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(clientai.ApiErrorModel{ErrorCode: &code})
	}))

	err := client.DeleteProject(t.Context(), projectID)
	require.Error(t, err)

	var notFound *apperror.NotFoundByIdError
	require.True(t, errors.As(err, &notFound))
	assert.Equal(t, "project", notFound.Resource)
	assert.Equal(t, projectID.String(), notFound.ID)
	assert.Equal(t, "project with ID "+projectID.String()+" not found", notFound.Error())
}

func TestDeleteProjectSuccess(t *testing.T) {
	t.Parallel()

	projectID := uuid.MustParse("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
	client := newTestClient6x(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/projects/"+projectID.String(), r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))

	require.NoError(t, client.DeleteProject(t.Context(), projectID))
}
