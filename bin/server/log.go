package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
)

const (
	logLevelTrace = "trace"
	logLevelDebug = "debug"
	logLevelInfo  = "info"
)

var logLevel = flag.String("loglevel", logLevelInfo, "trace / debug / info")

// InitLog ...
func logConfigure() {
	switch *logLevel {
	case logLevelTrace:
		log.SetLevel(log.TraceLevel)
	case logLevelDebug:
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
