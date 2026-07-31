package common

import (
	"context"
	"net/http"
	"sync"

	"golang.org/x/sync/singleflight"
)

type BaseClient struct {
	HttpClient    *http.Client
	JwtHttpClient *http.Client

	tokenMu      sync.RWMutex
	apiToken     string
	accessToken  string
	refreshToken string

	Initialized bool
	WithRetry   bool

	jwtRefresh singleflight.Group
}

func NewBaseClient() *BaseClient {
	return &BaseClient{
		HttpClient:    &http.Client{},
		JwtHttpClient: &http.Client{},
	}
}

func (c *BaseClient) Reset() {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()

	c.HttpClient = &http.Client{}
	c.JwtHttpClient = &http.Client{}
	c.apiToken = ""
	c.accessToken = ""
	c.refreshToken = ""
	c.Initialized = false
	c.WithRetry = false
	c.jwtRefresh = singleflight.Group{}
}

func (c *BaseClient) GetAPIToken() string {
	c.tokenMu.RLock()
	defer c.tokenMu.RUnlock()

	return c.apiToken
}

func (c *BaseClient) SetAPIToken(token string) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()
	c.apiToken = token
}

func (c *BaseClient) GetAccessToken() string {
	c.tokenMu.RLock()
	defer c.tokenMu.RUnlock()

	return c.accessToken
}

func (c *BaseClient) SetAccessToken(token string) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()
	c.accessToken = token
}

func (c *BaseClient) GetRefreshToken() string {
	c.tokenMu.RLock()
	defer c.tokenMu.RUnlock()

	return c.refreshToken
}

func (c *BaseClient) SetAuthTokens(accessToken, refreshToken string) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()
	c.accessToken = accessToken
	c.refreshToken = refreshToken
}

// DoJWTRefresh runs fn once for concurrent callers; others wait for the same result.
func (c *BaseClient) DoJWTRefresh(fn func() error) error {
	_, err, _ := c.jwtRefresh.Do("jwt", func() (any, error) {
		return nil, fn()
	})

	return err
}

func (a *BaseClient) AddJWTToHeader(_ context.Context, req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+a.GetAccessToken())

	return nil
}
