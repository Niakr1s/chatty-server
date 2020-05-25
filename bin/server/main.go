package main

import (
	"flag"

	"github.com/niakr1s/chatty-server/app/server"

	log "github.com/sirupsen/logrus"
)

func main() {
	flag.Parse()
	logConfigure()

	server := server.NewMemoryServer().WithPool()
	log.Fatal(server.ListenAndServe())
}
