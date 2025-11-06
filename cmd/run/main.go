package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/POSIdev-community/aictl/internal/application"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	app := application.NewApplication()

	app.Run(ctx)
}
