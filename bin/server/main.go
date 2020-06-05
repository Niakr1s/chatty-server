package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/constants"
	"github.com/niakr1s/chatty-server/app/server"

	log "github.com/sirupsen/logrus"
)

const (
	logLevelTrace = "trace"
	logLevelDebug = "debug"
	logLevelInfo  = "info"
)

var logLevel = flag.String("loglevel", logLevelInfo, "trace / debug / info")
var dev = flag.Bool("dev", false, fmt.Sprintf("dev mode, can be used with %s to use persistent database", constants.EnvDatabaseURL))

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
	rand.Seed(time.Now().UTC().UnixNano())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error
	var s *server.Server
	switch *dev {
	case false: // it's prod
		log.Infof("Initializing prod server...")
		s, err = server.NewProdServer(ctx)
	default: // it's dev
		log.Infof("Initializing dev server...")
		s, err = server.NewDevServer(ctx)
	}
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(s.ListenAndServe())
}
