package switcher

import (
	"log"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/Saser/aukela/config"
)

type Switcher struct {
	logger *log.Logger
	xu     *xgbutil.XUtil
	config *config.Config

	focused <-chan xproto.Window

	done chan struct{}
}

func New(
	logger *log.Logger,
	xu *xgbutil.XUtil,
	config *config.Config,
	focused <-chan xproto.Window,
) *Switcher {
	return &Switcher{
		logger: logger,
		xu:     xu,
		config: config,

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
			s.logger.Printf("getting class of window %d", window)
			wmClass, err := icccm.WmClassGet(s.xu, window)
			if err != nil {
				s.logger.Printf("getting class of window %d failed: %+v", window, err)
				continue
			}
			s.logger.Printf("got class of window %d: %+v", window, wmClass)
			class := wmClass.Class
			if spec, ok := s.config.Classes[class]; ok {
				s.logger.Printf(`found spec for class "%s": %+v`, class, spec)
			} else {
				s.logger.Printf(`no spec for class "%s", using default spec`, class)
			}
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
