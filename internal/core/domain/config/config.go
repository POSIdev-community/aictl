package config

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

type Config struct {
	uri        Uri
	token      string
	tlsSkip    bool
	caCertPath string
	projectId  uuid.UUID
	branchId   uuid.UUID
}

func NewConfig(uri Uri, token string, tlsSkip bool, projectId, branchId uuid.UUID) *Config {
	return &Config{
		uri:       uri,
		token:     token,
		tlsSkip:   tlsSkip,
		projectId: projectId,
		branchId:  branchId,
	}
}

func (cfg *Config) Token() string {
	return cfg.token
}

func (cfg *Config) SetToken(token string) error {
	if token == "" {
		return validation.NewRequiredError("token")
	}

	cfg.token = token

	return nil
}

func (cfg *Config) Uri() Uri {
	return cfg.uri
}

func (cfg *Config) UriString() string {
	return cfg.uri.value
}

func (cfg *Config) SetURI(rawUri string) error {
	uri, err := NewUri(rawUri)
	if err != nil {
		cfg.uri = Uri{}

		return fmt.Errorf("set Uri error: %w", err)
	}

	cfg.uri = uri

	return nil
}

func (cfg *Config) TLSSkip() bool {
	return cfg.tlsSkip
}

// SetTLSSkip enables or disables TLS verification skip.
// Enabling tls-skip fails if cacert is already set; clear it first via ClearCACertPath / ctx unset.
func (cfg *Config) SetTLSSkip(tlsSkip bool) error {
	if tlsSkip && cfg.caCertPath != "" {
		return validation.NewError("cannot enable 'tls-skip' while 'cacert' is set; unset cacert first")
	}

	cfg.tlsSkip = tlsSkip

	return nil
}

func (cfg *Config) CACertPath() string {
	return cfg.caCertPath
}

// SetCACertPath stores a non-empty path to a PEM CA file.
// Empty path is rejected; use ClearCACertPath to unset.
// Fails if tls-skip is already enabled; clear it first via SetTLSSkip(false) / ctx unset --tls-skip.
func (cfg *Config) SetCACertPath(path string) error {
	if path == "" {
		return validation.NewRequiredError("cacert")
	}
	if cfg.tlsSkip {
		return validation.NewError("cannot set 'cacert' while 'tls-skip' is enabled; unset tls-skip first")
	}

	cfg.caCertPath = path

	return nil
}

// ClearCACertPath removes the configured CA certificate path.
func (cfg *Config) ClearCACertPath() {
	cfg.caCertPath = ""
}

// ApplyCACertPath sets or clears the path without requiring non-empty
// (used when loading context / preserving value on unset).
func (cfg *Config) ApplyCACertPath(path string) error {
	if path != "" && cfg.tlsSkip {
		return validation.NewError("cannot use 'cacert' together with 'tls-skip'")
	}

	cfg.caCertPath = path

	return nil
}

func (cfg *Config) ProjectId() uuid.UUID {
	return cfg.projectId
}

func (cfg *Config) SetProjectId(projectIdFlag string) error {
	if projectIdFlag == "" {
		return validation.NewRequiredError("project-id")
	}

	projectId, err := uuid.Parse(projectIdFlag)
	if err != nil {
		return validation.NewFieldError("project-id", fmt.Sprintf("'%s' invalud uuid", projectIdFlag))
	}

	cfg.projectId = projectId

	return nil
}

func (cfg *Config) UpdateProjectId(projectIdFlag string) error {
	var err error
	if projectIdFlag != "" {
		err = cfg.SetProjectId(projectIdFlag)
		if err != nil {
			return err
		}
	} else {
		err = cfg.ValidateProjectId()
		if err != nil {
			return err
		}
	}

	return nil
}

func (cfg *Config) BranchId() uuid.UUID {
	return cfg.branchId
}

func (cfg *Config) SetBranchId(branchIdFlag string) error {
	if branchIdFlag == "" {
		return validation.NewRequiredError("branch-id")
	}

	branchId, err := uuid.Parse(branchIdFlag)
	if err != nil {
		return validation.NewFieldError("branch-id", fmt.Sprintf("'%s' invalud uuid", branchId))
	}

	cfg.branchId = branchId

	return nil
}

func (cfg *Config) UpdateBranchId(branchIdFlag string) error {
	var err error
	if branchIdFlag != "" {
		err = cfg.SetBranchId(branchIdFlag)
		if err != nil {
			return err
		}
	} else {
		err = cfg.ValidateBranchId()
		if err != nil {
			return err
		}
	}

	return nil
}

func (cfg *Config) Validate() error {
	if err := cfg.uri.validate(); err != nil {
		return validation.NewRequiredError("uri")
	}

	if cfg.token == "" {
		return validation.NewRequiredError("token")
	}

	if cfg.tlsSkip && cfg.caCertPath != "" {
		return validation.NewError("cannot use 'cacert' together with 'tls-skip'")
	}

	return nil
}

func (cfg *Config) ValidateProjectId() error {
	if cfg.ProjectId() == uuid.Nil {
		return validation.NewRequiredError("projectId")
	}

	return nil
}

func (cfg *Config) ValidateBranchId() error {
	if cfg.BranchId() == uuid.Nil {
		return validation.NewRequiredError("branchId")
	}

	return nil
}
