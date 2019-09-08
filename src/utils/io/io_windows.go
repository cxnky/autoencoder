package io

import "os"

// IsFileLocked checks whether a file has been fully transferred/written to ensure that we don't get any errors when encoding
func IsFileLocked(fileName string) bool {
	file, err := os.Open(fileName)

	if err != nil {
		return true
	}

	file.Close()
	return false

}
