package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
)

var (
	binFile, lastHash string
	bs                *buildStatus
	watchChan         chan string
)

func init() {
	binFile = path.Join(tmpPath(), `server`)
	bs = &buildStatus{errFile: path.Join(tmpPath(), `error`)}
	watchChan = make(chan string, 1000)
}

func main() {
	flag.Parse()

	go watch()
	handleSig()
	start()
}

func start() {
	var build, server *exec.Cmd
	debounce := newDebouncer(watchChan, 5)

	for {
		build = exec.Command(`go`, `build`, `-i` ,`-o`, binFile)
		build.Stdout = os.Stdout
		build.Stderr = os.Stderr

		if err := build.Run(); err != nil {
			bs.error()
			log.Printf(`failed to build source => %v`, err)
		} else {
			bs.success()
			if should, err := shouldRestart(); err != nil {
				log.Printf("failed to check server binary => %v", err)
			} else if should {
				if server != nil && server.Process != nil {
					log.Println(`[go-watch] restarting...`)
					if err := server.Process.Kill(); err != nil {
						log.Fatalf(`failed to terminate server process => %v`, err)
					}
				} else {
					log.Println(`[go-watch] start`)
				}

				server = exec.Command(binFile, flag.Args()...)
				server.Stdout = os.Stdout
				server.Stderr = os.Stderr

				if err := server.Start(); err != nil {
					log.Fatalf(`failed to start server process => %v`, err)
				}
			}
		}

		debounce()
	}
}

func shouldRestart() (bool, error) {
	if h, err := binHash(); err != nil {
		return false, err
	} else if h != lastHash {
		lastHash = h
		return true, nil
	}

	return false, nil
}

func binHash() (string, error) {
	f, err := os.Open(binFile)
	if err != nil {
		return ``, err
	}
	defer fileClose(f)

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return ``, err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
