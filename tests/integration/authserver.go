// Package integration holds in-process HTTP stubs and integration-style unit tests.
package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

const (
	SigninPath  = "/api/auth/signin"
	RefreshPath = "/api/auth/refreshToken"
)

// Step is one canned response for signin or refresh.
type Step struct {
	Status      int
	AccessToken string
	RefreshTok  string
}

// AuthServer is an httptest stub that serves queued responses for auth endpoints.
type AuthServer struct {
	*httptest.Server

	mu sync.Mutex

	signinSteps  []Step
	refreshSteps []Step

	APIToken string

	SigninCalls  int
	RefreshCalls int
	LastAPIToken string
}

// NewAuthServer starts a stub server. Queued steps are consumed in order per endpoint.
func NewAuthServer(t *testing.T, apiToken string, signin, refresh []Step) *AuthServer {
	t.Helper()

	s := &AuthServer{
		APIToken:     apiToken,
		signinSteps:  append([]Step(nil), signin...),
		refreshSteps: append([]Step(nil), refresh...),
	}
	s.Server = httptest.NewServer(http.HandlerFunc(s.serve))
	t.Cleanup(s.Close)

	return s
}

func (s *AuthServer) serve(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch r.URL.Path {
	case SigninPath:
		s.SigninCalls++
		s.LastAPIToken = r.Header.Get("Access-Token")
		s.writeStep(w, s.nextSigninLocked())
	case RefreshPath:
		s.RefreshCalls++
		s.writeStep(w, s.nextRefreshLocked())
	default:
		http.NotFound(w, r)
	}
}

func (s *AuthServer) nextSigninLocked() Step {
	if len(s.signinSteps) == 0 {
		return Step{Status: http.StatusInternalServerError}
	}
	step := s.signinSteps[0]
	s.signinSteps = s.signinSteps[1:]

	return step
}

func (s *AuthServer) nextRefreshLocked() Step {
	if len(s.refreshSteps) == 0 {
		return Step{Status: http.StatusInternalServerError}
	}
	step := s.refreshSteps[0]
	s.refreshSteps = s.refreshSteps[1:]

	return step
}

func (s *AuthServer) writeStep(w http.ResponseWriter, step Step) {
	status := step.Status
	if status == 0 {
		status = http.StatusOK
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if status >= 400 {
		_, _ = w.Write([]byte(`{}`))

		return
	}

	body := map[string]string{}
	if step.AccessToken != "" {
		body["accessToken"] = step.AccessToken
	}
	if step.RefreshTok != "" {
		body["refreshToken"] = step.RefreshTok
	}
	_ = json.NewEncoder(w).Encode(body)
}

// OKSignin is a successful sign-in response.
func OKSignin(access, refresh string) Step {
	return Step{Status: http.StatusOK, AccessToken: access, RefreshTok: refresh}
}

// OKRefresh is a successful refresh response.
func OKRefresh(access string) Step {
	return Step{Status: http.StatusOK, AccessToken: access}
}

// Unauthorized is a 401 response.
func Unauthorized() Step {
	return Step{Status: http.StatusUnauthorized}
}

// ServerError is a 5xx response.
func ServerError(status int) Step {
	if status < 500 {
		status = http.StatusServiceUnavailable
	}

	return Step{Status: status}
}
