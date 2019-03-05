package main

import (
	"os"
	"sync"
)

type buildStatus struct {
	errFile string
	mux     sync.Mutex
}

func (bs *buildStatus) error() {
	bs.mux.Lock()
	defer bs.mux.Unlock()
	os.Create(bs.errFile)
}

func (bs *buildStatus) success() {
	bs.mux.Lock()
	defer bs.mux.Unlock()
	os.Remove(bs.errFile)
}
