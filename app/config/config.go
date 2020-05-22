package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

// Config ...
type Config struct {
	ServerListenAddress string
}

// C contains configuration for app
// Config filename : "config.toml"
var C Config

const configFilepath = "config.toml"

func init() {
	if _, err := toml.DecodeFile(configFilepath, &C); err != nil {
		log.Print(err)
	} else {
		log.Printf("config file succesfully loaded from %s", configFilepath)
	}
}
