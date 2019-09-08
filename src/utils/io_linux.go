package utils

import (
	"github.com/cxnky/autoencoder/src/logger"
	"os"
)

var previousFileSize int64 = 0

// IsFileLocked checks whether a file has been fully transferred/written to ensure that we don't get any errors when encoding
func IsFileLocked(fileName string) bool {
	fi, err := os.Stat(fileName)

	if err != nil {
		logger.Error("Error whilst getting file info: " + err.Error())
	}

	size := fi.Size()

	// File is no longer being written to
	if previousFileSize == size {
		return false
	} else {
		previousFileSize = size
		return true
	}

}
