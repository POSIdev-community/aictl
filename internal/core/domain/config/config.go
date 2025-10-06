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

func (c *Config) Token() string {
	return c.token
}

func (c *Config) SetToken(token string) error {
	if token == "" {
		return errs.NewValidationRequiredError("token")
	}

	c.token = token

	return nil
}

func (c *Config) Uri() Uri {
	return c.uri
}

func (c *Config) UriString() string {
	return c.uri.value
}

func (c *Config) SetURI(rawUri string) error {

	uri, err := NewUri(rawUri)
	if err != nil {
		c.uri = Uri{}

		return fmt.Errorf("set Uri error: %w", err)
	}

	c.uri = uri

	return nil
}

func (c *Config) TLSSkip() bool {
	return c.tlsSkip
}

func (c *Config) SetTLSSkip(tlsSkip bool) {
	c.tlsSkip = tlsSkip
}

func (c *Config) ProjectId() uuid.UUID {
	return c.projectId
}

func (c *Config) SetProjectId(projectIdFlag string) error {
	if projectIdFlag == "" {
		return errs.NewValidationRequiredError("project-id")
	}

	projectId, err := uuid.Parse(projectIdFlag)
	if err != nil {
		return errs.NewValidationFieldError("project-id", fmt.Sprintf("'%s' invalud uuid", projectIdFlag))
	}

	c.projectId = projectId

	return nil
}

func (c *Config) BranchId() uuid.UUID {
	return c.branchId
}

func (c *Config) SetBranchId(branchIdFlag string) error {
	if branchIdFlag == "" {
		return errs.NewValidationRequiredError("branch-id")
	}

	branchId, err := uuid.Parse(branchIdFlag)
	if err != nil {
		return errs.NewValidationFieldError("branch-id", fmt.Sprintf("'%s' invalud uuid", branchIdFlag))
	}

	c.branchId = branchId

	return nil
}

func (c *Config) Validate(projectIdRequired, branchIdRequired bool) error {
	if err := c.uri.validate(); err != nil {
		return errs.NewValidationRequiredError("uri")
	}

	if c.token == "" {
		return errs.NewValidationRequiredError("token")
	}

	if projectIdRequired {
		if c.ProjectId() == uuid.Nil {
			return errs.NewValidationRequiredError("projectId")
		}
	}

	if branchIdRequired {
		if c.BranchId() == uuid.Nil {
			return errs.NewValidationRequiredError("branchId")
		}
	}

	return nil
}
