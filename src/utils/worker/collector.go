package worker

import (
	"github.com/cxnky/autoencoder/src/config"
	"github.com/cxnky/autoencoder/src/logger"
)

var WorkQueue = make(chan WorkRequest, config.Configuration.MaxQueueLength)

func QueueWork(request WorkRequest) {
	logger.Info("Work request queued")
	WorkQueue <- request
}
