package config

import (
	"encoding/json"
	"github.com/cxnky/vid-encoder/src/logger"
	"io/ioutil"
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

	Configuration = config

}
