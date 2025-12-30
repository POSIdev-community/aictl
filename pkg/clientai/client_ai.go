package clientai

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
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
	err := a.initialize(cfg)
	if err != nil {
		return fmt.Errorf("initialize: %w", err)
	}

	if err := a.updateAccessTokenJwt(ctx, cfg); err != nil {
		return fmt.Errorf("update jwt: %w", err)
	}

	return nil
}

func (a *AiClient) initialize(cfg *config.Config) error {
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

	return nil
}

func (a *AiClient) InitializeLikeUser(ctx context.Context, cfg *config.Config, username, password string) error {
	err := a.initialize(cfg)
	if err != nil {
		return fmt.Errorf("initialize: %w", err)
	}

	err = a.userLogin(ctx, cfg, username, password)
	if err != nil {
		return fmt.Errorf("user login: %w", err)
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

func (a *AiClient) updateAccessTokenJwt(ctx context.Context, cfg *config.Config) error {
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

func (a *AiClient) userLogin(ctx context.Context, cfg *config.Config, login, password string) error {
	body := userLoginRequest{
		Login:    login,
		Password: password,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	url := cfg.UriString() + "/api/auth/userLogin?scopeType=Web"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.jwtHttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	var loginResp userLoginResponse
	if err := json.Unmarshal(respBody, &loginResp); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	if loginResp.AccessToken == "" || loginResp.RefreshToken == "" {
		return fmt.Errorf("no access token in response")
	}

	a.accessToken = loginResp.AccessToken
	a.refreshToken = loginResp.RefreshToken

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
