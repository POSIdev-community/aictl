package common

import (
	"context"
	"net/http"

	"golang.org/x/sync/singleflight"
)

type BaseClient struct {
	HttpClient    *http.Client
	JwtHttpClient *http.Client

	AccessToken  string
	RefreshToken string
	Initialized  bool
	WithRetry    bool

	jwtRefresh singleflight.Group
}

func NewBaseClient() *BaseClient {
	return &BaseClient{
		HttpClient:    &http.Client{},
		JwtHttpClient: &http.Client{},
	}
}

func (c *BaseClient) Reset() {
	c.HttpClient = &http.Client{}
	c.JwtHttpClient = &http.Client{}
	c.AccessToken = ""
	c.RefreshToken = ""
	c.Initialized = false
	c.WithRetry = false
	c.jwtRefresh = singleflight.Group{}
}

// DoJWTRefresh runs fn once for concurrent callers; others wait for the same result.
func (c *BaseClient) DoJWTRefresh(fn func() error) error {
	_, err, _ := c.jwtRefresh.Do("jwt", func() (any, error) {
		return nil, fn()
	})

	return err
}

func (a *BaseClient) AddJWTToHeader(_ context.Context, req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+a.AccessToken)

	return nil
}
