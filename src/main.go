package main

import (
	"encoding/json"
	"fmt"
	"github.com/cxnky/autoencoder/src/config"
	"github.com/cxnky/autoencoder/src/logger"
	"github.com/cxnky/autoencoder/src/utils/io"
	"github.com/cxnky/autoencoder/src/utils/worker"
	"github.com/fsnotify/fsnotify"
	"github.com/gen2brain/beeep"
	"path/filepath"
	"time"
)

func main() {
	logger.InitialiseLogger()
	config.ReadConfig()

	if len(config.Configuration.WatchDirectories) == 0 {
		logger.Fatal("You need to specify at least one watch directory!")
	}

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

							time.Sleep(100 * time.Millisecond)

						}
					}

					workData := worker.EncodingJob{
						FilePath: event.Name,
					}

					bytes, err := json.Marshal(workData)

					if err != nil {
						logger.Error("Unable to marshal JSON: " + err.Error())
						break
					}

					// File is no longer locked, we can encode now
					workRequest := worker.WorkRequest{
						Data: bytes,
						Type: "video",
					}

					worker.QueueWork(workRequest)
					break

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

	watchingString := ""

	if len(config.Configuration.WatchDirectories) == 1 {
		watchingString = "1 directory"
	} else {
		watchingString = fmt.Sprintf("%d directories", len(config.Configuration.WatchDirectories))
	}

	err = beeep.Notify("AutoEncoder", "AutoEncoder is now watching "+watchingString, "")

	if err != nil {
		logger.Warning("Unable to show notification, AutoEncoder is now ready.")
	}

	<-done

}
