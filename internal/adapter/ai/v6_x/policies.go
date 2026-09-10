package v6_x

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/apperror"
	"github.com/POSIdev-community/aictl/internal/core/domain/policystate"
)

func (a *ClientAI6x) GetProjectPolicies(ctx context.Context, projectId uuid.UUID) (io.ReadCloser, error) {
	response, err := a.GetApiProjectsProjectIdSecurityPolicies(ctx, projectId, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get project policies request: %w", err)
	}

	if err = CheckResponse(response, "policies"); err != nil {
		_ = response.Body.Close()

		return nil, fmt.Errorf("ai adapter get project policies: %w", err)
	}

	return response.Body, nil
}

func (a *ClientAI6x) SetProjectPolicies(ctx context.Context, projectId uuid.UUID, rawJSON []byte) error {
	response, err := a.PutApiProjectsProjectIdSecurityPoliciesWithBody(
		ctx, projectId, "application/json", bytes.NewReader(rawJSON), a.AddJWTToHeader,
	)
	if err != nil {
		return fmt.Errorf("ai adapter set project policies request: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if err = CheckResponse(response, "policies"); err != nil {
		return fmt.Errorf("ai adapter set project policies: %w", err)
	}

	return nil
}

func (a *ClientAI6x) GetScanPolicyState(ctx context.Context, projectId, scanId uuid.UUID) (policystate.State, error) {
	response, err := a.GetApiProjectsProjectIdScanResultsScanResultIdStatisticWithResponse(
		ctx, projectId, scanId, a.AddJWTToHeader,
	)
	if err != nil {
		return "", fmt.Errorf("ai adapter get scan policy request: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	if err = CheckResponseByModel(statusCode, body, response.JSON400); err != nil {
		return "", fmt.Errorf("ai adapter get scan policy: %w", err)
	}

	if response.JSON200 == nil {
		return "", apperror.NewEmptyResponseError("scan statistic")
	}

	// On AIE 6.3+ policy state moved under securityPolicyStatus. When policies were
	// not evaluated, the API returns securityPolicyStatus: null (previously top-level
	// policyState: "None").
	if response.JSON200.SecurityPolicyStatus == nil ||
		response.JSON200.SecurityPolicyStatus.PolicyState == nil {
		return policystate.None, nil
	}

	return policystate.State(*response.JSON200.SecurityPolicyStatus.PolicyState), nil
}
