package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"

	log "github.com/sirupsen/logrus"
)

// CookieStoreSecretKey is key for Cookie store
// provided by env arg SECRET_KEY
var CookieStoreSecretKey string

// Config ...
type Config struct {
	ServerListenAddress string

	CleanInactiveUsersInterval duration
	InactivityTimeout          duration

	CookieMaxAge int

	RequestTimeout  duration
	ResponseTimeout duration

	Chats []string

	LastMessages int
}

// C contains configuration for app
// Config filename : "config.toml"
var C *Config

var configFilepath = flag.String("config", "config.toml", "provide custom config path")

// InitConfig inits Config and CookieStoreSecretKey.
func InitConfig() {
	C = new(Config)
	C = NewDefaultConfig()
	if _, err := toml.DecodeFile(*configFilepath, C); err != nil {
		log.Warnf("couldn't load config: %v, initializing default config", configFilepath, err)
	} else {
		log.Infof("config file succesfully loaded from %s", configFilepath)
	}

	CookieStoreSecretKey = os.Getenv("SECRET_KEY")
	if CookieStoreSecretKey == "" {
		log.Warnf("Using default secret key, for security reasons provide it via env SECRET_KEY arg")
		CookieStoreSecretKey = "1234567890"
	}
}

func init() {
	C = NewDefaultConfig()
}

// NewDefaultConfig ...
func NewDefaultConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Infof("env port: %v", os.Getenv("PORT"))
	return &Config{
		ServerListenAddress:        fmt.Sprintf(":%s", port),
		CleanInactiveUsersInterval: duration{time.Second * 60},
		InactivityTimeout:          duration{time.Second * 60},

		CookieMaxAge: 60 * 60 * 24 * 7, // week

		RequestTimeout:  duration{time.Second * 15},
		ResponseTimeout: duration{time.Second * 30},

		Chats: []string{"Main"},

		LastMessages: 50,
	}
}
