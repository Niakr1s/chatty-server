package main

import (
	"flag"
	"server2/app/server"

	log "github.com/sirupsen/logrus"
)

func main() {
	flag.Parse()
	logConfigure()

	server := server.NewMemoryServer()
	log.Fatal(server.ListenAndServe())
}
