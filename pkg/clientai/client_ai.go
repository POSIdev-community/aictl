package clientai

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/logger"
)

type AiClient struct {
	*ClientWithResponses
	jwtClient *ClientWithResponses

	httpClient    *http.Client
	jwtHttpClient *http.Client

	accessToken  string
	refreshToken string
}

func NewAiClient() (*AiClient, error) {
	httpClient := &http.Client{}
	jwtHTTPClient := &http.Client{}

	aiClient := AiClient{
		httpClient:    httpClient,
		jwtHttpClient: jwtHTTPClient,
	}

	return &aiClient, nil
}

func (a *AiClient) Initialize(ctx context.Context, cfg *config.Config) error {
	client, err := NewClientWithResponses(cfg.UriString(), WithHTTPClient(a.httpClient))
	if err != nil {
		return fmt.Errorf("new client: %w", err)
	}
	a.ClientWithResponses = client

	a.jwtClient, err = NewClientWithResponses(cfg.UriString(), WithHTTPClient(a.jwtHttpClient))
	if err != nil {
		return fmt.Errorf("new jwt client: %w", err)
	}

	transport := &http.Transport{}
	if cfg.TLSSkip() {
		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{}
		}

		transport.TLSClientConfig.InsecureSkipVerify = true
	}

	a.httpClient.Transport = transport.Clone()
	a.jwtHttpClient.Transport = transport.Clone()

	if err := a.getJWT(ctx, cfg); err != nil {
		return fmt.Errorf("update jwt: %w", err)
	}

	return nil
}

func (a *AiClient) AddJwtRetry() {
	a.httpClient.Transport = NewRetryRoundTripper(a.httpClient.Transport, http.StatusUnauthorized, a.refreshJWT)
}

func (a *AiClient) AddJWTToHeader(_ context.Context, req *http.Request) error {
	req.Header.Add("Authorization", "Bearer "+a.accessToken)

	return nil
}

func (a *AiClient) getJWT(ctx context.Context, cfg *config.Config) error {
	response, err := a.jwtClient.GetApiAuthSigninWithResponse(ctx, func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Access-Token", cfg.Token())

		return nil
	})
	if err != nil {
		return fmt.Errorf("get api auth signin: %w", err)
	}

	if err = CheckResponseByModel(response.StatusCode(), string(response.Body), response.JSON400); err != nil {
		return err
	}

	a.accessToken = *response.JSON200.AccessToken
	a.refreshToken = *response.JSON200.RefreshToken

	return nil
}

func (a *AiClient) refreshJWT(ctx context.Context, req *http.Request) error {
	log := logger.FromContext(ctx)

	response, err := a.jwtClient.GetApiAuthRefreshTokenWithResponse(ctx, func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", "Bearer "+a.refreshToken)

		return nil
	})
	if err != nil {
		return fmt.Errorf("get api auth signin: %w", err)
	}

	if err = CheckResponse(response.HTTPResponse, "jwt refresh"); err != nil {
		return err
	}

	if response.JSON200.AccessToken == nil {
		log.StdErrf("Got empty access token")

		return fmt.Errorf("no access token")
	}

	a.accessToken = *response.JSON200.AccessToken

	req.Header.Set("Authorization", "Bearer "+a.accessToken)

	return nil
}

func (a *AiClient) CreateAgentToken(ctx context.Context, serverURL, login, password, agentName string, tlsSkip bool) (string, error) {
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
