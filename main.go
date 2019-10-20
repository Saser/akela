package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/BurntSushi/xgbutil"
)

func withCancelOnInterrupt(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cancel()
	}()
	return ctx, cancel
}

func main() {
	logger := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)

	// Set up a connection to the X server.
	logger.Println("setting up connection to X server")
	_, err := xgbutil.NewConn()
	if err != nil {
		logger.Fatalf("setting up connection to X server failed: %+v", err)
	}
	logger.Println("set up connection to X server")

	ctx, cancel := withCancelOnInterrupt(context.Background())
	defer cancel()
	logger.Println("waiting for interrupt")
	<-ctx.Done()
	logger.Println("interrupted, goodbye")
}
