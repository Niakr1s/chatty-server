package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	runServer()
}

func runServer() {
	log.Infof("Initializing server, dev=%v", *dev)
	server, err := server.New(*dev)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})

	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt)
	signal.Notify(exit, syscall.SIGTERM)

	go func() {
		s := <-exit
		log.Infof("Got %v, exiting...", s)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Errorf("server shutdown fail: %v", err)
		}
		done <- struct{}{}
	}()

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("listen err: %v", err)
		}
	}()

	<-done

	log.Infof("server shutdown succesfully")
}
