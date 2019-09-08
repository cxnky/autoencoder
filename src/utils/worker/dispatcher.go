package worker

import (
	"github.com/cxnky/autoencoder/src/logger"
	"strconv"
)

var WorkerQueue chan chan WorkRequest

func StartDispatcher(nworkers int) {
	WorkerQueue = make(chan chan WorkRequest, nworkers)

	// Create all of the workers
	for i := 0; i < nworkers; i++ {
		logger.Info("Starting worker #" + strconv.Itoa(i+1))
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				go func() {
					worker := <-WorkerQueue
					worker <- work
				}()
			}
		}
	}()

}
