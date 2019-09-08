package worker

import (
	"encoding/json"
	"github.com/cxnky/autoencoder/src/logger"
)

type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
}

func NewWorker(id int, workerQueue chan chan WorkRequest) Worker {
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
	}

	return worker

}

type EncodingJob struct {
	FilePath string `json:"file_path"`
	FileType string `json:"file_type"`
}

func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				var encodingJob EncodingJob
				err := json.Unmarshal(work.Data, &encodingJob)

				if err != nil {
					logger.Error("Unable to parse job JSON: " + err.Error())
					break
				}

				// todo: encode the file here
				// todo: ask the user what format they want

			}

		}
	}()
}
