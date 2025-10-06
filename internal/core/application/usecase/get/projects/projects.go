package projects

import (
	"context"
	"fmt"
	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/regexfilter"
	"github.com/POSIdev-community/aictl/internal/core/port"
	"github.com/POSIdev-community/aictl/pkg/errs"
)

type UseCase struct {
	aiAdapter  port.Ai
	cliAdapter port.Cli
}

func NewUseCase(aiAdapter port.Ai, cliAdapter port.Cli) (*UseCase, error) {
	if aiAdapter == nil {
		return nil, errs.NewValidationRequiredError("aiAdapter")
	}

	if cliAdapter == nil {
		return nil, errs.NewValidationRequiredError("cliAdapter")
	}

	return &UseCase{aiAdapter, cliAdapter}, nil
}

func (u *UseCase) Execute(ctx context.Context, filter regexfilter.RegexFilter, quite bool) error {
	projects, err := u.aiAdapter.GetProjects(ctx)
	if err != nil {
		return fmt.Errorf("get projects: %w", err)
	}

	filteredProjects := make([]project.Project, 0, len(projects))
	if filter.Empty() {
		filteredProjects = projects
	} else {
		for _, p := range projects {
			matched := filter.Execute(p.Name)

			if matched {
				filteredProjects = append(filteredProjects, p)
			}
		}
	}

	if quite {
		u.cliAdapter.ShowProjectsQuite(filteredProjects)
	} else {
		u.cliAdapter.ShowProjects(filteredProjects)
	}

	return nil
}
