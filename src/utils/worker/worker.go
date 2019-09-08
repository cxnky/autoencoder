package worker

import (
	"encoding/json"
	"fmt"
	"github.com/cxnky/autoencoder/src/config"
	"github.com/cxnky/autoencoder/src/logger"
	"github.com/xfrr/goffmpeg/transcoder"
	"path/filepath"
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

				fileExt := filepath.Ext(encodingJob.FilePath)

				if fileExt == ".crdownload" {
					logger.Info("Stopping work request as it is for a partial download file")
					break
				}

				// todo: look into job duplication issue
				transcoder := new(transcoder.Transcoder)
				fileName := filepath.Base(encodingJob.FilePath)

				err = transcoder.Initialize(encodingJob.FilePath, config.Configuration.EncodeDirectory+"/"+fileName+"."+config.Configuration.OutputFormat)

				if err != nil {
					logger.Error("Unable to transcode: " + err.Error())
					break
				}

				done := transcoder.Run(true)
				progress := transcoder.Output()

				for msg := range progress {
					logger.Debug(fmt.Sprintf("Progress: %f, bitrate: %s, speed: %s", msg.Progress, msg.CurrentBitrate, msg.Speed))
				}

				err = <-done

			}

		}
	}()
}
