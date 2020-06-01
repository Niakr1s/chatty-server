package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/server"

	log "github.com/sirupsen/logrus"
)

const (
	logLevelTrace = "trace"
	logLevelDebug = "debug"
	logLevelInfo  = "info"
)

var logLevel = flag.String("loglevel", logLevelInfo, "trace / debug / info")

func logConfigure() {
	switch *logLevel {
	case logLevelTrace:
		log.SetLevel(log.TraceLevel)
	case logLevelDebug:
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	log.Printf("Log level set to %v", log.GetLevel())
}

func main() {
	flag.Parse()
	logConfigure()
	config.InitConfig()
	rand.Seed(time.Now().UnixNano())

	server := server.NewMemoryServer()
	log.Fatal(server.ListenAndServe())
}
