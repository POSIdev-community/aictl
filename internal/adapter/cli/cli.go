package cli

import (
	"fmt"
	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/scan"
	"github.com/POSIdev-community/aictl/internal/core/port"
	"strings"
)

var _ port.Cli = &Cli{}

type Cli struct{}

func NewCli() *Cli {
	return &Cli{}
}

func (cli *Cli) AskConfirmation(question string) (bool, error) {
	fmt.Printf("%s [y/n]: ", question)

	var answer string
	_, err := fmt.Scan(&answer)
	if err != nil {
		return false, err
	}

	return strings.ToLower(answer) == "y" ||
		strings.ToLower(answer) == "yes", nil
}

func (cli *Cli) ShowProjects(projects []project.Project) {
	const format = "%-36s\t%s\n"

	fmt.Printf(format, "ID", "NAME")

	for _, p := range projects {
		fmt.Printf(format, p.Id, p.Name)
	}
}

func (cli *Cli) ShowProjectsQuite(projects []project.Project) {
	for _, p := range projects {
		fmt.Println(p.Id)
	}
}

func (cli *Cli) ShowText(text string) {
	fmt.Println(text)
}

func (cli *Cli) ShowScans(scans []scan.Scan) {
	const format = "%-36s\t%-36s\n"

	fmt.Printf(format, "ID", "SETTINGS ID")

	for _, p := range scans {
		fmt.Printf(format, p.Id, p.SettingsId)
	}
}
