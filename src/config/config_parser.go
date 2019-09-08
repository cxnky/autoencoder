package config

import (
	"encoding/json"
	"github.com/cxnky/autoencoder/src/logger"
	"io/ioutil"
	"os"
)

type Config struct {
	WatchDirectories []string `json:"watch_directories"`
	EncodeDirectory  string   `json:"encode_directory"`
	DeleteOriginal   bool     `json:"delete_original"`
}

var Configuration Config

func ReadConfig() {
	logger.Info("Reading configuration file")

	bytes, err := ioutil.ReadFile("config.json")

	if err != nil {
		logger.Fatal(err.Error())
	}

	var config Config
	err = json.Unmarshal(bytes, &config)

	if err != nil {
		logger.Fatal(err.Error())
	}

	// Check if the chosen directories exist
	logger.Info("Validating your config")

	if _, err := os.Stat(config.EncodeDirectory); os.IsNotExist(err) {
		logger.Fatal(config.EncodeDirectory + " does not exist on this system")
	}

	for _, dir := range config.WatchDirectories {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			logger.Fatal(dir + " does not exist on this system")
		}
	}

	logger.Info("Config successfully validated")

	Configuration = config

}
