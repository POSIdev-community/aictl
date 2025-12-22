package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/POSIdev-community/aictl/pkg/fshelper"
	. "github.com/dave/jennifer/jen"
)

const innerPresenterPath = "internal/presenter"
const innerUseCasePath = "internal/core/application/usecase"

// Generate presenter and usecase boilerplate code
func main() {
	newCmdPath := []string{"set", "project", "settings"}

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if err := generatePresenterCode(currentDir, newCmdPath); err != nil {
		panic(err)
	}

	if err := generateUseCaseCode(currentDir, newCmdPath); err != nil {
		panic(err)
	}
}

func generatePresenterCode(currentDir string, newCmdPath []string) error {
	presenterPath := path.Join(currentDir, innerPresenterPath)

	cmdPath := path.Join(presenterPath, newCmdPath[0])

	if !fshelper.PathExists(cmdPath) {
		err := os.Mkdir(cmdPath, 0755)
		if err != nil {
			return fmt.Errorf("make dir: %v", err)
		}
	}

	filename := strings.Join(newCmdPath, "_") + ".go"
	newCmdFilePath := path.Join(cmdPath, filename)

	if err := getPresenterCode(newCmdPath, newCmdFilePath); err != nil {
		return err
	}

	return nil
}

func getPresenterCode(cmdPath []string, filepath string) error {
	use := cmdPath[len(cmdPath)-1]
	useCase := strings.Join(cmdPath, " ")
	short := capitalize(useCase)

	cmdName := createCmdName(cmdPath)

	f := NewFile(cmdPath[0])

	// imports
	cobra := "github.com/spf13/cobra"
	application := "github.com/POSIdev-community/aictl/internal/core/application"
	config := "github.com/POSIdev-community/aictl/internal/core/domain/config"
	utils := "github.com/POSIdev-community/aictl/internal/presenter/.utils"

	f.ImportName(cobra, "cobra")
	f.ImportName(application, "application")
	f.ImportName(config, "config")
	f.ImportName(utils, ".utils")

	// body
	f.Func().Id("New"+cmdName+"Cmd").Params(
		Id("cfg").Op("*").Qual(config, "Config"),
		Id("depsContainer").Op("*").Qual(application, "DependenciesContainer"),
	).Op("*").Qual(cobra, "Command").Block(
		Id("cmd").Op(":=").Op("&").Qual(cobra, "Command").Values(Dict{
			Id("Use"):   Lit(use),
			Id("Short"): Lit(short),
			Id("PreRunE"): Func().Params(
				Id("cmd").Op("*").Qual(cobra, "Command"),
				Id("args").Op("[]").String(),
			).Error().Block(
				Return().Nil()),
			Id("RunE"): Func().Params(
				Id("cmd").Op("*").Qual(cobra, "Command"),
				Id("args").Op("[]").String(),
			).Error().Block(
				Id("ctx").Op(":=").Id("cmd").Dot("Context").Call(),
				Line(),
				Id("useCase").Op(",").Id("err").Op(":=").Id("depsContainer").Dot(cmdName+"UseCase").Call(Id("ctx"), Id("cfg")),
				If(Err().Op("!=").Nil()).Block(
					Return(Qual("fmt", "Errorf").Call(Lit("presenter "+useCase+" useCase error: %w"), Id("err"))),
				),
				Line(),
				If(Err().Op(":=").Id("useCase").Dot("Execute").Call(Id("ctx")), Err().Op("!=").Nil()).Block(
					Id("cmd").Dot("SilenceUsage").Op("=").True(),
					Line(),
					Return(Qual("fmt", "Errorf").Call(Lit("presenter "+useCase+": %w"), Id("err"))),
				),
				Line(),
				Return(Nil()),
			),
		}),
		Line(),
		Return(Id("cmd")),
	)

	if err := f.Save(filepath); err != nil {
		return fmt.Errorf("save file: %v", err)
	}

	return nil
}

func createCmdName(cmdPath []string) string {
	var builder strings.Builder

	for _, p := range cmdPath {
		builder.WriteString(capitalize(p))
	}

	return builder.String()
}

func generateUseCaseCode(currentDir string, newCmdPath []string) error {
	useCasePath := path.Join(currentDir, innerUseCasePath)

	allParts := append([]string{useCasePath}, newCmdPath...)
	cmdPath := filepath.Join(allParts...)
	err := os.MkdirAll(cmdPath, 0755)
	if err != nil {
		return fmt.Errorf("make dir: %v", err)
	}
	filename := newCmdPath[len(newCmdPath)-1] + ".go"
	newCmdFilePath := path.Join(cmdPath, filename)

	if err := getUseCaseCode(newCmdPath, newCmdFilePath); err != nil {
		return err
	}

	return nil
}

func getUseCaseCode(cmdPath []string, filepath string) error {
	packageName := cmdPath[len(cmdPath)-1]

	f := NewFile(packageName)

	// Добавляем импорты
	f.ImportName("context", "context")
	f.ImportName("fmt", "fmt")
	f.ImportName("github.com/POSIdev-community/aictl/internal/core/port", "port")
	f.ImportName("github.com/POSIdev-community/aictl/pkg/errs", "errs")

	// Генерируем структуру UseCase
	f.Type().Id("UseCase").Struct(
		Id("aiAdapter").Qual("github.com/POSIdev-community/aictl/internal/core/port", "Ai"),
		Id("cliAdapter").Qual("github.com/POSIdev-community/aictl/internal/core/port", "Cli"),
	)

	f.Line()

	// Генерируем функцию NewUseCase
	f.Func().Id("NewUseCase").Params(
		Id("aiAdapter").Qual("github.com/POSIdev-community/aictl/internal/core/port", "Ai"),
		Id("cliAdapter").Qual("github.com/POSIdev-community/aictl/internal/core/port", "Cli"),
	).Params(Op("*").Id("UseCase"), Error()).Block(
		If(Id("aiAdapter").Op("==").Nil()).Block(
			Return(Nil(), Qual("github.com/POSIdev-community/aictl/pkg/errs", "NewValidationRequiredError").Call(Lit("aiAdapter"))),
		),
		Line(),
		If(Id("cliAdapter").Op("==").Nil()).Block(
			Return(Nil(), Qual("github.com/POSIdev-community/aictl/pkg/errs", "NewValidationRequiredError").Call(Lit("cliAdapter"))),
		),
		Line(),
		Return(
			Op("&").Id("UseCase").Values(Dict{
				Id("aiAdapter"):  Id("aiAdapter"),
				Id("cliAdapter"): Id("cliAdapter"),
			}),
			Nil(),
		),
	)

	f.Line()

	// Генерируем метод Execute
	f.Func().Params(Id("u").Op("*").Id("UseCase")).Id("Execute").Params(
		Id("ctx").Qual("context", "Context"),
	).Error().Block(
		Return(Nil()),
	)

	if err := f.Save(filepath); err != nil {
		return fmt.Errorf("save file: %v", err)
	}

	return nil
}
