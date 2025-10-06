package main

import (
	"context"
	"github.com/POSIdev-community/aictl/internal/application"
	"os"
	"os/signal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	app := application.NewApplication()

	app.Run(ctx)
}
