package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

const (
	VersionPath = "/api/versions/package/current"
	LicensePath = "/api/license"
)

// LicenseStep is one canned response for GET /api/license.
type LicenseStep struct {
	Status  int
	IsValid *bool
}

// InitServer is an httptest stub for AI client initialization: signin, version, license.
type InitServer struct {
	*httptest.Server

	mu sync.Mutex

	apiToken          string
	version           string
	signinSteps       []Step
	signinFallback    Step
	hasSigninFallback bool
	licenseSteps      []LicenseStep

	SigninCalls  int
	LicenseCalls int
	LastAPIToken string
}

// InitServerConfig configures a new InitServer.
type InitServerConfig struct {
	APIToken string
	Version  string
	Signin   []Step
	License  []LicenseStep
}

// NewInitServer starts a stub server for adapter initialization tests.
func NewInitServer(t *testing.T, cfg InitServerConfig) *InitServer {
	t.Helper()

	s := &InitServer{
		apiToken:     cfg.APIToken,
		version:      cfg.Version,
		signinSteps:  append([]Step(nil), cfg.Signin...),
		licenseSteps: append([]LicenseStep(nil), cfg.License...),
	}
	s.Server = httptest.NewServer(http.HandlerFunc(s.serve))
	t.Cleanup(s.Close)

	return s
}

func (s *InitServer) serve(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch r.URL.Path {
	case SigninPath:
		s.SigninCalls++
		s.LastAPIToken = r.Header.Get("Access-Token")
		writeAuthStep(w, s.nextSigninLocked())
	case VersionPath:
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(s.version))
	case LicensePath:
		s.LicenseCalls++
		writeLicenseStep(w, s.nextLicenseLocked())
	default:
		http.NotFound(w, r)
	}
}

func (s *InitServer) nextSigninLocked() Step {
	if len(s.signinSteps) == 0 {
		if s.hasSigninFallback {
			return s.signinFallback
		}

		return Step{Status: http.StatusInternalServerError}
	}

	step := s.signinSteps[0]
	if step.Status == 0 || step.Status == http.StatusOK {
		s.signinFallback = step
		s.hasSigninFallback = true
	}
	if len(s.signinSteps) > 1 {
		s.signinSteps = s.signinSteps[1:]
	}

	return step
}

func (s *InitServer) nextLicenseLocked() LicenseStep {
	if len(s.licenseSteps) == 0 {
		return LicenseStep{Status: http.StatusInternalServerError}
	}

	step := s.licenseSteps[0]
	s.licenseSteps = s.licenseSteps[1:]

	return step
}

func writeAuthStep(w http.ResponseWriter, step Step) {
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

func writeLicenseStep(w http.ResponseWriter, step LicenseStep) {
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

	body := map[string]bool{}
	if step.IsValid != nil {
		body["isValid"] = *step.IsValid
	}
	_ = json.NewEncoder(w).Encode(body)
}

// LicenseOK is a successful license check response.
func LicenseOK() LicenseStep {
	valid := true

	return LicenseStep{Status: http.StatusOK, IsValid: &valid}
}

// LicenseInvalid is a 200 response with isValid=false.
func LicenseInvalid() LicenseStep {
	valid := false

	return LicenseStep{Status: http.StatusOK, IsValid: &valid}
}
