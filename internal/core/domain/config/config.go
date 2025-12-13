package config

import (
	"fmt"

	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
)

type Config struct {
	uri       Uri
	token     string
	tlsSkip   bool
	projectId uuid.UUID
	branchId  uuid.UUID
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
		return errs.NewValidationRequiredError("token")
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

func (cfg *Config) SetTLSSkip(tlsSkip bool) {
	cfg.tlsSkip = tlsSkip
}

func (cfg *Config) ProjectId() uuid.UUID {
	return cfg.projectId
}

func (cfg *Config) SetProjectId(projectIdFlag string) error {
	if projectIdFlag == "" {
		return errs.NewValidationRequiredError("project-id")
	}

	projectId, err := uuid.Parse(projectIdFlag)
	if err != nil {
		return errs.NewValidationFieldError("project-id", fmt.Sprintf("'%s' invalid uuid", projectIdFlag))
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
		return errs.NewValidationRequiredError("branch-id")
	}

	branchId, err := uuid.Parse(branchIdFlag)
	if err != nil {
		return errs.NewValidationFieldError("branch-id", fmt.Sprintf("'%s' invalud uuid", branchId))
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
		return errs.NewValidationRequiredError("uri")
	}

	if cfg.token == "" {
		return errs.NewValidationRequiredError("token")
	}

	return nil
}

func (cfg *Config) ValidateProjectId() error {
	if cfg.ProjectId() == uuid.Nil {
		return errs.NewValidationRequiredError("projectId")
	}

	return nil
}

func (cfg *Config) ValidateBranchId() error {
	if cfg.BranchId() == uuid.Nil {
		return errs.NewValidationRequiredError("branchId")
	}

	return nil
}
