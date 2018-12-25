package main

import (
	"sync"
	"time"
)

type watchStatus struct {
	received bool
	mux      sync.Mutex
}

func (w *watchStatus) setReceived(b bool) {
	w.mux.Lock()
	defer w.mux.Unlock()
	w.received = b
}

func (w *watchStatus) isReceived() bool {
	w.mux.Lock()
	defer w.mux.Unlock()
	return w.received
}

func newDebouncer(sender chan string, d int) func() {
	duration := time.Duration(d) * time.Second
	status := watchStatus{}
	timer := time.NewTimer(duration)

	go func() {
		for {
			<-sender
			timer.Reset(duration)
			status.setReceived(true)
		}
	}()

	return func() {
		for {
			<-timer.C
			timer.Reset(duration)

			if status.isReceived() {
				status.setReceived(false)
				break
			}
		}
	}
}
