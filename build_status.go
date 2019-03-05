package main

import (
	"log"
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
	if _, err := os.Create(bs.errFile); err != nil {
		log.Printf(`failed to create error file => %+v`, err)
	}
}

func (bs *buildStatus) success() {
	bs.mux.Lock()
	defer bs.mux.Unlock()
	if err := os.Remove(bs.errFile); err != nil {
		log.Printf(`failed to remove error file => %+v`, err)
	}
}
