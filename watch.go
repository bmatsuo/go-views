package main

import (
	"os"

	"github.com/howeyc/fsnotify"
)

func watch(sig chan<- string, pathmatch func(string) bool, dirs ...string) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		errLogger.Println("error creating watcher: ", err)
		return nil, err
	}

	go func() {
		for cont := true; cont; {
			select {
			case ev := <-watcher.Event:
				outLogger.Println("filesystem event:", ev)
				path := ev.Name
				if ev.IsCreate() {
					info, err := os.Stat(path)
					if err != nil {
						errLogger.Println("stat error: ", err)
						continue
					}
					if info.IsDir() {
						outLogger.Println("watching: ", path)
						err := watcher.Watch(path)
						if err != nil {
							errLogger.Println("add watch error: ", err)
						}
					}
				}

				if pathmatch(path) {
					sig <- path
				}
			case err := <-watcher.Error:
				errLogger.Println("watcher error:", err)
			}
		}
	}()

	return watcher, err
}
