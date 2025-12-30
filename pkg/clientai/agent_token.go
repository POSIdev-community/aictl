package clientai

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// userLoginRequest represents the request body for user login
type userLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// userLoginResponse represents the response from user login
type userLoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// createTokenRequest represents the request body for creating access token
type createTokenRequest struct {
	Name   string   `json:"name"`
	Scopes []string `json:"scopes"`
}

// createTokenResponse represents the response from creating access token
type createTokenResponse struct {
	Token string `json:"token"`
}

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

	// Step 1: Login with user credentials to get admin access token
	adminToken, err := userLogin(ctx, httpClient, serverURL, login, password)
	if err != nil {
		return "", fmt.Errorf("user login: %w", err)
	}

	// Step 2: Create agent token with ScanAgent scope
	agentToken, err := createAccessToken(ctx, httpClient, serverURL, adminToken, agentName)
	if err != nil {
		return "", fmt.Errorf("create access token: %w", err)
	}

	return agentToken, nil
}

func userLogin(ctx context.Context, client *http.Client, serverURL, login, password string) (string, error) {
	body := userLoginRequest{
		Login:    login,
		Password: password,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	url := serverURL + "/api/auth/userLogin?scopeType=Web"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	var loginResp userLoginResponse
	if err := json.Unmarshal(respBody, &loginResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if loginResp.AccessToken == "" {
		return "", fmt.Errorf("no access token in response")
	}

	return loginResp.AccessToken, nil
}

func createAccessToken(ctx context.Context, client *http.Client, serverURL, adminToken, agentName string) (string, error) {
	body := createTokenRequest{
		Name:   agentName,
		Scopes: []string{"ScanAgent"},
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	url := serverURL + "/api/auth/accessToken"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	var tokenResp createTokenResponse
	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if tokenResp.Token == "" {
		return "", fmt.Errorf("no token in response")
	}

	return tokenResp.Token, nil
}
