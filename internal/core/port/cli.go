package port

import (
	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/scan"
)

type Cli interface {
	AskConfirmation(question string) (bool, error)

	ShowProjects(projects []project.Project)
	ShowProjectsQuite(projects []project.Project)

	ShowText(text string)

	ShowScans(scans []scan.Scan)
}
