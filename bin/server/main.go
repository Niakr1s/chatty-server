package main

import (
	"log"
	userdb "server2/app/db/user"
	"server2/app/server"
	"server2/app/store"
)

func main() {
	server := server.NewServer(store.NewStore(userdb.NewMemoryDB()))
	log.Fatal(server.ListenAndServe())
}
