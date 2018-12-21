package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

var (
	settings = map[string]string{
		"root":      ".",
		"tmp_path":  "./tmp",
		"valid_ext": ".go, .tpl, .tmpl, .html, .toml, .yml",
		"ignored":   "assets, tmp, vendor",
	}
)

func handleSig() {
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

func isTmpDir(path string) bool {
	absolutePath, _ := filepath.Abs(path)
	absoluteTmpPath, _ := filepath.Abs(tmpPath())

	return absolutePath == absoluteTmpPath
}

func isIgnoredFolder(path string) bool {
	paths := strings.Split(path, "/")
	if len(paths) <= 0 {
		return false
	}

	for _, e := range strings.Split(settings["ignored"], ",") {
		if strings.TrimSpace(e) == paths[0] {
			return true
		}
	}
	return false
}

func isWatchedFile(path string) bool {
	absolutePath, _ := filepath.Abs(path)
	absoluteTmpPath, _ := filepath.Abs(tmpPath())

	if strings.HasPrefix(absolutePath, absoluteTmpPath) {
		return false
	}

	ext := filepath.Ext(path)

	for _, e := range strings.Split(settings["valid_ext"], ",") {
		if strings.TrimSpace(e) == ext {
			return true
		}
	}

	return false
}

func root() string {
	return settings["root"]
}

func tmpPath() string {
	return settings["tmp_path"]
}

func fileClose(file *os.File) {
	if err := file.Close(); err != nil {
		log.Printf("failed to close file: %s => %v\n", file.Name(), err)
	}
}

func closex() {
	_ = os.RemoveAll(tmpPath())
	os.Exit(1)
}
