package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/POSIdev-community/aictl/internal/application"
)

func main() {
	_, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	app, err := application.NewApplication()
	if err != nil {
		log.Fatal(err)
	}

	outDir := "./doc/gen"
	if len(os.Args) > 1 {
		outDir = os.Args[1]
	}

	if err := app.GenerateDoc(outDir); err != nil {
		log.Fatal(err)
	}
}
