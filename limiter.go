package main

import (
	"sync"
	"time"
)

type watchStatus struct {
	watched bool
	mux     sync.Mutex
}

func (w *watchStatus) setWatched(b bool) {
	w.mux.Lock()
	defer w.mux.Unlock()
	w.watched = b
}

func (w *watchStatus) isWatched() bool {
	w.mux.Lock()
	defer w.mux.Unlock()
	return w.watched
}

func newLimitter(sender chan string, d int) func() {
	status := watchStatus{}

	go func() {
		for {
			<-sender
			status.setWatched(true)
		}
	}()

	tick := time.Tick(time.Duration(d) * time.Second)

	return func() {
		for {
			<-tick
			if status.isWatched() {
				status.setWatched(false)
				break
			}
		}
	}
}
