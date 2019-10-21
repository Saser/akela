package main

import "log"

type Switcher struct {
	logger *log.Logger

	done chan struct{}
}

func NewSwitcher(logger *log.Logger) *Switcher {
	return &Switcher{
		logger: logger,

		done: make(chan struct{}, 1),
	}
}

func (s *Switcher) Start() error {
	s.logger.Println("starting switcher")
	<-s.done
	s.logger.Println("received done signal")
	s.logger.Println("stopped switcher")
	return nil
}

func (s *Switcher) Stop() {
	s.logger.Println("sending done signal")
	s.done <- struct{}{}
	s.logger.Println("sent done signal")
}
