package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"winqroo/config"
)

func main() {
	config.Init()
	serverConfig := config.GetInstance()

	startApp(serverConfig) // Run as web server
}

func startApp(serverConfig *config.Config) {
	app := NewApp(serverConfig)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failed to start app: ", err)
	}
}
