package config

import (
	"time"

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
var C *Config

const configFilepath = "config.toml"

func init() {
	C = new(Config)
	if _, err := toml.DecodeFile(configFilepath, C); err != nil {
		log.Infof("couldn't load from %s: %v, initializing default config", configFilepath, err)
		C = NewDefaultConfig()
	} else {
		log.Infof("config file succesfully loaded from %s", configFilepath)
	}
}

// NewDefaultConfig ...
func NewDefaultConfig() *Config {
	return &Config{
		ServerListenAddress:        "127.0.0.1:8080",
		CleanInactiveUsersInterval: duration{time.Second * 60},
		InactivityTimeout:          duration{time.Second * 60},
	}
}
