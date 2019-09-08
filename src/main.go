package main

import (
	"github.com/cxnky/vid-encoder/src/config"
	"github.com/cxnky/vid-encoder/src/logger"
	"github.com/cxnky/vid-encoder/src/utils"
	"github.com/fsnotify/fsnotify"
	"path/filepath"
	"time"
)

func main() {
	logger.InitialiseLogger()
	config.ReadConfig()

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
					if filepath.Ext(event.Name) != ".crdownload" {
						for {
							if utils.IsFileLocked(event.Name) {
								logger.Debug("File is locked")
							} else {
								logger.Debug("File is not locked")
								break
							}

							time.Sleep(500 * time.Millisecond)

						}
					}

					// File is no longer locked, we can encode now

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
