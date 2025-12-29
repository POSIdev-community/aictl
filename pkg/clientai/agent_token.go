package clientai

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
)

// CreateAgentToken authenticates with user credentials and creates an access token
// for a scan agent. This is a standalone function because it uses a different
// authentication flow (user/password) than the regular API client (API token).
func CreateAgentToken(ctx context.Context, serverURL, login, password, agentName string, tlsSkip bool) (string, error) {
	httpClient := &http.Client{}
	if tlsSkip {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	client, err := NewClientWithResponses(serverURL, WithHTTPClient(httpClient))
	if err != nil {
		return "", fmt.Errorf("new client: %w", err)
	}

	// Step 1: Login with user credentials to get admin access token
	scopeType := AuthScopeWeb
	loginParams := &PostApiAuthUserLoginParams{
		ScopeType: &scopeType,
	}
	loginBody := PostApiAuthUserLoginJSONRequestBody{
		Login:    &login,
		Password: &password,
	}

	loginResp, err := client.PostApiAuthUserLoginWithResponse(ctx, loginParams, loginBody)
	if err != nil {
		return "", fmt.Errorf("user login: %w", err)
	}

	if err = CheckResponseByModel(loginResp.StatusCode(), string(loginResp.Body), loginResp.JSON400); err != nil {
		return "", fmt.Errorf("user login: %w", err)
	}

	if loginResp.JSON200 == nil || loginResp.JSON200.AccessToken == nil {
		return "", fmt.Errorf("user login: no access token in response")
	}

	adminToken := *loginResp.JSON200.AccessToken

	// Step 2: Create agent token with ScanAgent scope
	scopes := []AccessTokenScopeType{AccessTokenScopeTypeScanAgent}
	tokenBody := PostApiAuthAccessTokenJSONRequestBody{
		Name:   &agentName,
		Scopes: &scopes,
	}

	tokenResp, err := client.PostApiAuthAccessTokenWithResponse(ctx, tokenBody, func(_ context.Context, req *http.Request) error {
		req.Header.Add("Authorization", "Bearer "+adminToken)
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("create access token: %w", err)
	}

	if err = CheckResponseByModel(tokenResp.StatusCode(), string(tokenResp.Body), tokenResp.JSON400); err != nil {
		return "", fmt.Errorf("create access token: %w", err)
	}

	if tokenResp.JSON200 == nil || tokenResp.JSON200.Token == nil {
		return "", fmt.Errorf("create access token: no token in response")
	}

	return *tokenResp.JSON200.Token, nil
}
