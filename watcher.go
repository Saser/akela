package main

import (
	"log"
)

type Watcher struct {
	logger *log.Logger

	done chan struct{}
}

func NewWatcher(
	logger *log.Logger,
) *Watcher {
	return &Watcher{
		logger: logger,

		done: make(chan struct{}, 1),
	}
}

func (w *Watcher) Start() error {
	w.logger.Println("starting watcher")
	<-w.done
	w.logger.Println("received done signal")
	w.logger.Println("stopped watcher")
	return nil
}

func (w *Watcher) Stop() {
	w.logger.Println("sending done signal")
	w.done <- struct{}{}
	w.logger.Println("sent done signal")
}
