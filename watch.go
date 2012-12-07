package main

import (
	"os"

	"github.com/howeyc/fsnotify"
)

func watch(sig chan<- *fsnotify.FileEvent, match func(*fsnotify.FileEvent) bool) (w *fsnotify.Watcher, err error) {
	w, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	go _watch(w, sig, match)
	return
}

func _watch(watcher *fsnotify.Watcher, sig chan<- *fsnotify.FileEvent, match func(*fsnotify.FileEvent) bool) {
	for cont := true; cont; {
		select {
		case ev := <-watcher.Event:
			Debug(2, "event:", ev)
			if ev.IsCreate() {
				info, err := os.Stat(ev.Name)
				if err != nil {
					Debug(2, "stat error: ", err)
					continue
				}
				if info.IsDir() {
					Debug(1, "watching: ", ev.Name)
					err := watcher.Watch(ev.Name)
					if err != nil {
						Error("watch error: ", err)
					}
				}
			}

			if match(ev) {
				sig <- ev
			}
		case err := <-watcher.Error:
			Print("watcher error:", err)
		}
	}
	close(sig)
}
