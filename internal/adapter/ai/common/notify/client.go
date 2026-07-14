package notify

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

const (
	hubPath           = "/notifyApi/notifications"
	recordSeparator   = '\x1e'
	handshakeTimeout  = 15 * time.Second
	defaultPingPeriod = 15 * time.Second
)

// Client connects to PT AI notification hub via ASP.NET Core SignalR JSON protocol.
type Client struct {
	baseURL     string
	accessToken string
	httpClient  *http.Client
	tlsSkip     bool
	clientID    uuid.UUID

	mu sync.RWMutex
}

type Options struct {
	BaseURL     string
	AccessToken string
	HTTPClient  *http.Client
	TLSSkip     bool
	ClientID    uuid.UUID
}

func NewClient(opts Options) (*Client, error) {
	if opts.BaseURL == "" {
		return nil, fmt.Errorf("base url is required")
	}
	if opts.AccessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}

	httpClient := opts.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: handshakeTimeout}
	}

	clientID := opts.ClientID
	if clientID == uuid.Nil {
		clientID = uuid.New()
	}

	return &Client{
		baseURL:     strings.TrimRight(opts.BaseURL, "/"),
		accessToken: opts.AccessToken,
		httpClient:  httpClient,
		tlsSkip:     opts.TLSSkip,
		clientID:    clientID,
	}, nil
}

// SetAccessToken updates the JWT used for negotiate and websocket dial (e.g. after refresh).
func (c *Client) SetAccessToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.accessToken = token
}

func (c *Client) getAccessToken() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.accessToken
}

// Subscribe opens a WebSocket to the notification hub and returns a channel of hub messages.
// The channel is closed when ctx is cancelled, the connection ends, or a fatal error occurs.
// non-nil errorC receives the terminal error (if any) once; caller may ignore it.
func (c *Client) Subscribe(ctx context.Context) (<-chan Message, <-chan error, error) {
	connToken, err := c.negotiate(ctx)
	if err != nil {
		return nil, nil, err
	}

	wsURL, err := c.websocketURL(connToken)
	if err != nil {
		return nil, nil, err
	}

	dialOpts := &websocket.DialOptions{
		HTTPClient: c.httpClient,
	}
	if c.tlsSkip {
		dialOpts.HTTPClient = cloneHTTPClientWithTLSSkip(c.httpClient)
	}

	conn, resp, err := websocket.Dial(ctx, wsURL, dialOpts)
	if err != nil {
		if resp != nil && isUnauthorized(resp.StatusCode) {
			return nil, nil, newAuthError(resp.StatusCode, "")
		}

		return nil, nil, fmt.Errorf("websocket dial: %w", err)
	}
	conn.SetReadLimit(1 << 20)

	if err := c.handshake(ctx, conn); err != nil {
		_ = conn.Close(websocket.StatusInternalError, "handshake failed")

		return nil, nil, err
	}

	if err := c.subscribeScanNotifications(ctx, conn); err != nil {
		_ = conn.Close(websocket.StatusInternalError, "subscribe failed")

		return nil, nil, err
	}

	out := make(chan Message, 16)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)
		c.readLoop(ctx, conn, out, errc)
	}()

	return out, errc, nil
}

// subscriptionOnNotification mirrors ptai-ee-tools SubscriptionOnNotification hub payload.
type subscriptionOnNotification struct {
	NotificationTypeName string      `json:"notificationTypeName"`
	IDs                  []uuid.UUID `json:"ids"`
	CreatedDate          time.Time   `json:"createdDate"`
}

func (c *Client) subscribeScanNotifications(ctx context.Context, conn *websocket.Conn) error {
	// Match generic-client-lib ApiClient.subscribe: fire-and-forget SubscribeOnNotification
	// for ScanProgress and ScanCompleted (ids empty = all events of that type).
	for _, typ := range []string{TargetScanProgress, TargetScanCompleted} {
		if err := c.subscribeOnNotification(ctx, conn, typ, nil); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) subscribeOnNotification(ctx context.Context, conn *websocket.Conn, notificationType string, ids []uuid.UUID) error {
	writeCtx, cancel := context.WithTimeout(ctx, handshakeTimeout)
	defer cancel()

	if ids == nil {
		ids = []uuid.UUID{}
	}

	payload, err := json.Marshal(map[string]any{
		"type":   messageTypeInvocation,
		"target": "SubscribeOnNotification",
		"arguments": []any{subscriptionOnNotification{
			NotificationTypeName: notificationType,
			IDs:                  ids,
			CreatedDate:          time.Now().UTC(),
		}},
	})
	if err != nil {
		return fmt.Errorf("marshal SubscribeOnNotification: %w", err)
	}
	payload = append(payload, recordSeparator)

	if err := conn.Write(writeCtx, websocket.MessageText, payload); err != nil {
		return fmt.Errorf("write SubscribeOnNotification %s: %w", notificationType, err)
	}

	return nil
}

type negotiateResponse struct {
	ConnectionID    string `json:"connectionId"`
	ConnectionToken string `json:"connectionToken"`
}

func (c *Client) negotiate(ctx context.Context) (string, error) {
	u, err := url.Parse(c.baseURL + hubPath + "/negotiate")
	if err != nil {
		return "", fmt.Errorf("parse negotiate url: %w", err)
	}

	q := u.Query()
	q.Set("negotiateVersion", "1")
	q.Set("clientId", c.clientID.String())
	q.Set("cache", "true")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("create negotiate request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.getAccessToken())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("negotiate request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("read negotiate response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if isUnauthorized(resp.StatusCode) {
			return "", newAuthError(resp.StatusCode, string(body))
		}

		return "", fmt.Errorf("negotiate: unexpected status %d: %s", resp.StatusCode, truncate(string(body), 256))
	}

	var parsed negotiateResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", fmt.Errorf("decode negotiate response: %w", err)
	}

	token := parsed.ConnectionToken
	if token == "" {
		token = parsed.ConnectionID
	}
	if token == "" {
		return "", fmt.Errorf("negotiate: empty connection token")
	}

	return token, nil
}

func (c *Client) websocketURL(connectionToken string) (string, error) {
	u, err := url.Parse(c.baseURL + hubPath)
	if err != nil {
		return "", fmt.Errorf("parse websocket url: %w", err)
	}

	switch u.Scheme {
	case "https":
		u.Scheme = "wss"
	case "http":
		u.Scheme = "ws"
	default:
		return "", fmt.Errorf("unsupported uri scheme %q", u.Scheme)
	}

	q := u.Query()
	q.Set("clientId", c.clientID.String())
	q.Set("cache", "true")
	q.Set("id", connectionToken)
	q.Set("access_token", c.getAccessToken())
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (c *Client) handshake(ctx context.Context, conn *websocket.Conn) error {
	handshakeCtx, cancel := context.WithTimeout(ctx, handshakeTimeout)
	defer cancel()

	payload := []byte(`{"protocol":"json","version":1}` + string(recordSeparator))
	if err := conn.Write(handshakeCtx, websocket.MessageText, payload); err != nil {
		return fmt.Errorf("write handshake: %w", err)
	}

	_, data, err := conn.Read(handshakeCtx)
	if err != nil {
		return fmt.Errorf("read handshake: %w", err)
	}

	for _, frame := range splitRecords(data) {
		if len(bytes.TrimSpace(frame)) == 0 {
			continue
		}
		var resp struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(frame, &resp); err != nil {
			return fmt.Errorf("decode handshake response: %w", err)
		}
		if resp.Error != "" {
			return fmt.Errorf("signalr handshake error: %s", resp.Error)
		}
	}

	return nil
}

func (c *Client) readLoop(ctx context.Context, conn *websocket.Conn, out chan<- Message, errc chan<- error) {
	defer func() { _ = conn.Close(websocket.StatusNormalClosure, "") }()

	pingTicker := time.NewTicker(defaultPingPeriod)
	defer pingTicker.Stop()

	rawFrames := make(chan []byte, 8)
	readDone := make(chan error, 1)

	go func() {
		defer close(rawFrames)
		for {
			_, data, err := conn.Read(ctx)
			if err != nil {
				readDone <- err

				return
			}
			select {
			case rawFrames <- data:
			case <-ctx.Done():
				return
			}
		}
	}()

	sendErr := func(err error) {
		if err == nil || ctx.Err() != nil {
			return
		}
		select {
		case errc <- err:
		default:
		}
	}

	for {
		select {
		case <-ctx.Done():
			sendErr(ctx.Err())

			return
		case <-pingTicker.C:
			ping := []byte(`{"type":6}` + string(recordSeparator))
			writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			err := conn.Write(writeCtx, websocket.MessageText, ping)
			cancel()
			if err != nil {
				sendErr(fmt.Errorf("write ping: %w", err))

				return
			}
		case err := <-readDone:
			sendErr(err)

			return
		case data, ok := <-rawFrames:
			if !ok {
				return
			}
			for _, frame := range splitRecords(data) {
				msg, parsed, err := parseFrame(frame)
				if err != nil || !parsed {
					continue
				}
				if msg.Type == messageTypePing {
					continue
				}
				if msg.Type == messageTypeClose {
					sendErr(fmt.Errorf("signalr close: %s", msg.Error))

					return
				}
				if msg.NeedSyncClientState {
					if err := c.subscribeScanNotifications(ctx, conn); err != nil {
						sendErr(fmt.Errorf("re-subscribe after NeedSyncClientState: %w", err))

						return
					}

					continue
				}
				select {
				case out <- msg:
				case <-ctx.Done():
					sendErr(ctx.Err())

					return
				}
			}
		}
	}
}

func splitRecords(data []byte) [][]byte {
	parts := bytes.Split(data, []byte{recordSeparator})
	out := make([][]byte, 0, len(parts))
	for _, p := range parts {
		if len(p) == 0 {
			continue
		}
		out = append(out, p)
	}

	return out
}

func cloneHTTPClientWithTLSSkip(src *http.Client) *http.Client {
	clone := *src
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if src.Transport != nil {
		if t, ok := src.Transport.(*http.Transport); ok {
			transport = t.Clone()
		}
	}
	if transport.TLSClientConfig == nil {
		transport.TLSClientConfig = &tls.Config{}
	}
	transport.TLSClientConfig.InsecureSkipVerify = true
	clone.Transport = transport

	return &clone
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}

	return s[:n] + "..."
}
