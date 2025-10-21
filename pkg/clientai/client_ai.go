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

	accessToken  string
	refreshToken string
}

func NewAiClient(ctx context.Context, cfg *config.Config) (*AiClient, error) {
	httpClient := &http.Client{}
	jwtHTTPClient := &http.Client{}

	transport := &http.Transport{}
	if cfg.TLSSkip() {
		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{}
		}

		transport.TLSClientConfig.InsecureSkipVerify = true
	}

	httpClient.Transport = transport.Clone()
	jwtHTTPClient.Transport = transport.Clone()

	client, err := NewClientWithResponses(cfg.UriString(), WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("new client: %w", err)
	}

	jwtClient, err := NewClientWithResponses(cfg.UriString(), WithHTTPClient(jwtHTTPClient))
	if err != nil {
		return nil, fmt.Errorf("new jwt client: %w", err)
	}

	aiClient := AiClient{
		ClientWithResponses: client,
		jwtClient:           jwtClient,
	}

	httpClient.Transport = NewRetryRoundTripper(httpClient.Transport, http.StatusUnauthorized, aiClient.refreshJWT)

	if err := aiClient.getJWT(ctx, cfg, jwtClient); err != nil {
		return nil, fmt.Errorf("update jwt: %w", err)
	}

	return &aiClient, nil
}

func (a *AiClient) AddJWTToHeader(_ context.Context, req *http.Request) error {
	req.Header.Add("Authorization", "Bearer "+a.accessToken)

	return nil
}

func (a *AiClient) getJWT(ctx context.Context, cfg *config.Config, client *ClientWithResponses) error {
	response, err := client.GetApiAuthSigninWithResponse(ctx, func(ctx context.Context, req *http.Request) error {
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
	log := logger.FromContext(ctx).Named("ai_client.refresh_jwt")
	log.Debug("Refreshing JWT")

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
		log.Error("Got empty access token")

		return fmt.Errorf("no access token")
	}

	a.accessToken = *response.JSON200.AccessToken

	req.Header.Set("Authorization", "Bearer "+a.accessToken)

	log.Debug("JWT refreshed")

	return nil
}
