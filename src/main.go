package main

import (
	"github.com/cxnky/vid-encoder/src/config"
	"github.com/cxnky/vid-encoder/src/logger"
	"github.com/fsnotify/fsnotify"
)

func main() {
	logger.InitialiseLogger()
	config.ReadConfig()
	//bootstrapper.StartBootstrap()

	logger.Info("Initialising file system watcher")
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		logger.Fatal(err.Error())
	}

	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Create == fsnotify.Create {
					// todo: detect lock on both windows and linux
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				logger.Error(err.Error())

			}
		}
	}()

	for _, dir := range config.Configuration.WatchDirectories {
		err = watcher.Add(dir)

		if err != nil {
			logger.Error("Unable to watch directory " + dir)
		} else {
			logger.Info("Watching " + dir)
		}
	}

	<-done

}
