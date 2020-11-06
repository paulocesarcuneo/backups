package quit

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Control struct {
	subs   []chan interface{}
	lock   sync.Mutex
	active bool
}

func NewControl() Control {
	return Control{
		active: true,
		lock:   sync.Mutex{},
		subs:   nil,
	}
}

func (control *Control) Sub() (chan interface{}, error) {
	control.lock.Lock()
	defer control.lock.Unlock()

	var ch chan interface{}
	if !control.active {
		return nil, errors.New("Control: Already quitted")
	} else {
		ch = make(chan interface{}, 1)
		control.subs = append(control.subs, ch)
		return ch, nil
	}
}

func (control *Control) UnSub(quitter chan interface{}) {
	control.lock.Lock()
	defer control.lock.Unlock()

	var updated []chan interface{} = nil
	for _, ch := range control.subs {
		if ch == quitter {
			continue
		}
		updated = append(updated, ch)

	}
	control.subs = updated
}

func (control *Control) finish() {
	control.lock.Lock()
	defer control.lock.Unlock()

	control.active = false
	for _, ch := range control.subs {
		ch <- true
		<-ch
	}
}

func (control *Control) WaitForTermSignal() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT)
	<-sigc
	log.Println("Signal Received")
	control.finish()
}
