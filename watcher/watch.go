package watcher

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
)

func StartWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Write) {
					eventCh <- "reload"
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	if *root == "." {
		wd, _ := os.Getwd()
		err = watcher.Add(wd)
	} else {
		err = watcher.Add(*root)
	}

	if err != nil {
		log.Fatal(err)
	}

	<-make(chan struct{})
}
