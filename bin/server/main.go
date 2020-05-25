package main

import (
	"flag"

	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/server"

	log "github.com/sirupsen/logrus"
)

func main() {
	flag.Parse()
	logConfigure()

	server := server.NewMemoryServer()
	db.StartCleanInactiveUsers(server.Store.LoggedDB,
		config.C.CleanInactiveUsersInterval.Duration,
		config.C.InactivityTimeout.Duration)

	log.Fatal(server.ListenAndServe())
}
