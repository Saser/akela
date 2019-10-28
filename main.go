package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/Saser/aukela/config"
	"golang.org/x/sync/errgroup"
)

var (
	configPath = flag.String("config", "", "path to configuration file")
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

	flag.Parse()
	if *configPath == "" {
		logger.Fatalf("no configuration file given, exiting")
	}

	logger.Printf("parsing configuration file at %s", *configPath)
	configFile, err := os.Open(*configPath)
	if err != nil {
		logger.Fatal(err)
	}
	config, err := config.Parse(configFile)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Printf("parsed configuration file at %s", *configPath)

	// Set up a connection to the X server.
	logger.Println("setting up connection to X server")
	xu, err := xgbutil.NewConn()
	if err != nil {
		logger.Fatalf("setting up connection to X server failed: %+v", err)
	}
	logger.Println("set up connection to X server")

	// Create a `Watcher`.
	logger.Println("creating watcher")
	watcher := NewWatcher(logger, xu)
	logger.Println("created watcher")

	// Create a `Switcher`.
	logger.Println("creating switcher")
	switcher := NewSwitcher(logger, xu, config, watcher.Focused())
	logger.Println("created switcher")

	// Set up an errgroup with a derived context that is cancelled when an
	// interrupt is received or when one errgroup task returns a non-nil error.
	ctx, cancel := withCancelOnInterrupt(context.Background())
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)

	// Run the watcher in the errgroup.
	g.Go(func() error {
		logger.Println("starting watcher")
		if err := watcher.Start(); err != nil {
			return err
		}
		logger.Println("watcher finished")
		return nil
	})
	cleanupWatcher := func() {
		logger.Println("stopping watcher")
		watcher.Stop()
		logger.Println("stopped watcher")
	}

	// Run the switcher in the errgroup.
	g.Go(func() error {
		logger.Println("starting switcher")
		if err := switcher.Start(); err != nil {
			return err
		}
		logger.Println("switcher finished")
		return nil
	})
	cleanupSwitcher := func() {
		logger.Println("stopping switcher")
		switcher.Stop()
		logger.Println("stopped switcher")
	}

	// Run the X event loop in the errgroup.
	g.Go(func() error {
		logger.Println("starting event loop")
		xevent.Main(xu)
		logger.Println("event loop finished")
		return nil
	})
	cleanupEventLoop := func() {
		logger.Println("stopping event loop")
		xevent.Quit(xu)
		logger.Println("stopped event loop")
	}

	logger.Println("waiting for context cancellation")
	<-ctx.Done()
	logger.Println("context cancelled")
	cleanupEventLoop()
	cleanupSwitcher()
	cleanupWatcher()
	if err := g.Wait(); err != nil {
		logger.Fatalf("error in errgroup: %+v", err)
	}
	logger.Println("goodbye")
}
