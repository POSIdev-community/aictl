package clientai

import (
	"context"
	"crypto/tls"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"net/http"
)

type AiClient struct {
	*ClientWithResponses

	jwt       string
	bearerJwt string
}

func NewAiClient(ctx context.Context, cfg *config.Config) (*AiClient, error) {
	httpClient := &http.Client{}
	if cfg.TLSSkip() {
		httpClient.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}} //nolint:gosec // switch for test
	}

	client, err := NewClientWithResponses(cfg.UriString(), WithHTTPClient(httpClient))

	jwt, err := getJwt(ctx, cfg, client)
	if err != nil {
		return nil, err
	}

	return &AiClient{
		client,
		jwt,
		"Bearer " + jwt,
	}, nil
}

func (a *AiClient) AddJwtToHeader(_ context.Context, req *http.Request) error {
	req.Header.Add("Authorization", a.bearerJwt)

	return nil
}

func getJwt(ctx context.Context, cfg *config.Config, client *ClientWithResponses) (string, error) {
	response, err := client.GetApiAuthSigninWithResponse(ctx, func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Access-Token", cfg.Token())

		return nil
	})
	if err != nil {
		return "", err
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	model := response.JSON400
	if err = CheckResponseByModel(statusCode, body, model); err != nil {
		return "", err
	}

	return *response.JSON200.AccessToken, nil
}
