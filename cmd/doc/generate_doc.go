package main

import (
	"context"
	"github.com/POSIdev-community/aictl/internal/application"
	"log"
	"os"
	"os/signal"
)

func main() {
	_, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	app := application.NewApplication()

	if err := app.GenerateDoc("./doc"); err != nil {
		log.Fatal(err)
	}
}
