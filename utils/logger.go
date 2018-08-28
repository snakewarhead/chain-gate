package utils

import (
	log "github.com/jeanphorn/log4go"
)

const (
	configPath = "./resources/log.json"
	// configFile = "/home/lin/Work/go-bench/src/github.com/snakewarhead/chain-gate/resources/log.json"

	category = "gate"
)

var (
	// Logger is export,usage: utils.Logger.Debug("test")
	Logger *log.Filter
)

func init() {
	log.LoadConfiguration(configPath)

	Logger = log.LOGGER(category)
}
