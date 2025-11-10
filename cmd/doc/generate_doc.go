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

	app := application.NewApplication()

	if err := app.GenerateDoc("./doc"); err != nil {
		log.Fatal(err)
	}
}
