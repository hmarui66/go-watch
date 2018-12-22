package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				if err = handleEvent(watcher, ev); err != nil {
					log.Printf("failed to handle event: %v => %v", ev, err)
				}
			case err = <-watcher.Errors:
				log.Printf("error: %s", err)
			}
		}
	}()

	err = filepath.Walk(root(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf(`failed to walk a path: %s => %v`, path, err)
			return err
		}

		return addPath(watcher, path, info)
	})

	if err == filepath.SkipDir {
		log.Println(err)
	} else if err != nil {
		log.Fatal(err)
	}
}

func addPath(watcher *fsnotify.Watcher, path string, info os.FileInfo) error {
	if info.IsDir() && !isTmpDir(path) {
		if len(path) > 1 && strings.HasPrefix(filepath.Base(path), `.`) {
			return filepath.SkipDir
		}

		if isIgnoredFolder(path) {
			return filepath.SkipDir
		}

		return watcher.Add(path)
	}

	return nil
}

func handleEvent(watcher *fsnotify.Watcher, ev fsnotify.Event) error {
	if isWatchedFile(ev.Name) {
		watchChan <- ev.String()
		return nil
	}

	if ev.Op != fsnotify.Create {
		return nil
	}

	fi, err := os.Lstat(ev.Name)
	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return nil
	}

	if ev.Op == fsnotify.Create {
		if err = addPath(watcher, ev.Name, fi); err != nil {
			return err
		}
	}

	return nil
}
