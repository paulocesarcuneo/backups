package quit

import (
	"sync"
	"log"
)

type Control struct {
	subs []chan interface{}
	lock sync.Mutex
	active bool
}

var control = Control{subs: nil, active:true, lock: sync.Mutex{}}

func Sub() chan interface{} {
	var ch chan interface{}
	control.lock.Lock()
	if !control.active {
		ch = nil
	} else {
		ch = make(chan interface{})
		control.subs = append(control.subs, ch)
	}
	control.lock.Unlock()
	if ch == nil {
		panic("Already Quited")
	}
	// log.Println("subs", len(control.subs))
	return ch
}

func filter(ss []string, test func(string) bool) (ret []string) {
    for _, s := range ss {
        if test(s) {
            ret = append(ret, s)
        }
    }
    return
}

func UnSub(quitter chan interface{}) {
	control.lock.Lock()
	var updated []chan interface{} = nil
	for _, ch := range control.subs {
		if ch == quitter {
			continue
		}
		updated = append(updated, ch)

	}
	control.subs = updated
	control.lock.Unlock()
}

func Quit(quitter chan interface{}) {
	control.lock.Lock()
	control.active = false
	for i, ch := range control.subs {
		log.Println("send quit", i)
		if ch == quitter {
			continue
		}
		ch<- true
	}
	control.lock.Unlock()
}
