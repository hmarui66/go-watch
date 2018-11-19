package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func watch() {
	fileCount := 0
	filepath.Walk(root(), func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && !isTmpDir(path) {
			if len(path) > 1 && strings.HasPrefix(filepath.Base(path), `.`) {
				return filepath.SkipDir
			}

			if isIgnoredFolder(path) {
				return filepath.SkipDir
			}

			watchFolder(path)
			fileCount++
		}
		return err
	})
	fmt.Printf("watching target file count: %d\n", fileCount)
}

func watchFolder(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				if isWatchedFile(ev.Name) {
					watchChan <- ev.String()
				}
			case err := <-watcher.Errors:
				log.Printf("error: %s", err)
			}
		}
	}()

	err = watcher.Add(path)

	if err != nil {
		log.Fatal(err)
	}
}
