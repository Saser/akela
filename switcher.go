package main

import (
	"log"

	"github.com/BurntSushi/xgb/xproto"
)

type Switcher struct {
	logger *log.Logger

	focused <-chan xproto.Window

	done chan struct{}
}

func NewSwitcher(
	logger *log.Logger,
	focused <-chan xproto.Window,
) *Switcher {
	return &Switcher{
		logger: logger,

		focused: focused,

		done: make(chan struct{}, 1),
	}
}

func (s *Switcher) Start() error {
	s.logger.Println("starting switcher")
loop:
	for {
		select {
		case window := <-s.focused:
			s.logger.Printf("received focus on window %d", window)
		case <-s.done:
			s.logger.Println("received done signal")
			break loop
		}
	}
	s.logger.Println("stopped switcher")
	return nil
}

func (s *Switcher) Stop() {
	s.logger.Println("sending done signal")
	s.done <- struct{}{}
	s.logger.Println("sent done signal")
}
