package main

import (
	"github.com/cxnky/autoencoder/src/config"
	"github.com/cxnky/autoencoder/src/logger"
	"github.com/cxnky/autoencoder/src/utils/io"
	"github.com/cxnky/autoencoder/src/utils/worker"
	"github.com/fsnotify/fsnotify"
	"github.com/gen2brain/beeep"
	"path/filepath"
	"strconv"
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
							if io.IsFileLocked(event.Name) {
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

	logger.Info("Starting dispatcher")
	worker.StartDispatcher(config.Configuration.ConcurrentWorkers)

	err = beeep.Notify("AutoEncoder", "AutoEncoder is now watching "+strconv.Itoa(len(config.Configuration.WatchDirectories))+" directories", "")

	if err != nil {
		logger.Warning("Unable to show notification, AutoEncoder is now ready.")
	}

	<-done

}
