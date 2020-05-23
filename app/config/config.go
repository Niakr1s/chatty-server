package config

import (
	"github.com/BurntSushi/toml"

	log "github.com/sirupsen/logrus"
)

// Config ...
type Config struct {
	ServerListenAddress string

	CleanInactiveUsersInterval duration
	InactivityTimeout          duration
}

// C contains configuration for app
// Config filename : "config.toml"
var C Config

const configFilepath = "config.toml"

func init() {
	if _, err := toml.DecodeFile(configFilepath, &C); err != nil {
		log.Error(err)
	} else {
		log.Printf("config file succesfully loaded from %s", configFilepath)
	}
}
