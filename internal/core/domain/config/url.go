package config

import (
	"net/url"
	"strings"

	"github.com/POSIdev-community/aictl/pkg/errs"
)

type Uri struct {
	value string

	createByConstructor bool
}

func NewUri(value string) (Uri, error) {
	if value == "" {
		return Uri{}, errs.NewValidationRequiredError("uri")
	}

	value = strings.TrimRight(value, "/")

	if _, err := url.ParseRequestURI(value); err != nil {
		return Uri{}, errs.NewValidationInvalidError("uri")
	}

	return Uri{value: value, createByConstructor: true}, nil
}

func (u Uri) validate() error {
	if u.createByConstructor {
		return nil
	}

	return errs.NewValidationInvalidError("uri")
}
