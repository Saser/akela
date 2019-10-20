package main

import (
	"context"
	"log"
	"os"
	"os/signal"
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
	ctx, cancel := withCancelOnInterrupt(context.Background())
	defer cancel()
	logger.Println("waiting for interrupt")
	<-ctx.Done()
	logger.Println("interrupted, goodbye")
}
