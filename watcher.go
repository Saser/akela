package main

import (
	"log"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type Watcher struct {
	logger *log.Logger
	xu     *xgbutil.XUtil

	done chan struct{}
}

func NewWatcher(
	logger *log.Logger,
	xu *xgbutil.XUtil,
) *Watcher {
	return &Watcher{
		logger: logger,
		xu:     xu,

		done: make(chan struct{}, 1),
	}
}

func (w *Watcher) focusInFun(xu *xgbutil.XUtil, e xevent.FocusInEvent) {
	w.logger.Printf("FocusIn: %+v", e)
}

func (w *Watcher) watch(window xproto.Window) {
	w.logger.Printf("adding callback for FocusIn events on window %d", window)
	xwindow.New(w.xu, window).Listen(xproto.EventMaskFocusChange)
	xevent.FocusInFun(w.focusInFun).Connect(w.xu, window)
	w.logger.Printf("added callback for FocusIn events on window %d", window)
}

func (w *Watcher) unwatch(window xproto.Window) {
	w.logger.Printf("detaching all callbacks on window %d", window)
	xwindow.New(w.xu, window).Listen(xproto.EventMaskNoEvent)
	xevent.Detach(w.xu, window)
	w.logger.Printf("detached all callbacks on window %d", window)
}

func (w *Watcher) Start() error {
	// Set up callbacks for window creation and destruction events noticed by
	// the root window.
	root := xwindow.New(w.xu, w.xu.RootWin())
	w.logger.Printf("root window ID: %d", root.Id)
	root.Listen(xproto.EventMaskSubstructureNotify)
	w.logger.Println("adding callbacks for CreateNotify and DestroyNotify events on root window")
	xevent.CreateNotifyFun(func(xu *xgbutil.XUtil, e xevent.CreateNotifyEvent) {
		w.watch(e.Window)
	}).Connect(w.xu, root.Id)
	xevent.DestroyNotifyFun(func(xu *xgbutil.XUtil, e xevent.DestroyNotifyEvent) {
		w.unwatch(e.Window)
	}).Connect(w.xu, root.Id)
	defer func() {
		w.logger.Println("detaching all callbacks on root window")
		xevent.Detach(w.xu, root.Id)
		w.logger.Println("detached all callbacks on root window")
	}()
	w.logger.Println("added callbacks for CreateNotify and DestroyNotify events on root window")

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
