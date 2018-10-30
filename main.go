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
	"os/signal"
	"syscall"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	lastBinHash, cache string
)

func main() {
	handleSignal()

	cache = `./tmp/` + randStr(24)

	var server *exec.Cmd

	for {
		exec.Command(`go`, `build`, `-o`, cache).Run()
		if shouldRestart() {

			if server != nil && server.Process != nil {
				log.Println(`[go-watch] restarting...`)
				server.Process.Kill()
			} else {
				log.Println(`[go-watch] start`)
			}

			server = exec.Command(`./` + cache)

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
		time.Sleep(5 * time.Second)
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func handleSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		for {
			s := <-ch
			switch s {
			case syscall.SIGINT:
				closex()
			case syscall.SIGHUP:
				closex()
			case syscall.SIGTERM:
				closex()
			case syscall.SIGQUIT:
				closex()
			}
		}
	}()
}

func shouldRestart() bool {
	h := binHash(cache)
	if h != lastBinHash {
		lastBinHash = h
		return true
	}

	return false
}

func binHash(cache string) string {
	f, err := os.Open(`./` + cache)
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

func closex() {
	os.RemoveAll(`./tmp`)
	os.Exit(1)
}
