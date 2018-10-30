package main

import (
	"bufio"
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
	binFile = fmt.Sprintf(`%s/watch-%s`, tmpPath(), randStr(12))
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
		exec.Command(`go`, `build`, `-o`, binFile).Run()
		if shouldRestart() {

			if server != nil && server.Process != nil {
				log.Println(`[go-watch] restarting...`)
				server.Process.Kill()
			} else {
				log.Println(`[go-watch] start`)
			}

			server = exec.Command(binFile)

			stdout, err := server.StdoutPipe()
			if err != nil {
				log.Fatal(err)
			}
			scanner := bufio.NewScanner(stdout)
			go func() {
				for scanner.Scan() {
					fmt.Println(scanner.Text())
				}
			}()

			stderr, err := server.StderrPipe()
			if err != nil {
				log.Fatal(err)
			}
			errScanner := bufio.NewScanner(stderr)
			go func() {
				for errScanner.Scan() {
					fmt.Println(errScanner.Text())
				}
			}()

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
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))

}
