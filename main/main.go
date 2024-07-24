package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/makasim/gogame/internal/app"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	go handleSignals(cancel)

	cfg := app.Config{}

	if err := app.New(cfg).Run(ctx); err != nil {
		log.Printf("ERROR: %v", err)
		os.Exit(1)
	}
}

func handleSignals(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	<-signals
	log.Printf("INFO: got signal; canceling context")
	cancel()

	<-signals
	log.Printf("WARN: got second signal; force exiting")
	os.Exit(1)
}
