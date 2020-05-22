package main

import (
	"server2/app/server"

	log "github.com/sirupsen/logrus"
)

func main() {
	server := server.NewMemoryServer()
	log.Fatal(server.ListenAndServe())
}
