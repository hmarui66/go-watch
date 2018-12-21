package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

var (
	binFile, lastHash string
	watchChan         chan string
)

func init() {
	rand.Seed(time.Now().UnixNano())
	binFile = tmpPath() + `/go-watch-tmp`
	watchChan = make(chan string, 1000)
}

func main() {
	handleSig()
	watch()
	start()
}

func start() {
	var server *exec.Cmd

	for {
		buildCmd := exec.Command(`go`, `build`, `-o`, binFile)
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr

		if err := buildCmd.Run(); err != nil {
			log.Println(`failed to build source`)
		}

		if shouldRestart() {

			if server != nil && server.Process != nil {
				log.Println(`[go-watch] restarting...`)
				if err := server.Process.Kill(); err != nil {
					log.Fatal(`failed to terminate server process`)
				}
			} else {
				log.Println(`[go-watch] start`)
			}

			server = exec.Command(binFile)
			server.Stdout = os.Stdout
			server.Stderr = os.Stderr

			if err := server.Start(); err != nil {
				log.Fatal(err)
			}
		}
		<-watchChan
	}
}

func shouldRestart() bool {
	h := binHash()
	if h != lastHash {
		lastHash = h
		return true
	}

	return false
}

func binHash() string {
	f, err := os.Open(binFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fileClose(f)

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
